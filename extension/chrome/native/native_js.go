//+build wasm

package native

import (
	"encoding/base64"
	"fmt"
	"io"
	"syscall/js"
)

// NewApp returns a application instance with a name, as specified in the extension manifest.
func NewApp(name string) *App {
	return &App{name: name}
}

// Msg is a type of the message.
type Msg = map[string]interface{}

// App represents a native messaging connector running on the host.
type App struct {
	name string
}

// Send runs a native extension, passes a single message to it and waits for the response.
// Native connector will likely be killed after the response is received.
func (app *App) Send(o Msg) js.Value {
	resp := make(chan js.Value, 1)
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		cb.Release()
		if len(args) != 0 {
			resp <- args[0]
		} else {
			resp <- js.Undefined()
		}
		return nil
	})
	js.Global().Get("chrome").Get("runtime").Call("sendNativeMessage", app.name, o, cb)
	return <-resp
}

// SendBinary sends a single binary message and receives a response.
//
// The format of the binary messages is not defined in any spec,
// thus a custom message format is used to transfer binary data.
func (app *App) SendBinary(p []byte) ([]byte, error) {
	resp := app.Send(Msg{
		"d": base64.StdEncoding.EncodeToString(p),
	})
	if resp == js.Undefined() {
		return nil, fmt.Errorf("application failed")
	}
	return base64.StdEncoding.DecodeString(resp.Get("d").String())
}

// Dial runs a native extension connector to send multiple messages.
// This is more efficient than App.Send because an extension won't be killed after each received message.
func (app *App) Dial() (*Conn, error) {
	port := js.Global().Get("chrome").Get("runtime").Call("connectNative", app.name)
	c := &Conn{
		app: app, port: port,
		r: make(chan js.Value, 1),
	}
	c.cRecv = js.FuncOf(c.recv)
	c.cDisc = js.FuncOf(c.disconnect)
	port.Get("onMessage").Call("addListener", c.cRecv)
	port.Get("onDisconnect").Call("addListener", c.cDisc)
	return c, nil
}

// Conn represents a connection to a native extension.
type Conn struct {
	app  *App
	port js.Value

	err error
	r   chan js.Value

	cRecv, cDisc js.Func
}

func (c *Conn) recv(_ js.Value, args []js.Value) interface{} {
	c.r <- args[0]
	return nil
}
func (c *Conn) disconnect(_ js.Value, _ []js.Value) interface{} {
	c.cDisc.Release()
	c.cRecv.Release()
	c.err = io.EOF
	close(c.r)
	return nil
}

// Recv receives a single message.
func (c *Conn) Recv() (js.Value, error) {
	if c.err != nil {
		return js.Undefined(), c.err
	}
	v, ok := <-c.r
	if !ok {
		return js.Undefined(), c.err
	}
	return v, nil
}

// Send sends a single message.
func (c *Conn) Send(m Msg) error {
	c.port.Call("postMessage", m)
	return nil
}

// RecvBinary receives a single binary message.
//
// The format of the binary messages is not defined in any spec,
// thus a custom message format is used to transfer binary data.
func (c *Conn) RecvBinary() ([]byte, error) {
	v, err := c.Recv()
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(v.Get("d").String())
}

// SendBinary sends a single binary message.
//
// The format of the binary messages is not defined in any spec,
// thus a custom message format is used to transfer binary data.
func (c *Conn) SendBinary(p []byte) error {
	return c.Send(Msg{
		"d": base64.StdEncoding.EncodeToString(p),
	})
}
