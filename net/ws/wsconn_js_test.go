// +build js

package ws

import (
	"net/rpc"
	"strings"
	"testing"

	"github.com/dennwc/dom/js"
	"github.com/stretchr/testify/require"
)

func TestWebSocketsJS(t *testing.T) {
	url := js.Get("window", "location").String()
	url = "ws://" + strings.TrimPrefix(url, "http://") + "/ws"
	c, err := Dial(url)
	require.NoError(t, err)
	defer c.Close()

	cli := rpc.NewClient(c)

	var out string
	err = cli.Call("S.Hello", "Alice", &out)
	require.NoError(t, err)
	require.Equal(t, "Hello Alice", out)
}

func TestWebSocketsJSFail(t *testing.T) {
	_, err := Dial("ws://localhost:80")
	require.NotNil(t, err)
	require.Equal(t, "ws.dial: connection closed with code 1006", err.Error())
}
