//+build js,wasm

package webrtc

import (
	"encoding/json"

	"github.com/dennwc/dom/js"
)

type Listener interface {
	Accept() ([]byte, error)
	Close() error
}

type OfferListener interface {
	Listener
	Answer(data []byte) error
}

func Listen(lis OfferListener) (*Peers, error) {
	c := newPeerConnection()
	return &Peers{c: c, lis: lis, ans: lis}, nil
}

type Signaling interface {
	Broadcast(data []byte) (Listener, error)
}

func Discover(net Signaling) (*Peers, error) {
	c := newPeerConnection()
	c.NewDataChannel(primaryChan)
	// prepare to collect local ICEs
	collectICEs := c.CollectICEs()

	// create an offer and activate it
	offer, err := c.CreateOffer()
	if err != nil {
		c.Close()
		return nil, err
	}
	err = c.SetLocalDescription(offer)
	if err != nil {
		c.Close()
		return nil, err
	}
	// collect all local ICE candidates
	ices, err := collectICEs()
	if err != nil {
		c.Close()
		return nil, err
	}
	local := peerDesc{
		Desc: offer, ICEs: ices,
	}
	// encode and broadcast
	data, err := json.Marshal(local)
	if err != nil {
		c.Close()
		return nil, err
	}
	lis, err := net.Broadcast(data)
	if err != nil {
		c.Close()
		return nil, err
	}
	return &Peers{c: c, lis: lis, local: local}, nil
}

const primaryChan = "primary"

type peerDesc struct {
	Desc js.Value   `json:"desc"`
	ICEs []js.Value `json:"ices"`
}
