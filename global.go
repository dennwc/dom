package dom

import (
	"image"

	"github.com/dennwc/dom/js"
)

var (
	Doc  = GetDocument()
	Body = getFirstWithTag("body")
	Head = getFirstWithTag("head")
)

func getFirstWithTag(tag string) *HTMLElement {
	list := Doc.GetElementsByTagName(tag)
	if len(list) == 0 {
		return nil
	}
	return list[0].AsHTMLElement()
}

// Value is an alias for js.Wrapper.
//
// Derprecated: use js.Wrapper
type Value = js.Wrapper

func ConsoleLog(args ...interface{}) {
	js.Get("console").Call("log", args...)
}

func Loop() {
	select {}
}

type Point = image.Point
type Rect = image.Rectangle
