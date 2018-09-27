//+build !wasm

// Package native provides an API for Native Messaging for Chrome extensions.
//
// See https://developer.chrome.com/apps/nativeMessaging for more details.
package native

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"os"
	"sync"
)

var (
	rmu sync.Mutex
)

// Recv receives a message and decodes it into dst.
func Recv(dst interface{}) error {
	rmu.Lock()
	defer rmu.Unlock()
	var r io.Reader = os.Stdin
	var b [4]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return err
	}
	size := binary.LittleEndian.Uint32(b[:])
	r = io.LimitReader(r, int64(size))
	return json.NewDecoder(r).Decode(dst)
}

var (
	wmu sync.Mutex
	buf = new(bytes.Buffer)
)

// Send sends a message.
func Send(obj interface{}) error {
	wmu.Lock()
	defer wmu.Unlock()
	buf.Reset()
	buf.Write([]byte{0, 0, 0, 0})
	err := json.NewEncoder(buf).Encode(obj)
	if err != nil {
		return err
	}
	data := buf.Bytes()
	binary.LittleEndian.PutUint32(data[:4], uint32(buf.Len()-4))
	_, err = os.Stdout.Write(data)
	return err
}

// RecvBinary receives a single binary message.
//
// The format of the binary messages is not defined in any spec,
// thus a custom message format is used to transfer binary data.
func RecvBinary() ([]byte, error) {
	var m struct {
		Data []byte `json:"d"`
	}
	err := Recv(&m)
	return m.Data, err
}

// SendBinary sends a single binary message.
//
// The format of the binary messages is not defined in any spec,
// thus a custom message format is used to transfer binary data.
func SendBinary(p []byte) error {
	return Send(map[string][]byte{
		"d": p,
	})
}

var accepted bool

// Accept accepts a single client connection.
//
// Notes that communication is done via stdio/stdout, thus an application should not
// read/write any data to/from these streams.
func Accept() io.ReadWriter {
	if accepted {
		panic("accept can only be called once")
	}
	accepted = true
	return &conn{}
}

type conn struct {
	rbuf bytes.Buffer
}

func (conn) Write(p []byte) (int, error) {
	err := SendBinary(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *conn) Read(p []byte) (int, error) {
	if c.rbuf.Len() != 0 {
		return c.rbuf.Read(p)
	}
	data, err := RecvBinary()
	if err != nil {
		return 0, err
	}
	n := copy(p, data)
	c.rbuf.Write(data[n:])
	return n, nil
}
