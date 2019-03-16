package dom

import (
	"fmt"

	"github.com/dennwc/dom/js"
)

var _ Node = (*Element)(nil)

func AsElement(v js.Value) *Element {
	if !v.Valid() {
		return nil
	}
	return &Element{NodeBase{v: v}}
}

func AsNodeList(v js.Value) NodeList {
	if !v.Valid() {
		return nil
	}
	arr := make(NodeList, v.Length())
	for i := range arr {
		arr[i] = AsElement(v.Index(i))
	}
	return arr
}

var _ Node = (*Element)(nil)

type Element struct {
	NodeBase
}

type Position string

const (
	BeforeBegin Position = "beforebegin"
	BeforeEnd   Position = "beforeend"
	AfterBegin  Position = "afterbegin"
	AfterEnd    Position = "afterend"
)

// Properties

// Attributes returns a NamedNodeMap object containing the assigned attributes of the corresponding HTML element.
// func (e *Element) Attributes() NamedNodeMap {
// 	return e.v.Get("attributes")
// }

// ClassList returns a DOMTokenList containing the list of class attributes.
func (e *Element) ClassList() *TokenList {
	return AsTokenList(e.v.Get("classList"))
}

// ClassName is a DOMString representing the class of the element.
func (e *Element) ClassName() string {
	return e.v.Get("className").String()
}

// SetClassName is a DOMString representing the class of the element.
func (e *Element) SetClassName(v string) {
	e.v.Set("className", v)
}

// ClientHeight returns a Number representing the inner height of the element.
func (e *Element) ClientHeight() int {
	return e.v.Get("clientHeight").Int()
}

// ClientLeft returns a Number representing the width of the left border of the element.
func (e *Element) ClientLeft() int {
	return e.v.Get("clientLeft").Int()
}

// ClientTop returns a Number representing the width of the top border of the element.
func (e *Element) ClientTop() int {
	return e.v.Get("clientTop").Int()
}

// ClientWidth returns a Number representing the inner width of the element.
func (e *Element) ClientWidth() int {
	return e.v.Get("clientWidth").Int()
}

// ComputedName returns a DOMString containing the label exposed to accessibility.
func (e *Element) ComputedName() string {
	return e.v.Get("computedName").String()
}

// ComputedRole returns a DOMString containing the ARIA role that has been applied to a particular element.
func (e *Element) ComputedRole() string {
	return e.v.Get("computedRole").String()
}

// Id is a DOMString representing the id of the element.
func (e *Element) Id() string {
	return e.v.Get("id").String()
}

// SetId is a DOMString representing the id of the element.
func (e *Element) SetId(v string) {
	e.v.Set("id", v)
}

// InnerHTML is a DOMString representing the markup of the element's content.
func (e *Element) InnerHTML() string {
	return e.v.Get("innerHTML").String()
}

// SetInnerHTML is a DOMString representing the markup of the element's content.
func (e *Element) SetInnerHTML(v string) {
	e.v.Set("innerHTML", v)
}

// LocalName a DOMString representing the local part of the qualified name of the element.
func (e *Element) LocalName() string {
	return e.v.Get("localName").String()
}

// NamespaceURI the namespace URI of the element, or null if it is no namespace.
func (e *Element) NamespaceURI() string {
	return e.v.Get("namespaceURI").String()
}

// OuterHTML is a DOMString representing the markup of the element including its content. When used as a setter, replaces the element with nodes parsed from the given string.
func (e *Element) OuterHTML() string {
	return e.v.Get("outerHTML").String()
}

// SetOuterHTML is a DOMString representing the markup of the element including its content. When used as a setter, replaces the element with nodes parsed from the given string.
func (e *Element) SetOuterHTML(v string) {
	e.v.Set("outerHTML", v)
}

// Prefix a DOMString representing the namespace prefix of the element, or null if no prefix is specified.
func (e *Element) Prefix() string {
	return e.v.Get("prefix").String()
}

// ScrollHeight returns a Number representing the scroll view height of an element.
func (e *Element) ScrollHeight() int {
	return e.v.Get("scrollHeight").Int()
}

// ScrollLeft is a Number representing the left scroll offset of the element.
func (e *Element) ScrollLeft() int {
	return e.v.Get("scrollLeft").Int()
}

// SetScrollLeft is a Number representing the left scroll offset of the element.
func (e *Element) SetScrollLeft(v int) {
	e.v.Set("scrollLeft", v)
}

