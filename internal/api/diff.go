package api

import (
	"strconv"
	"strings"
)

const diffContext = 3

// maxLCSCells bounds the table of the line diff. A change region larger than
// this (a rewrite of a huge file) is reported as one block of removed lines
// followed by one block of added lines.
const maxLCSCells = 4_000_000

type diffOp struct {
	op       string // "ctx", "add" or "rem"
	text     string
	old, new int // 1-based line numbers; 0 when the line is not in that side
}

// splitLines cuts text into lines. Each line keeps a mark for a missing final
// newline, so that adding or removing the last newline shows as a change.
func splitLines(s string) (text []string, key []string) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, "\n")
	for i, p := range parts {
		last := i == len(parts)-1
		if last && p == "" {
			break
		}
		text = append(text, p)
		if last {
			key = append(key, p+"\x00")
		} else {
			key = append(key, p)
		}
	}
	return text, key
}

// diffOps computes a line diff of a to b: the common prefix and suffix are cut
// off, then a longest-common-subsequence table orders the rest.
func diffOps(a, b string) []diffOp {
	at, ak := splitLines(a)
	bt, bk := splitLines(b)
	pre := 0
	for pre < len(ak) && pre < len(bk) && ak[pre] == bk[pre] {
		pre++
	}
	suf := 0
	for suf < len(ak)-pre && suf < len(bk)-pre && ak[len(ak)-1-suf] == bk[len(bk)-1-suf] {
		suf++
	}
	var ops []diffOp
	for i := 0; i < pre; i++ {
		ops = append(ops, diffOp{"ctx", at[i], i + 1, i + 1})
	}
	ma, mb := ak[pre:len(ak)-suf], bk[pre:len(bk)-suf]
	oa, ob := at[pre:len(at)-suf], bt[pre:len(bt)-suf]
	n, m := len(ma), len(mb)
	if n*m > maxLCSCells || (n+1)*(m+1) > maxLCSCells {
		for i := 0; i < n; i++ {
			ops = append(ops, diffOp{"rem", oa[i], pre + i + 1, 0})
		}
		for j := 0; j < m; j++ {
			ops = append(ops, diffOp{"add", ob[j], 0, pre + j + 1})
		}
	} else {
		// l[i][j] is the length of the LCS of ma[i:] and mb[j:]
		w := m + 1
		l := make([]int32, (n+1)*w)
		for i := n - 1; i >= 0; i-- {
			for j := m - 1; j >= 0; j-- {
				if ma[i] == mb[j] {
					l[i*w+j] = l[(i+1)*w+j+1] + 1
				} else if l[(i+1)*w+j] >= l[i*w+j+1] {
					l[i*w+j] = l[(i+1)*w+j]
				} else {
					l[i*w+j] = l[i*w+j+1]
				}
			}
		}
		i, j := 0, 0
		for i < n || j < m {
			switch {
			case i < n && j < m && ma[i] == mb[j]:
				ops = append(ops, diffOp{"ctx", oa[i], pre + i + 1, pre + j + 1})
				i++
				j++
			case i < n && (j == m || l[(i+1)*w+j] >= l[i*w+j+1]):
				ops = append(ops, diffOp{"rem", oa[i], pre + i + 1, 0})
				i++
			default:
				ops = append(ops, diffOp{"add", ob[j], 0, pre + j + 1})
				j++
			}
		}
	}
	for k := 0; k < suf; k++ {
		ops = append(ops, diffOp{"ctx", at[len(at)-suf+k], len(at) - suf + k + 1, len(bt) - suf + k + 1})
	}
	return ops
}

// Diff builds the diff member of a dry-run response: the change from the old
// file text to the new, as unified hunks with three lines of context. path and
// newPath are the file before and after (they differ for a close). Line text is
// the line without its LF; a CR before the LF stays in the text.
func Diff(path, newPath, old, new string) Obj {
	ops := diffOps(old, new)
	hunks := []Obj{}
	i := 0
	for i < len(ops) {
		for i < len(ops) && ops[i].op == "ctx" {
			i++
		}
		if i == len(ops) {
			break
		}
		from := i - diffContext
		if from < 0 {
			from = 0
		}
		// extend over changes that are close enough to share context
		end, lastChange := i, i
		for end < len(ops) {
			if ops[end].op != "ctx" {
				lastChange = end
			} else if end-lastChange > 2*diffContext {
				break
			}
			end++
		}
		to := lastChange + 1 + diffContext
		if to > len(ops) {
			to = len(ops)
		}
		hunks = append(hunks, hunk(ops[from:to]))
		i = to
	}
	return O("path", path, "new_path", newPath, "hunks", hunks)
}

func hunk(ops []diffOp) Obj {
	var oldStart, newStart, oldN, newN int
	lines := make([]Obj, 0, len(ops))
	for _, o := range ops {
		var ov, nv any
		if o.old > 0 {
			ov = o.old
			oldN++
			if oldStart == 0 {
				oldStart = o.old
			}
		}
		if o.new > 0 {
			nv = o.new
			newN++
			if newStart == 0 {
				newStart = o.new
			}
		}
		lines = append(lines, O("op", o.op, "text", o.text, "old", ov, "new", nv))
	}
	// a side with no lines (an empty file) starts at 0
	return O("header", "@@ -"+span(oldStart, oldN)+" +"+span(newStart, newN)+" @@", "lines", lines)
}

func span(start, n int) string {
	if n == 1 {
		return strconv.Itoa(start)
	}
	return strconv.Itoa(start) + "," + strconv.Itoa(n)
}
