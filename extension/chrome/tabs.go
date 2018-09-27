//+build js

package chrome

import "github.com/dennwc/dom/js"

type AllTabs interface {
	js.JSRef
	GetCurrent() Tab
	GetSelected(window WindowID) Tab
	GetAllInWindow(window WindowID) []Tab
}

func Tabs() AllTabs {
	return jsTabs{v: chrome.Get("tabs")}
}

type jsTabs struct {
	v js.Value
}

func (t jsTabs) JSRef() js.Ref {
	return t.v.JSRef()
}

func (t jsTabs) callAsync(name string, args ...interface{}) js.Value {
	ch := make(chan js.Value, 1)
	cb := js.NewEventCallback(func(v js.Value) {
		ch <- v
	})
	defer cb.Release()
	args = append(args, cb)
	t.v.Call(name, args...)
	return <-ch
}
func (t jsTabs) GetCurrent() Tab {
	v := t.callAsync("getCurrent")
	if !v.Valid() {
		return nil
	}
	return jsTab{v}
}
func (t jsTabs) GetSelected(window WindowID) Tab {
	var win interface{}
	if window != 0 {
		win = int(window)
	}
	v := t.callAsync("getSelected", win)
	if !v.Valid() {
		return nil
	}
	return jsTab{v}
}
func (t jsTabs) GetAllInWindow(window WindowID) []Tab {
	var win interface{}
	if window != 0 {
		win = int(window)
	}
	v := t.callAsync("getAllInWindow", win)
	vals := v.Slice()
	tabs := make([]Tab, 0, len(vals))
	for _, v := range vals {
		tabs = append(tabs, jsTab{v})
	}
	return tabs
}

type Tab interface {
	js.JSRef
	ID() int
	Active() bool
	Incognito() bool
	Highlighted() bool
	Pinned() bool
	Selected() bool
	Index() int
	WindowID() int
	URL() string
	Title() string
	Size() (w, h int)
}

type jsTab struct {
	v js.Value
}

func (t jsTab) JSRef() js.Ref {
	return t.v.JSRef()
}

func (t jsTab) ID() int {
	return t.v.Get("id").Int()
}

func (t jsTab) Active() bool {
	return t.v.Get("active").Bool()
}

func (t jsTab) Incognito() bool {
	return t.v.Get("incognito").Bool()
}

func (t jsTab) Highlighted() bool {
	return t.v.Get("highlighted").Bool()
}

func (t jsTab) Pinned() bool {
	return t.v.Get("pinned").Bool()
}

func (t jsTab) Selected() bool {
	return t.v.Get("selected").Bool()
}

func (t jsTab) WindowID() int {
	return t.v.Get("windowId").Int()
}

func (t jsTab) Index() int {
	return t.v.Get("index").Int()
}

func (t jsTab) Title() string {
	return t.v.Get("title").String()
}

func (t jsTab) URL() string {
	return t.v.Get("url").String()
}

func (t jsTab) Size() (w, h int) {
	w = t.v.Get("width").Int()
	h = t.v.Get("height").Int()
	return
}
