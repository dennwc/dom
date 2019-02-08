//+build wasm

package storage

import (
	"syscall/js"
)

func getStorage(name string) Storage {
	s := js.Global().Get("window").Get(name)
	if s == js.Null() || s == js.Undefined() {
		return nil
	}
	return jsStorage{s}
}

func Local() Storage {
	return getStorage("localStorage")
}

func Session() Storage {
	return getStorage("sessionStorage")
}

type jsStorage struct {
	v js.Value
}

func (s jsStorage) Length() int {
	return s.v.Get("length").Int()
}

func (s jsStorage) Key(ind int) string {
	return s.v.Call("key", ind).String()
}

func (s jsStorage) GetItem(key string) (string, bool) {
	v := s.v.Call("getItem", key)
	if v == js.Null() || v == js.Undefined() {
		return "", false
	}
	return v.String(), true
}

func (s jsStorage) SetItem(key, val string) {
	s.v.Call("setItem", key, val)
}

func (s jsStorage) RemoveItem(key string) {
	s.v.Call("removeItem", key)
}

func (s jsStorage) Clear() {
	s.v.Call("clear")
}
