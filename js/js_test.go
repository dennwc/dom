//+build wasm,js

package js

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUndefined(t *testing.T) {
	var v Value
	require.True(t, v.IsUndefined())
}

func TestEmptyErrorObj(t *testing.T) {
	e := Error{NewObject().Ref}
	require.Equal(t, "JavaScript error: undefined", e.Error())
}
