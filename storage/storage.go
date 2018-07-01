package storage

import (
	"syscall/js"
)

type Storage interface {
	// Length returns an integer representing the number of data items stored in the Storage object.
	Length() int
	// Key will return the name of the nth key in the storage.
	Key(ind int) string
	// GetItem will return that key's value.
	GetItem(key string) (string, bool)
	// SetItem will add that key to the storage, or update that key's value if it already exists.
	SetItem(key, val string)
	// RemoveItem will remove that key from the storage.
	RemoveItem(key string)
	// Clear will empty all keys out of the storage.
	Clear()
}

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
