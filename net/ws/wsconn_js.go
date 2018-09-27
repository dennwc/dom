//+build js

package ws

import (
	"bytes"
	"io"
	"net"
	"time"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/js"
)

// Dial connects to a WebSocket on a specified URL.
func Dial(addr string) (net.Conn, error) {
	c := &jsConn{
		open: make(chan struct{}),
		done: make(chan struct{}),
		errc: make(chan error, 1),
		msg:  make(chan js.Value, 1),
	}
	c.ws = js.New("WebSocket", addr)
	c.installCallbacks()
	select {
	case err := <-c.errc:
		return nil, err
	case <-c.open:
	}
	return c, nil
}

type jsConn struct {
	ws  js.Value
	cbs *js.CallbackGroup

	open chan struct{}
	errc chan error
	done chan struct{}
	err  error

	msg chan js.Value

	rbuf bytes.Buffer
}

func (c *jsConn) installCallbacks() {
	c.cbs = c.ws.NewCallbackGroup()
	c.cbs.Set("onopen", c.onOpen)
	c.cbs.Set("onerror", c.onError)
	c.cbs.Set("onclose", c.onClose)
	msg := js.NewFunction("onMsg", "onErr", `
return function(e){
	const r = new FileReader();
	r.addEventListener('loadend', function(){
		onMsg(new Uint8Array(r.result));
	})
	r.addEventListener('onerror', onErr);
	r.readAsArrayBuffer(e.data);
}
`)
	cb := js.NewCallback(c.onMessage)
	c.cbs.Add(cb)
	c.ws.Set("onmessage", msg.Invoke(cb, cb))
}

func (c *jsConn) onError(args []js.Value) {
	err := js.Error{Value: args[0]}
	select {
	case c.errc <- err:
	default:
		c.err = err
	}
	c.Close()
}

func (c *jsConn) onMessage(args []js.Value) {
	e := args[0]
	if e.InstanceOfClass("Uint8Array") {
		select {
		case c.msg <- e:
		case <-c.done:
		}
		return
	}
	dom.ConsoleLog("read error:", e)
	c.onError([]js.Value{e})
}

func (c *jsConn) onOpen(_ []js.Value) {
	close(c.open)
}

func (c *jsConn) onClose(args []js.Value) {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
}

func (c *jsConn) send(data []byte) {
	arr := js.TypedArrayOf(data)
	jarr := js.New("Uint8Array", arr)
	arr.Release()
	c.ws.Call("send", jarr)
}

func (c *jsConn) close() {
	c.ws.Call("close")
}

func (c *jsConn) Read(b []byte) (int, error) {
	for {
		if c.rbuf.Len() != 0 {
			return c.rbuf.Read(b)
		} else if c.err != nil {
			return 0, c.err
		}
		select {
		case err := <-c.errc:
			c.err = err
			return 0, c.err
		case <-c.done:
			return 0, io.EOF
		case arr := <-c.msg:
			sz := arr.Get("length").Int()

			data := make([]byte, sz)
			m := js.MMap(data)
			err := m.CopyFrom(arr)
			m.Release()
			if err != nil {
				c.err = err
				return 0, err
			}
			c.rbuf.Write(data)
		}
	}
}

func (c *jsConn) Write(b []byte) (int, error) {
	if c.err != nil {
		return 0, c.err
	}
	select {
	case err := <-c.errc:
		c.err = err
		return 0, c.err
	default:
	}
	c.send(b)
	return len(b), nil
}

func (c *jsConn) LocalAddr() net.Addr {
	return wsAddr{}
}

func (c *jsConn) RemoteAddr() net.Addr {
	return wsAddr{}
}

func (c *jsConn) SetDeadline(t time.Time) error {
	return nil // TODO
}

func (c *jsConn) SetReadDeadline(t time.Time) error {
	return nil // TODO
}

func (c *jsConn) SetWriteDeadline(t time.Time) error {
	return nil // TODO
}

func (c *jsConn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	c.close()
	c.cbs.Release()
	return c.err
}
