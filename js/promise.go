//+build wasm,js

package js

import (
	"context"
)

var (
	_ Wrapper = (*Promise)(nil)
)

// NewPromise runs a given function asynchronously by converting it to JavaScript promise.
// Promise will be resolved if functions returns a nil error and will be rejected otherwise.
func NewPromise(fnc func() ([]interface{}, error)) Value {
	var initFunc Callback
	initFunc = NewCallbackAsync(func(args []Value) {
		initFunc.Release()
		resolve, reject := args[0], args[1]
		res, err := fnc()
		if err != nil {
			if w, ok := err.(Wrapper); ok {
				reject.Invoke(w)
			} else {
				reject.Invoke(err.Error())
			}
		} else {
			resolve.Invoke(res...)
		}
	})
	return New("Promise", initFunc)
}

// Promise represents a JavaScript Promise.
type Promise struct {
	v    Value
	done <-chan struct{}
	res  []Value
	err  error
}

// JSValue implements Wrapper interface.
func (p *Promise) JSValue() Ref {
	return p.v.JSValue()
}

// Await for the promise to resolve.
func (p *Promise) Await() ([]Value, error) {
	<-p.done
	return p.res, p.err
}

// AwaitContext for the promise to resolve or context to be canceled.
func (p *Promise) AwaitContext(ctx context.Context) ([]Value, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.done:
		return p.res, p.err
	}
}

// Promised returns converts the value into a Promise.
func (v Value) Promised() *Promise {
	done := make(chan struct{})
	p := &Promise{
		v: v, done: done,
	}
	var then, catch Callback
	then = NewCallback(func(v []Value) {
		then.Release()
		catch.Release()
		p.res = v
		close(done)
	})
	catch = NewCallback(func(v []Value) {
		then.Release()
		catch.Release()
		var e Value
		if len(v) != 0 {
			e = v[0]
		}
		if e.Ref == undefined {
			e = NewObject()
		}
		p.err = Error{e.Ref}
		close(done)
	})
	v.Call("then", then).Call("catch", catch)
	return p
}

// Await wait for the promise to be resolved or rejected.
// A shorthand for calling Await on the promise returned by Promised.
func (v Value) Await() ([]Value, error) {
	return v.Promised().Await()
}
