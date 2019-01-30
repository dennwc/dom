package dom

import "github.com/dennwc/dom/js"

type TokenList struct {
	v js.Value
}

func (t *TokenList) Add(class ...interface{}) {
	t.v.Call("add", class...)
}

func (t *TokenList) Remove(class ...interface{}) {
	t.v.Call("remove", class...)
}
