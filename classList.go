package dom

import "github.com/dennwc/dom/js"

type ClassList struct {
	v js.Value
}

func (c *ClassList) Add(class ...interface{}) {
	c.v.Call("add", class...)
}

func (c *ClassList) Remove(class ...interface{}) {
	c.v.Call("remove", class...)
}
