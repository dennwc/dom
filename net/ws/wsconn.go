//+build !js

// Package ws provides a functionality similar to Go net package on top of WebSockets.
package ws

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func newServer(def http.Handler) *wsServer {
	return &wsServer{
		def:  def,
		stop: make(chan struct{}),
		errc: make(chan error, 1),
		conn: make(chan net.Conn),
	}
}

// Listen listens for incoming connections on a given URL and server other requests with def handler.
func Listen(addr string, def http.Handler) (net.Listener, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	srv := newServer(def)
	mux := http.NewServeMux()
	// TODO: support HTTPS
	srv.h = &http.Server{
		Handler: mux,
	}
	if u.Path == "" {
		u.Path = "/"
	}
	mux.HandleFunc(u.Path, srv.handleWS)
	if def != nil && strings.Trim(u.Path, "/") != "" {
		mux.Handle("/", def)
	}
	lis, err := net.Listen("tcp", u.Host)
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(srv.errc)
		srv.errc <- srv.h.Serve(lis)
	}()
	return srv, nil
}

type wsServer struct {
	def  http.Handler
	h    *http.Server
	stop chan struct{}
	errc chan error
	conn chan net.Conn
}

func (s *wsServer) Accept() (net.Conn, error) {
	select {
	case <-s.stop:
		return nil, fmt.Errorf("server stopped")
	case err := <-s.errc:
		_ = s.Close()
		return nil, err
	case c := <-s.conn:
		return c, nil
	}
}

func (s *wsServer) Close() error {
	select {
	case <-s.stop:
	default:
		close(s.stop)
	}
	if s.h != nil {
		return s.h.Close()
	}
	return nil
}

func (s *wsServer) Addr() net.Addr {
	return wsAddr{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: make it configurable
	},
}

func (s *wsServer) handleWS(w http.ResponseWriter, r *http.Request) {
	if s.def != nil && !websocket.IsWebSocketUpgrade(r) {
		s.def.ServeHTTP(w, r)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	done := make(chan struct{})
	c := &wsConn{c: conn, done: done}
	defer conn.Close()
	select {
	case <-s.stop:
		return
	case s.conn <- c:
		select {
		case <-s.stop:
		case <-done:
		}
	}
}

var dialer = &websocket.Dialer{}

// Dial connects to a WebSocket on a specified URL.
func Dial(addr string) (net.Conn, error) {
	conn, _, err := dialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}
	return &wsConn{c: conn}, nil
}

type wsConn struct {
	c    *websocket.Conn
	cur  io.Reader
	done chan struct{}
}

func (c *wsConn) Read(b []byte) (int, error) {
	for {
		if c.cur != nil {
			n, err := c.cur.Read(b)
			if err == nil || err != io.EOF {
				return n, err
			} else if err == io.EOF && n != 0 {
				return n, nil
			}
			// EOF, n == 0
			c.cur = nil
		}
		_, r, err := c.c.NextReader()
		if err != nil {
			return 0, err
		}
		c.cur = r
	}
}

func (c *wsConn) Write(b []byte) (int, error) {
	// TODO: buffer writes
	err := c.c.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *wsConn) Close() error {
	if c.done == nil {
		return c.c.Close()
	}
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return nil
}

func (c *wsConn) LocalAddr() net.Addr {
	return wsAddr{}
}

func (c *wsConn) RemoteAddr() net.Addr {
	return wsAddr{}
}

func (c *wsConn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	if err := c.SetWriteDeadline(t); err != nil {
		return err
	}
	return nil
}

func (c *wsConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

func (c *wsConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}
