//+build wasm

package js

type Error struct {
	Value
}

func (e Error) Error() string {
	return "error: " + e.String()
}
