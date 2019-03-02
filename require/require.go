//+build wasm,js

package require

import (
	"errors"
	"strings"
	"sync"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/js"
)

var required = make(map[string]error)

func appendAndWait(e *dom.Element) error {
	errc := make(chan error, 1)
	e.AddEventListener("error", func(e dom.Event) {
		errc <- js.NewError(e)
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

// Require adds a specified file (js or css) into the document and waits for it to load.
//
// The function relies on a file extension to detect the type. If there is no extension in
// the file path, use specific function like Stylesheet or Script. As an alternative,
// append a '#.js' or '#.css' suffix to a file path.
func Require(path string) error {
	if strings.HasSuffix(path, ".css") {
		return Stylesheet(path)
	} else if strings.HasSuffix(path, ".js") {
		return Script(path)
	}
	return errors.New("the file should have an extension specified (or '#.ext')")
}

// Stylesheet add a specified CSS file into the document and waits for it to load.
func Stylesheet(path string) error {
	if err, ok := required[path]; ok {
		return err
	}
	s := dom.NewElement("link")
	v := s.JSValue()
	v.Set("type", "text/css")
	v.Set("rel", "stylesheet")
	v.Set("href", path)
	err := appendAndWait(s)
	required[path] = err
	return err
}

// Script adds a specified JS file into the document and waits for it to load.
func Script(path string) error {
	if err, ok := required[path]; ok {
		return err
	}
	s := dom.NewElement("script")
	v := s.JSValue()
	v.Set("async", true)
	v.Set("src", path)
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

// RequireLazy is the same as Require, but returns a function that will load the file on the first call.
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

// StylesheetString loads a CSS stylesheet string into the DOM.
func StylesheetString(data string) {
	s := dom.NewElement("style")
	s.JSValue().Set("type", "text/css")
	s.SetInnerHTML(data)
	err := appendAndWait(s)
	if err != nil {
		panic(err)
	}
}
