// +build js

package require

import (
	"testing"

	"github.com/dennwc/dom/js"
	"github.com/stretchr/testify/require"
)

func TestRequireJS(t *testing.T) {
	err := Require("/env.js")
	require.NoError(t, err)
	require.Equal(t, "ok", js.Get("Val").String())
}

func TestRequireJSSyntaxError(t *testing.T) {
	t.SkipNow() // FIXME
	err := Require("/err.js")
	require.NotNil(t, err)
}

func TestRequireJSNotFound(t *testing.T) {
	err := Require("/na.js")
	require.NotNil(t, err)
}
