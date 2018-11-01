//+build wasm,js

package require

import (
	"strings"
	"sync"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/js"
)

var required = make(map[string]error)

func appendAndWait(e *dom.Element) error {
	errc := make(chan error, 1)
	e.AddEventListener("error", func(e dom.Event) {
		errc <- js.Error{Value: js.Value{Ref: e.JSValue()}}
	})
	done := make(chan struct{})
	e.AddEventListener("load", func(e dom.Event) {
		close(done) // TODO: load may happen before an error
	})
	dom.Head.AppendChild(e)
	select {
	case err := <-errc:
		return err
	case <-done:
	}
	return nil
}

// Require loads a specified file (js or css) into the document and waits for it to apply.
func Require(path string) error {
	if err, ok := required[path]; ok {
		return err
	}
	var s *dom.Element
	if strings.HasSuffix(path, ".css") {
		// stylesheet
		s = dom.NewElement("link")
		v := s.JSValue()
		v.Set("type", "text/css")
		v.Set("rel", "stylesheet")
		v.Set("href", path)
	} else {
		// script
		s = dom.NewElement("script")
		v := s.JSValue()
		v.Set("async", true)
		v.Set("src", path)
	}
	err := appendAndWait(s)
	required[path] = err
	return err
}

// MustRequire is the same as Require, but panics on an error.
func MustRequire(path string) {
	err := Require(path)
	if err != nil {
		panic(err)
	}
}

// MustRequireValue loads a specified file and returns a global value with a given name.
func MustRequireValue(name, path string) js.Value {
	MustRequire(path)
	return js.Get(name)
}

func RequireLazy(path string) func() error {
	var (
		once sync.Once
		err  error
	)
	return func() error {
		once.Do(func() {
			err = Require(path)
		})
		return err
	}
}

func StylesheetString(data string) {
	s := dom.NewElement("style")
	s.JSValue().Set("type", "text/css")
	s.SetInnerHTML(data)
	err := appendAndWait(s)
	if err != nil {
		panic(err)
	}
}
