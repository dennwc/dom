package dom

import "syscall/js"

type Node interface {
	JSValue() js.Value
}
