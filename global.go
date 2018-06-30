package dom

import "syscall/js"

var (
	global    = js.Global()
	null      = js.Null()
	undefined = js.Undefined()
)

var (
	Doc  = GetDocument()
	Body = Doc.GetElementsByTagName("body")[0]
)

type Value interface {
	JSValue() js.Value
}

func IsValid(v js.Value) bool {
	return v != null && v != undefined
}

func ConsoleLog(args ...interface{}) {
	for i, a := range args {
		if v, ok := a.(Value); ok {
			args[i] = v.JSValue()
		}
	}
	global.Get("console").Call("log", args...)
}
