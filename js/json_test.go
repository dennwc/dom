//+build wasm,js

package js

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValueMarshalJSON(t *testing.T) {
	t.Run("object", func(t *testing.T) {
		v := NewObject()
		v.Set("k", "v")
		require.True(t, v.Valid())

		data, err := v.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `{"k":"v"}`, string(data))
	})
	t.Run("undefined", func(t *testing.T) {
		var v Value
		data, err := v.MarshalJSON()
		require.NoError(t, err)
		require.Equal(t, `null`, string(data))
	})
}

func TestValueUnmarshalJSON(t *testing.T) {
	var v Value
	err := json.Unmarshal([]byte(`{"k":"v"}`), &v)
	require.NoError(t, err)

	require.Equal(t, "v", v.Get("k").String())
}

func TestValueUnmarshalJSONError(t *testing.T) {
	var v Value
	// Call unmarshal directly, because json.Unmarshal will try to validate an input.
	err := v.UnmarshalJSON([]byte(`{"k":"v"`))
	require.NotNil(t, err)
}
