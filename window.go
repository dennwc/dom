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

func (w *Window) JSRef() js.Ref {
	return w.v.JSRef()
}

func (w *Window) AddEventListener(typ string, h EventHandler) {
	w.v.Call("addEventListener", typ, js.NewEventCallback(func(v js.Value) {
		h(convertEvent(v))
	}))
}

func (w *Window) Open(url, windowName, windowFeatures string) {
	w.v.Call("open", url, windowName, windowFeatures)
}

func (w *Window) SetLocation(url string) {
	w.v.Set("location", url)
}

func (w *Window) OnResize(fnc func(e Event)) {
	w.AddEventListener("resize", fnc)
}
