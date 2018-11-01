//+build wasm,js

package dom

import (
	"fmt"

	"github.com/dennwc/dom/js"
)

type EventTarget interface {
	js.Wrapper
	AddEventListener(typ string, h EventHandler)
	// TODO: removeEventListener
	// TODO: dispatchEvent
}

type Event interface {
	js.Wrapper
	Bubbles() bool
	Cancelable() bool
	Composed() bool
	CurrentTarget() *Element
	DefaultPrevented() bool
	Target() *Element
	Type() string
	IsTrusted() bool
	Path() NodeList

	PreventDefault()
	StopPropagation()
	StopImmediatePropagation()
}

type EventHandler func(Event)

type EventConstructor func(e BaseEvent) Event

func RegisterEventType(typ string, fnc EventConstructor) {
	cl := js.Get(typ)
	if !cl.Valid() {
		panic(fmt.Errorf("class undefined: %q", typ))
	}
	eventClasses = append(eventClasses, eventClass{
		Class: cl, New: fnc,
	})
}

func init() {
	RegisterEventType("MouseEvent", func(e BaseEvent) Event {
		return &MouseEvent{e}
	})
}

type eventClass struct {
	Class js.Value
	New   EventConstructor
}

var (
	eventClasses []eventClass
)

func convertEvent(v js.Value) Event {
	e := BaseEvent{v: v}
	// TODO: get class name directly
	for _, cl := range eventClasses {
		if v.InstanceOf(cl.Class) {
			return cl.New(e)
		}
	}
	return &e
}

type BaseEvent struct {
	v js.Value
}

func (e *BaseEvent) getBool(name string) bool {
	return e.v.Get(name).Bool()
}
func (e *BaseEvent) Bubbles() bool {
	return e.getBool("bubbles")
}

func (e *BaseEvent) Cancelable() bool {
	return e.getBool("cancelable")
}

func (e *BaseEvent) Composed() bool {
	return e.getBool("composed")
}

func (e *BaseEvent) CurrentTarget() *Element {
	return AsElement(e.v.Get("currentTarget"))
}

func (e *BaseEvent) DefaultPrevented() bool {
	return e.getBool("defaultPrevented")
}

func (e *BaseEvent) IsTrusted() bool {
	return e.getBool("isTrusted")
}

func (e *BaseEvent) JSValue() js.Ref {
	return e.v.JSValue()
}

func (e *BaseEvent) Type() string {
	return e.v.Get("type").String()
}

func (e *BaseEvent) Target() *Element {
	return AsElement(e.v.Get("target"))
}

func (e *BaseEvent) Path() NodeList {
	return AsNodeList(e.v.Get("path"))
}

func (e *BaseEvent) PreventDefault() {
	e.v.Call("preventDefault")
}
func (e *BaseEvent) StopPropagation() {
	e.v.Call("stopPropagation")
}
func (e *BaseEvent) StopImmediatePropagation() {
	e.v.Call("stopImmediatePropagation")
}

type MouseEventHandler func(*MouseEvent)

type MouseEvent struct {
	BaseEvent
}

func (e *MouseEvent) getPos(nameX, nameY string) Point {
	x := e.v.Get(nameX).Int()
	y := e.v.Get(nameY).Int()
	return Point{X: x, Y: y}
}

func (e *MouseEvent) getPosPref(pref string) Point {
	return e.getPos(pref+"X", pref+"Y")
}

const (
	MouseLeft = MouseButton(0)
)

type MouseButton int

func (e *MouseEvent) Button() MouseButton {
	return MouseButton(e.v.Get("button").Int())
}

func (e *MouseEvent) ClientPos() Point {
	return e.getPosPref("client")
}

func (e *MouseEvent) OffsetPos() Point {
	return e.getPosPref("offset")
}

func (e *MouseEvent) PagePos() Point {
	return e.getPosPref("page")
}

func (e *MouseEvent) ScreenPos() Point {
	return e.getPosPref("screen")
}

func (e *MouseEvent) AltKey() bool {
	return e.v.Get("altKey").Bool()
}

func (e *MouseEvent) CtrlKey() bool {
	return e.v.Get("ctrlKey").Bool()
}

func (e *MouseEvent) ShiftKey() bool {
	return e.v.Get("shiftKey").Bool()
}

func (e *MouseEvent) MetaKey() bool {
	return e.v.Get("metaKey").Bool()
}
