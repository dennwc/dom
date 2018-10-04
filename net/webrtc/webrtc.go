//+build js,wasm

package webrtc

import "github.com/dennwc/dom/js"

type Listener interface {
	Accept() ([]byte, error)
	Close() error
}

type Signaling interface {
	Broadcast(data []byte) (Listener, error)
}

func New() *Local {
	l := &Local{}
	l.SetChannels("data")
	return l
}

type Local struct {
	chans []string
}

func (s *Local) SetChannels(names ...string) {
	s.chans = append([]string{}, names...)
}

type peerDesc struct {
	Desc js.Value   `json:"desc"`
	ICEs []js.Value `json:"ices"`
}
