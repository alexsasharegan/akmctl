package api

import (
	"os"
	"path/filepath"

	"github.com/alexsasharegan/akmctl/internal/migrate"
	"github.com/alexsasharegan/akmctl/internal/record"
)

// previewOps are the writes a preview can apply (ADR 261109-vgg7).
var previewOps = map[string]bool{"new": true, "set": true, "section": true, "tick": true, "remove": true, "log": true, "close": true, "claim": true, "release": true}

// newPlaceholder in a string of a write stands for the id of the record made by
// the most recent `new` before it.
const newPlaceholder = "$new"

// opPreview applies an ordered list of writes to a scratch copy of the
// repository and returns, for each, the response its dry run gives after the
// earlier writes. It never writes to the repository itself.
func opPreview(repo *record.Repo, info Info, q Request) (Obj, *Error) {
	var writes []map[string]any
	switch l := q["writes"].(type) {
	case []any:
		for _, x := range l {
			writes = append(writes, x.(map[string]any)) // validate has checked the type
		}
	case []map[string]any:
		writes = l
	}
	if len(writes) == 0 {
		return nil, invalid("writes", "writes must hold at least one write")
	}
	// check every write before copying anything, so a mistake names its position
	for i, w := range writes {
		op, _ := w["op"].(string)
		if !previewOps[op] {
			return nil, atWrite(i, invalid("writes", "writes[%d]: op must be one of new, set, section, tick, remove, log, close, claim, release", i))
		}
	}
	sc, cleanup := migrate.ScratchRepo(repo)
	if sc == nil {
		return nil, &Error{Code: "unreadable", Message: "cannot make a scratch copy of the repository"}
	}
	defer cleanup()

	etags := map[string]string{} // the etag of each record on the scratch copy after its last write
	newID := ""
	out := make([]Obj, 0, len(writes))
	for i, w := range writes {
		op, _ := w["op"].(string)
		req := Request{}
		for k, v := range w {
			if k != "op" {
				req[k] = substituteNew(v, newID)
			}
		}
		if _, ok := specs[op]["as_of"]; ok && req["as_of"] == nil && q["as_of"] != nil {
			req["as_of"] = q["as_of"]
		}
		if id, ok := req["id"].(string); ok && etags[id] != "" {
			if _, has := specs[op]["if_match"]; has {
				req["if_match"] = etags[id]
			}
		}
		if err := validate(op, req); err != nil {
			return nil, atWrite(i, err)
		}
		if op == "new" {
			resp, err := Run(sc, info, "new", req)
			if err != nil {
				return nil, atWrite(i, err)
			}
			rec, _ := resp.Get("record")
			id, _ := rec.(Obj).Get("id")
			path, _ := rec.(Obj).Get("path")
			newID, _ = id.(string)
			text, _ := os.ReadFile(filepath.Join(sc.Root, filepath.FromSlash(path.(string))))
			resp = append(resp, KV{"diff", Diff("", path.(string), "", string(text))}, KV{"note", "The id is chosen when the write is made; this one is only for the preview."})
			out = append(out, resp)
			continue
		}
		dry := Request{"dry_run": true}
		for k, v := range req {
			dry[k] = v
		}
		resp, err := Run(sc, info, op, dry)
		if err != nil {
			return nil, atWrite(i, err)
		}
		out = append(out, resp)
		// make the write on the scratch copy so the next one sees it
		real, rerr := Run(sc, info, op, req)
		if rerr != nil {
			return nil, atWrite(i, rerr)
		}
		if id, ok := req["id"].(string); ok {
			if rec, ok := real.Get("record"); ok {
				if et, ok := rec.(Obj).Get("etag"); ok {
					etags[id], _ = et.(string)
				}
			}
		}
	}
	return O("writes", out), nil
}

// atWrite adds the position of the write that failed to an error's details.
func atWrite(i int, e *Error) *Error {
	return &Error{Code: e.Code, Message: e.Message, Details: append(Obj{{"write", i}}, e.Details...)}
}

// substituteNew replaces the placeholder in every string of a value.
func substituteNew(v any, id string) any {
	if id == "" {
		return v
	}
	switch x := v.(type) {
	case string:
		if x == newPlaceholder {
			return id
		}
	case []any:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = substituteNew(e, id)
		}
		return out
	case []string:
		out := make([]any, len(x))
		for i, e := range x {
			out[i] = substituteNew(e, id)
		}
		return out
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, e := range x {
			out[k] = substituteNew(e, id)
		}
		return out
	}
	return v
}
