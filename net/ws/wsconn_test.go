// +build !js

package ws

import (
	"github.com/dennwc/dom/js/jstest"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/rpc"
	"testing"
)

type service struct{}

func (service) Hello(name string, out *string) error {
	*out = "Hello " + name
	return nil
}

func TestWebSockets(t *testing.T) {
	errc := make(chan error, 1)
	rs := rpc.NewServer()
	err := rs.RegisterName("S", service{})
	require.NoError(t, err)

	srv := newServer(nil)
	defer srv.Close()

	lis := srv

	go func() {
		for {
			c, err := lis.Accept()
			if err != nil {
				errc <- err
				return
			}
			go func() {
				defer c.Close()
				rs.ServeConn(c)
			}()
		}
	}()

	jstest.RunTestChrome(t, http.HandlerFunc(srv.handleWS))
	select {
	case err = <-errc:
		require.NoError(t, err)
	default:
	}
}
