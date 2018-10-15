package dom

import (
	"strings"

	"github.com/dennwc/dom/js"
)

var required = make(map[string]error)

func appendAndWait(e *Element) error {
	errc := make(chan error, 1)
	e.AddEventListener("error", func(e Event) {
		errc <- js.Error{Value: js.Value{Ref: e.JSRef()}}
	})
	done := make(chan struct{})
	e.AddEventListener("load", func(e Event) {
		close(done) // TODO: load may happen before an error
	})
	Head.AppendChild(e)
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
	var s *Element
	if strings.HasSuffix(path, ".css") {
		// stylesheet
		s = NewElement("link")
		s.v.Set("type", "text/css")
		s.v.Set("rel", "stylesheet")
		s.v.Set("href", path)
	} else {
		// script
		s = NewElement("script")
		s.v.Set("async", true)
		s.v.Set("src", path)
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