// ScrollLeftMax returns a Number representing the maximum left scroll offset possible for the element.
func (e *Element) ScrollLeftMax() int {
	return e.v.Get("scrollLeftMax").Int()
}

// ScrollTop a Number representing number of pixels the top of the document is scrolled vertically.
func (e *Element) ScrollTop() int {
	return e.v.Get("scrollTop").Int()
}

// SetScrollTop a Number representing number of pixels the top of the document is scrolled vertically.
func (e *Element) SetScrollTop(v int) {
	e.v.Set("scrollTop", v)
}

// ScrollTopMax returns a Number representing the maximum top scroll offset possible for the element.
func (e *Element) ScrollTopMax() int {
	return e.v.Get("scrollTopMax").Int()
}

// ScrollWidth returns a Number representing the scroll view width of the element.
func (e *Element) ScrollWidth() int {
	return e.v.Get("scrollWidth").Int()
}

// Shadow returns the open shadow root that is hosted by the element, or null if no open shadow root is present.
func (e *Element) ShadowRoot() *ShadowRoot {
	return AsShadowRoot(e.v.Get("shadowRoot"))
}

// Slot  returns the name of the shadow DOM slot the element is inserted in.
func (e *Element) Slot() string {
	return e.v.Get("slot").String()
}

// SetSlot  returns the name of the shadow DOM slot the element is inserted in.
func (e *Element) SetSlot(v string) {
	e.v.Set("slot", v)
}

// TabStop  is a Boolean indicating if the element can receive input focus via the tab key.
func (e *Element) TabStop() bool {
	return e.v.Get("tabStop").Bool()
}

// SetTabStop  is a Boolean indicating if the element can receive input focus via the tab key.
func (e *Element) SetTabStop(v bool) {
	e.v.Set("tabStop", v)
}

// TagName returns a String with the name of the tag for the given element.
func (e *Element) TagName() string {
	return e.v.Get("tagName").String()
}

// UndoManager returns the UndoManager associated with the element.
func (e *Element) UndoManager() js.Value {
	return e.v.Get("undoManager")
}

// UndoScope  is a Boolean indicating if the element is an undo scope host, or not.
func (e *Element) UndoScope() bool {
	return e.v.Get("undoScope").Bool()
}

// SetUndoScope  is a Boolean indicating if the element is an undo scope host, or not.
func (e *Element) SetUndoScope(v bool) {
	e.v.Set("undoScope", v)
}

// Methods

func (e *Element) SetAttribute(k string, v interface{}) {
	e.v.Call("setAttribute", k, fmt.Sprint(v))
}

func (e *Element) GetAttribute(k string) js.Value {
	return e.v.Call("getAttribute", k)
}

func (e *Element) RemoveAttribute(k string) {
	e.v.Call("removeAttribute", k)
}

func (e *Element) GetBoundingClientRect() Rect {
	rv := e.v.Call("getBoundingClientRect")
	x, y := rv.Get("x").Int(), rv.Get("y").Int()
	w, h := rv.Get("width").Int(), rv.Get("height").Int()
	return Rect{Min: Point{x, y}, Max: Point{x + w, y + h}}
}

func (e *Element) onMouseEvent(typ string, h MouseEventHandler) {
	e.AddEventListener(typ, func(e Event) {
		h(e.(*MouseEvent))
	})
}

func (e *Element) OnClick(h MouseEventHandler) {
	e.onMouseEvent("click", h)
}

func (e *Element) OnMouseDown(h MouseEventHandler) {
	e.onMouseEvent("mousedown", h)
}

func (e *Element) OnMouseMove(h MouseEventHandler) {
	e.onMouseEvent("mousemove", h)
}

func (e *Element) OnMouseUp(h MouseEventHandler) {
	e.onMouseEvent("mouseup", h)
}

type AttachShadowOpts struct {
	Open           bool
	DeligatesFocus bool
}

func (e *Element) AttachShadow(opts AttachShadowOpts) *ShadowRoot {
	m := map[string]interface{}{}
	if opts.Open {
		m["mode"] = "open"
	} else {
		m["mode"] = "closed"
	}
	m["delegatesFocus"] = opts.DeligatesFocus
	return AsShadowRoot(e.v.Call("attachShadow", js.ValueOf(m)))
}

// InsertAdjacentElement inserts a given element node at a given position relative to the element it is invoked upon.
func (e *Element) InsertAdjacentElement(position Position, newElement *Element) js.Value {
	return e.v.Call("insertAdjacentElement", string(position), newElement.v)
}
