//+build wasm,js

package js

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFuncJS(t *testing.T) {
	v := NewFuncJS("a", "b", `return a+b`)
	c := int(v.Invoke(1, 2).Int())
	require.Equal(t, 3, c)
}
