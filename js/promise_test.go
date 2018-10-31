//+build wasm,js

package js

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPromiseAwait(t *testing.T) {
	f1 := NewFuncJS(`return new Promise((resolve) => {
	setTimeout(resolve, 0)
})`)

	ch := make(chan []Value, 1)
	errc := make(chan error, 1)
	p := f1.Invoke().Promised()
	go func() {
		r, err := p.Await()
		if err != nil {
			errc <- err
		} else {
			ch <- r
		}
	}()

	select {
	case <-ch:
	case err := <-errc:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("deadlock")
	}
}

func TestPromiseAwaitReject(t *testing.T) {
	f1 := NewFuncJS(`return new Promise((resolve, reject) => {
	setTimeout(reject, 0)
})`)

	ch := make(chan []Value, 1)
	errc := make(chan error, 1)
	p := f1.Invoke().Promised()
	go func() {
		r, err := p.Await()
		if err != nil {
			errc <- err
		} else {
			ch <- r
		}
	}()

	select {
	case <-ch:
		t.Fatal("expected promise to be rejected")
	case err := <-errc:
		require.NotNil(t, err)
	case <-time.After(time.Second):
		t.Fatal("deadlock")
	}
}

func TestNewPromiseResolve(t *testing.T) {
	called := make(chan struct{})
	v := NewPromise(func() ([]interface{}, error) {
		close(called)
		return nil, nil
	})

	ch := make(chan []Value, 1)
	errc := make(chan error, 1)
	p := v.Promised()
	go func() {
		r, err := p.Await()
		if err != nil {
			errc <- err
		} else {
			ch <- r
		}
	}()

	select {
	case <-ch:
		select {
		case <-called:
		default:
			t.Fatal("function was not called")
		}
	case err := <-errc:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("deadlock")
	}
}

func TestNewPromiseReject(t *testing.T) {
	called := make(chan struct{})
	v := NewPromise(func() ([]interface{}, error) {
		close(called)
		return nil, errors.New("err")
	})

	ch := make(chan []Value, 1)
	errc := make(chan error, 1)
	p := v.Promised()
	go func() {
		r, err := p.Await()
		if err != nil {
			errc <- err
		} else {
			ch <- r
		}
	}()

	select {
	case <-ch:
		t.Fatal("expected promise to be rejected")
	case err := <-errc:
		require.NotNil(t, err)
		select {
		case <-called:
		default:
			t.Fatal("function was not called")
		}
	case <-time.After(time.Second):
		t.Fatal("deadlock")
	}
}
