package dom

import "syscall/js"

func GetWindow() *Window {
	win := global.Get("window")
	if !IsValid(win) {
		return nil
	}
	return &Window{v: win}
}

var _ EventTarget = (*Window)(nil)

type Window struct {
	v js.Value
}

func (w *Window) JSValue() js.Value {
	return w.v
}

func (w *Window) AddEventListener(typ string, fnc func(e Event)) {
	w.v.Call("addEventListener", typ, js.NewEventCallback(0, func(v js.Value) {
		fnc(convertEvent(v))
	}))
}

func (w *Window) OnResize(fnc func(e Event)) {
	w.AddEventListener("resize", fnc)
}
