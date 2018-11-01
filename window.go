//+build wasm,js

package dom

import "github.com/dennwc/dom/js"

func GetWindow() *Window {
	win := js.Get("window")
	if !win.Valid() {
		return nil
	}
	return &Window{v: win}
}

var _ EventTarget = (*Window)(nil)

type Window struct {
	v js.Value
}

func (w *Window) JSValue() js.Ref {
	return w.v.JSValue()
}

func (w *Window) AddEventListener(typ string, h EventHandler) {
	w.v.Call("addEventListener", typ, js.NewEventCallback(func(v js.Value) {
		h(convertEvent(v))
	}))
}

func (w *Window) OnResize(fnc func(e Event)) {
	w.AddEventListener("resize", fnc)
}
