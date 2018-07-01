package storage

import (
	"encoding/json"
	"errors"
)

var ErrNotFound = errors.New("not found")

func GetItemJSON(s Storage, key string, dst interface{}) error {
	v, ok := s.GetItem(key)
	if !ok {
		return ErrNotFound
	}
	return json.Unmarshal([]byte(v), dst)
}

func SetItemJSON(s Storage, key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	s.SetItem(key, string(data))
	return nil
}
