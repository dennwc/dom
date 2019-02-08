package dom

import (
	"strings"

	"github.com/dennwc/dom/js"
)

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

func (w *Window) Open(url, windowName, windowFeatures map[string]string) {
	w.v.Call("open", url, windowName, joinKeyValuePairs(windowFeatures))
}

func (w *Window) SetLocation(url string) {
	w.v.Set("location", url)
}

func (w *Window) OnResize(fnc func(e Event)) {
	w.AddEventListener("resize", fnc)
}

func joinKeyValuePairs(kvpair map[string]string, joiner string) (ret string) {
	if kvpair == nil {
		return ""
	}

	for k, v := range kvpair {
		ret += k + "=" + v + joiner
	}
	ret = strings.TrimSuffix(ret, joiner)
	return
}
