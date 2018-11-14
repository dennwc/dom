//+build wasm,js,go1.12

package js

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUndefined(t *testing.T) {
	var v Value
	require.True(t, v.IsUndefined())
}
