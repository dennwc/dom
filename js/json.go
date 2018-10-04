package js

import "encoding/json"

var (
	_ json.Marshaler   = Value{}
	_ json.Unmarshaler = (*Value)(nil)
)

var (
	jsonObj       = global.Get("JSON")
	jsonParse     = jsonObj.Get("parse")
	jsonStringify = jsonObj.Get("stringify")
)

func (v Value) MarshalJSON() ([]byte, error) {
	s := jsonStringify.Invoke(v.Ref).String()
	return []byte(s), nil
}

func (v *Value) UnmarshalJSON(p []byte) error {
	v.Ref = jsonParse.Invoke(string(p))
	return nil // TODO: set error properly
}
