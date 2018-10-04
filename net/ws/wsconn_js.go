//+build js

package ws

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/dennwc/dom/js"
)

var errClosed = errors.New("ws: connection closed")

// Dial connects to a WebSocket on a specified URL.
func Dial(addr string) (net.Conn, error) {
	c := &jsConn{
		events: make(chan event, 1),
		done:   make(chan struct{}),
		read:   make(chan struct{}),
	}
	c.openSocket(addr)
	ev := <-c.events
	if ev.Type == eventOpened {
		// connected - start event loop
		go c.loop()
		return c, nil
	}
	c.Close()
	if ev.Type == eventError {
		return nil, js.Error{Value: ev.Data}
	}
	return nil, fmt.Errorf("unexpected event: %v", ev.Type)
}

type jsConn struct {
	ws js.Value
	cb js.Callback

	events chan event
	done   chan struct{}
	read   chan struct{}

	mu   sync.Mutex
	err  error
	rbuf bytes.Buffer
}

type event struct {
	Type eventType
	Data js.Value
}

type eventType int

const (
	eventError  = eventType(0)
	eventOpened = eventType(1)
	eventClosed = eventType(2)
	eventData   = eventType(3)
)

func (c *jsConn) openSocket(addr string) {
	c.cb = js.NewCallback(func(v []js.Value) {
		ev := event{
			Type: eventType(v[0].Int()),
			Data: v[1],
		}
		select {
		case c.events <- ev:
		case <-c.done:
		}
	})
	setup := js.NewFunction("addr", "event", `
s = new WebSocket(addr);
s.onerror = (e) => {
	event(0, e);
}
s.onopen = (e) => {
	event(1, e);
}
s.onclose = (e) => {
	event(2, e);
}
s.onmessage = (m) => {
	const r = new FileReader();
	r.addEventListener('loadend', function(){
		const data = new Uint8Array(r.result);
		event(3, data);
	})
	r.addEventListener('onerror', (e) => {
		event(0, e);
	});
	r.readAsArrayBuffer(m.data);
}
return s;
`)
	c.ws = setup.Invoke(addr, c.cb)
}

func (c *jsConn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	c.ws.Call("close")
	c.cb.Release()
	return c.err
}

func (c *jsConn) wakeRead() {
	select {
	case c.read <- struct{}{}:
	default:
	}
}

func (c *jsConn) loop() {
	defer c.Close()
	for {
		select {
		case <-c.done:
			return
		case ev := <-c.events:
			switch ev.Type {
			case eventClosed:
				c.mu.Lock()
				c.err = errClosed
				c.mu.Unlock()
				c.wakeRead()
				return
			case eventError:
				c.mu.Lock()
				c.err = js.Error{Value: ev.Data}
				c.mu.Unlock()
				c.wakeRead()
				return
			case eventData:
				arr := ev.Data

				sz := arr.Get("length").Int()

				data := make([]byte, sz)
				m := js.MMap(data)
				err := m.CopyFrom(arr)
				m.Release()

				c.mu.Lock()
				if err == nil {
					c.rbuf.Write(data)
				} else {
					c.err = err
				}
				c.mu.Unlock()
				c.wakeRead()
				if err != nil {
					return
				}
			}
		}
	}
}

func cloneToJS(data []byte) js.Value {
	arr := js.TypedArrayOf(data)
	v := js.New("Uint8Array", arr)
	arr.Release()
	return v
}

func (c *jsConn) send(data []byte) {
	jarr := cloneToJS(data)
	c.ws.Call("send", jarr)
}

func (c *jsConn) Read(b []byte) (int, error) {
	for {
		var (
			n   int
			err error
		)
		c.mu.Lock()
		if c.rbuf.Len() != 0 {
			n, err = c.rbuf.Read(b)
		} else {
			err = c.err
		}
		c.mu.Unlock()
		if err != nil || n != 0 {
			return n, err
		}
		select {
		case <-c.done:
			return 0, io.EOF
		case <-c.read:
		}
	}
}

func (c *jsConn) Write(b []byte) (int, error) {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	if err != nil {
		return 0, err
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
