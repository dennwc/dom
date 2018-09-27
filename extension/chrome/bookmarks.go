//+build js

package chrome

import "github.com/dennwc/dom/js"

type AllBookmarks interface {
	GetTree() []BookmarkNode
}

func Bookmarks() AllBookmarks {
	v := chrome.Get("bookmarks")
	if !v.Valid() {
		return nil
	}
	return jsBookmarks{v}
}

type BookmarkNode interface {
	js.JSRef
	ID() string
	ParentID() string
	Index() int
	URL() string
	Title() string
	Children() []BookmarkNode
}

type jsBookmarks struct {
	v js.Value
}

func (b jsBookmarks) callAsync(name string, args ...interface{}) js.Value {
	ch := make(chan js.Value, 1)
	cb := js.NewEventCallback(func(v js.Value) {
		ch <- v
	})
	defer cb.Release()
	args = append(args, cb)
	b.v.Call(name, args...)
	return <-ch
}
func (b jsBookmarks) GetTree() []BookmarkNode {
	arr := b.callAsync("getTree").Slice()
	nodes := make([]BookmarkNode, 0, len(arr))
	for _, v := range arr {
		nodes = append(nodes, jsBookmarkNode{v})
	}
	return nodes
}

type jsBookmarkNode struct {
	v js.Value
}

func (b jsBookmarkNode) JSRef() js.Ref {
	return b.v.JSRef()
}

func (b jsBookmarkNode) ID() string {
	return b.v.Get("id").String()
}

func (b jsBookmarkNode) ParentID() string {
	return b.v.Get("id").String()
}

func (b jsBookmarkNode) Index() int {
	return b.v.Get("index").Int()
}

func (b jsBookmarkNode) URL() string {
	return b.v.Get("url").String()
}

func (b jsBookmarkNode) Title() string {
	return b.v.Get("title").String()
}

func (b jsBookmarkNode) Children() []BookmarkNode {
	arr := b.v.Get("children").Slice()
	nodes := make([]BookmarkNode, 0, len(arr))
	for _, v := range arr {
		nodes = append(nodes, jsBookmarkNode{v})
	}
	return nodes
}
