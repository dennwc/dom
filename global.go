package dom

import (
	"image"

	"github.com/dennwc/dom/js"
)

var (
	Doc  = GetDocument()
	Body = Doc.GetElementsByTagName("body")[0]
	Head = Doc.GetElementsByTagName("head")[0]
)

type Value = js.JSRef

func ConsoleLog(args ...interface{}) {
	js.Get("console").Call("log", args...)
}

func Loop() {
	select {}
}

type Point = image.Point
type Rect = image.Rectangle
