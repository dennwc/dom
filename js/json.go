//+build wasm,js

package js

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var (
	_ json.Marshaler   = Value{}
	_ json.Unmarshaler = (*Value)(nil)
)

var (
	jsonObj       Ref
	jsonParse     Ref
	jsonStringify Ref
	jsonOnce      sync.Once
)

func initJSON() {
	jsonObj = global.Get("JSON")
	if jsonObj == undefined {
		return
	}
	jsonParse = jsonObj.Get("parse")
	jsonStringify = jsonObj.Get("stringify")
}

// MarshalJSON encodes a value into JSON by using native JavaScript function (JSON.stringify).
func (v Value) MarshalJSON() ([]byte, error) {
	jsonOnce.Do(initJSON)
	if jsonStringify == undefined {
		return nil, errors.New("json encoding is not supported")
	}
	if v.Ref == undefined {
		return []byte("null"), nil
	}
	s := jsonStringify.Invoke(v.Ref).String()
	return []byte(s), nil
}

// UnmarshalJSON decodes a value from JSON by using native JavaScript functions (JSON.parse).
func (v *Value) UnmarshalJSON(p []byte) (err error) {
	jsonOnce.Do(initJSON)
	if jsonParse == undefined {
		return errors.New("json decoding is not supported")
	}
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	v.Ref = jsonParse.Invoke(string(p))
	return err
}
