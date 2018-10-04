package webrtc

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"

	"github.com/dennwc/dom/js"
)

type peerChannel struct {
	c    *peerConnection
	name string
	v    js.Value

	ready chan struct{}
	done  chan struct{}
	read  chan struct{}

	mu   sync.Mutex
	err  error
	rbuf bytes.Buffer
}

func (c *peerChannel) wakeReaders() {
	select {
	case c.read <- struct{}{}:
	default:
	}
}

func (c *peerChannel) handleEvent(e chanEvent) {
	switch e.Type {
	case eventError:
		c.mu.Lock()
		c.err = js.Error{Value: e.Data}
		c.mu.Unlock()
		c.Close()
	case eventClosed:
		c.mu.Lock()
		c.err = io.EOF
		c.mu.Unlock()
		c.Close()
	case eventOpened:
		close(c.ready)
	case eventMessage:
		c.mu.Lock()
		c.rbuf.WriteString(e.Data.String())
		c.mu.Unlock()
		c.wakeReaders()
	}
}

func (c *peerChannel) Read(b []byte) (int, error) {
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

func (c *peerChannel) Write(b []byte) (int, error) {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	if err != nil {
		return 0, err
	}
	// TODO: check if we can send byte arrays
	c.v.Call("send", string(b))
	return len(b), nil
}

func (c *peerChannel) LocalAddr() net.Addr {
	// TODO: pretty sure it's possible to get an address
	return peerAddr{}
}

func (c *peerChannel) RemoteAddr() net.Addr {
	// TODO: pretty sure it's possible to get an address
	return peerAddr{}
}

func (c *peerChannel) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.c.Close() // TODO: close only current channel
}

func (c *peerChannel) SetDeadline(t time.Time) error {
	// TODO
	return nil
}

func (c *peerChannel) SetReadDeadline(t time.Time) error {
	// TODO
	return nil
}

func (c *peerChannel) SetWriteDeadline(t time.Time) error {
	// TODO
	return nil
}

var _ net.Addr = peerAddr{}

type peerAddr struct{}

func (peerAddr) Network() string {
	return "webrtc"
}

func (peerAddr) String() string {
	return "webrtc://localhost"
}
