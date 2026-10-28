// Package api implements the AKM API operations over a loaded repository: a
// request is a map of named parameters and a response is one JSON object.
package api

import (
	"bytes"
	"encoding/json"
)

// KV is one member of an Obj.
type KV struct {
	K string
	V any
}

// Obj is a JSON object that keeps its member order, so output is stable.
type Obj []KV

// O builds an Obj from alternating keys and values.
func O(kv ...any) Obj {
	o := make(Obj, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		o = append(o, KV{kv[i].(string), kv[i+1]})
	}
	return o
}

// Get returns the member with key k.
func (o Obj) Get(k string) (any, bool) {
	for _, kv := range o {
		if kv.K == k {
			return kv.V, true
		}
	}
	return nil, false
}

func (o Obj) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, kv := range o {
		if i > 0 {
			b.WriteByte(',')
		}
		k, _ := json.Marshal(kv.K)
		b.Write(k)
		b.WriteByte(':')
		v, err := json.Marshal(kv.V)
		if err != nil {
			return nil, err
		}
		b.Write(v)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// marshal encodes without HTML escaping.
func Marshal(v any) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(b.Bytes(), "\n"), nil
}
