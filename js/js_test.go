//+build wasm,js

package js

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyErrorObj(t *testing.T) {
	e := Error{NewObject().Ref}
	require.Equal(t, "JavaScript error: undefined", e.Error())
}
