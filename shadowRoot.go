//+build wasm

package dom

import "github.com/dennwc/dom/js"

// https://developer.mozilla.org/en-US/docs/Web/API/ShadowRoot

func AsShadowRoot(v js.Value) *ShadowRoot {
	if !v.Valid() {
		return nil
	}
	return &ShadowRoot{NodeBase{v: v}}
}

var _ Node = (*ShadowRoot)(nil)

type ShadowRoot struct {
	NodeBase
}

func (r *ShadowRoot) IsOpen() bool {
	if r.v.Get("mode").String() == "open" {
		return true
	}
	return false
}

func (r *ShadowRoot) Host() *Element {
	return AsElement(r.v.Get("host"))
}

func (r *ShadowRoot) InnerHTML() string {
	return r.v.Get("innerHTML").String()
}

func (r *ShadowRoot) SetInnerHTML(s string) {
	r.v.Set("innerHTML", s)
}
