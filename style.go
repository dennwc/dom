//+build wasm,js

package dom

import "github.com/dennwc/dom/js"

type Style struct {
	v js.Value
}

func (s *Style) SetWidth(v Unit) {
	s.v.Set("width", v.String())
}

func (s *Style) SetHeight(v Unit) {
	s.v.Set("height", v.String())
}

func (s *Style) SetMarginsRaw(m string) {
	s.v.Set("margin", m)
}

func (s *Style) Set(k string, v interface{}) {
	s.v.Set(k, v)
}
