package js

// Promise is a function wrapping a JS Promise.
type Promise func() ([]Value, error)

// Await wait for the promise to be resolved or rejected.
// A shorthand for calling a promise function returned by Promised.
func (v Value) Await() ([]Value, error) {
	return v.Promised()()
}

// Promised returns a new function that awaits on the promise.
// Function will either return an slice of values, or an Error.
func (v Value) Promised() Promise {
	resolved := make(chan []Value, 1)
	rejected := make(chan []Value, 1)
	done := make(chan struct{}) // safeguard
	var (
		then, catch Callback
	)
	then = NewCallback(func(v []Value) {
		close(done)
		then.Release()
		catch.Release()
		resolved <- v
	})
	catch = NewCallback(func(v []Value) {
		close(done)
		then.Release()
		catch.Release()
		rejected <- v
	})
	v.Call("then", then).Call("catch", catch)
	return func() ([]Value, error) {
		select {
		case v := <-resolved:
			return v, nil
		case v := <-rejected:
			var e Value
			if len(v) != 0 {
				e = v[0]
			}
			return v, Error{Value: e}
		}
	}
}
