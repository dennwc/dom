package chrome

import "github.com/dennwc/dom/js"

type TabOptions struct {
	Active bool
}

type AllTabs interface {
	js.Wrapper
	GetCurrent() Tab
	GetSelected(window WindowID) Tab
	GetAllInWindow(window WindowID) []Tab
	Create(url string, opt *TabOptions) Tab
}

func Tabs() AllTabs {
	return tabs()
}

func tabs() jsTabs {
	return jsTabs{v: chrome.Get("tabs")}
}

type jsTabs struct {
	v js.Value
}

func (t jsTabs) JSValue() js.Ref {
	return t.v.JSValue()
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

func (t jsTabs) Create(url string, opt *TabOptions) Tab {
	obj := js.Obj{
		"url": url,
	}
	if opt != nil {
		obj["active"] = opt.Active
	}
	v := t.callAsync("create", obj)
	return jsTab{v}
}

type Tab interface {
	js.Wrapper
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

	ExecuteFile(path string) (js.Value, error)
	ExecuteCode(code string) (js.Value, error)
}

func AsTab(v js.Value) Tab {
	return jsTab{v: v}
}

type jsTab struct {
	v js.Value
}

func (t jsTab) JSValue() js.Ref {
	return t.v.JSValue()
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

func (t jsTab) executeScript(obj interface{}) ([]js.Value, error) {
	id := t.ID()
	ch := make(chan []js.Value, 1)
	errc := make(chan error, 1)
	cb := js.CallbackOf(func(args []js.Value) {
		err := lastError()
		if err != nil {
			errc <- err
			return
		}
		v := args[0]
		if !v.Valid() {
			ch <- nil
		} else {
			ch <- v.Slice()
		}
	})
	defer cb.Release()
	tabs().v.Call("executeScript", id, obj, cb)
	select {
	case err := <-errc:
		return nil, err
	case res := <-ch:
		return res, nil
	}
}

func (t jsTab) executeScriptOne(obj interface{}) (js.Value, error) {
	res, err := t.executeScript(obj)
	if err != nil || len(res) == 0 {
		return js.Value{}, err
	}
	return res[0], nil
}

func (t jsTab) ExecuteFile(path string) (js.Value, error) {
	return t.executeScriptOne(js.Obj{"file": path})
}

func (t jsTab) ExecuteCode(code string) (js.Value, error) {
	return t.executeScriptOne(js.Obj{"code": code})
}
