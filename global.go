package dom

import "github.com/dennwc/dom/js"

var (
	Doc  = GetDocument()
	Body = Doc.GetElementsByTagName("body")[0]
)

type Value = js.JSRef

func ConsoleLog(args ...interface{}) {
	js.Get("console").Call("log", args...)
}
