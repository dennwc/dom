//+build wasm,js

package dom

import "github.com/dennwc/dom/js"

type Node interface {
	EventTarget

	// properties

	BaseURI() string
	NodeName() string
	ChildNodes() NodeList
	ParentNode() Node
	ParentElement() *Element
	TextContent() string
	SetTextContent(s string)

	// methods

	AppendChild(n Node)
	Contains(n Node) bool
	IsEqualNode(n Node) bool
	IsSameNode(n Node) bool
	RemoveChild(n Node) Node
	ReplaceChild(n, old Node) Node
}

type NodeList []*Element

type NodeBase struct {
	v         js.Value
	callbacks []js.Callback
}

func (e *NodeBase) JSValue() js.Ref {
	return e.v.JSValue()
}

func (e *NodeBase) Remove() {
	e.ParentNode().RemoveChild(e)
	for _, c := range e.callbacks {
		c.Release()
	}
	e.callbacks = nil
}

func (e *NodeBase) AddEventListener(typ string, h EventHandler) {
	cb := js.NewEventCallback(func(v js.Value) {
		h(convertEvent(v))
	})
	e.callbacks = append(e.callbacks, cb)
	e.v.Call("addEventListener", typ, cb)
}

func (e *NodeBase) AddErrorListener(h func(err error)) {
	e.AddEventListener("error", func(e Event) {
		ConsoleLog(e.JSValue())
		h(js.Error{Value: js.Value{Ref: e.JSValue()}})
	})
}

func (e *NodeBase) BaseURI() string {
	return e.v.Get("baseURI").String()
}

func (e *NodeBase) NodeName() string {
	return e.v.Get("nodeName").String()
}

func (e *NodeBase) ChildNodes() NodeList {
	return AsNodeList(e.v.Get("childNodes"))
}

func (e *NodeBase) ParentNode() Node {
	return AsElement(e.v.Get("parentNode"))
}

func (e *NodeBase) ParentElement() *Element {
	return AsElement(e.v.Get("parentElement"))
}

func (e *NodeBase) TextContent() string {
	return e.v.Get("textContent").String()
}

func (e *NodeBase) SetTextContent(s string) {
	e.v.Set("textContent", s)
}

func (e *NodeBase) AppendChild(n Node) {
	e.v.Call("appendChild", n)
}

func (e *NodeBase) Contains(n Node) bool {
	return e.v.Call("contains", n).Bool()
}

func (e *NodeBase) IsEqualNode(n Node) bool {
	return e.v.Call("isEqualNode", n).Bool()
}

func (e *NodeBase) IsSameNode(n Node) bool {
	return e.v.Call("isSameNode", n).Bool()
}

func (e *NodeBase) RemoveChild(n Node) Node {
	return AsElement(e.v.Call("removeChild", n))
}

func (e *NodeBase) ReplaceChild(n, old Node) Node {
	return AsElement(e.v.Call("replaceChild", n, old))
}
