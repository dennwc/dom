package dom

import "syscall/js"

func GetDocument() *Document {
	doc := js.Global().Get("document")
	if doc == js.Null() || doc == js.Undefined() {
		panic("no document")
	}
	return &Document{v: doc}
}

type Document struct {
	v js.Value
}

func (d *Document) CreateElement(tag string) *Element {
	v := d.v.Call("createElement", tag)
	return &Element{v: v}
}
func (d *Document) CreateElementNS(ns string, tag string) *Element {
	v := d.v.Call("createElementNS", ns, tag)
	return &Element{v: v}
}
func (d *Document) GetElementsByTagName(tag string) []*Element {
	v := d.v.Call("getElementsByTagName", tag)
	arr := make([]*Element, v.Length())
	for i := range arr {
		arr[i] = &Element{v: v.Index(i)}
	}
	return arr
}
