package dom

import "syscall/js"

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
	v js.Value
}

func (e *NodeBase) JSValue() js.Value {
	return e.v
}

func (e *NodeBase) AddEventListener(typ string, h EventHandler) {
	e.v.Call("addEventListener", typ, js.NewEventCallback(0, func(v js.Value) {
		h(convertEvent(v))
	}))
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
	e.v.Call("appendChild", n.JSValue())
}

func (e *NodeBase) Contains(n Node) bool {
	return e.v.Call("contains", n.JSValue()).Bool()
}

func (e *NodeBase) IsEqualNode(n Node) bool {
	return e.v.Call("isEqualNode", n.JSValue()).Bool()
}

func (e *NodeBase) IsSameNode(n Node) bool {
	return e.v.Call("isSameNode", n.JSValue()).Bool()
}

func (e *NodeBase) RemoveChild(n Node) Node {
	return AsElement(e.v.Call("removeChild", n.JSValue()))
}

func (e *NodeBase) ReplaceChild(n, old Node) Node {
	return AsElement(e.v.Call("replaceChild", n.JSValue(), old.JSValue()))
}
