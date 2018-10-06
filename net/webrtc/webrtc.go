//+build js,wasm

package webrtc

import (
	"encoding/json"
	"net"

	"github.com/dennwc/dom/js"
)

// Peers represents a dynamic list of peers that were discovered via signalling.
type Peers interface {
	// Accept queries an information about next available peer. It won't connect to it automatically.
	Accept() (Peer, error)
	// Close ends a peer discovery process.
	Close() error
}

// Peer represents an information about a potential peer.
type Peer interface {
	// UID returns a optional user ID of this peer.
	UID() string
	// Dial establishes a new connection to this peer.
	Dial() (net.Conn, error)
}

// New creates a new local peer with a given ID that will use a specific server for peer discovery.
func New(uid string, s Signalling) *Local {
	return &Local{
		uid: uid, s: s,
	}
}

// Local represents an information about local peer.
type Local struct {
	uid string
	s   Signalling
}

// Listen starts a passive peer discovery process by waiting for incoming discovery requests.
func (l *Local) Listen() (Peers, error) {
	offers, err := l.s.Listen(l.uid)
	if err != nil {
		return nil, err
	}
	c := newPeerConnection()
	return &offerStream{self: l.uid, c: c, offers: offers}, nil
}

// Discover starts an active peer discovery process by broadcasting a discovery request.
func (l *Local) Discover() (Peers, error) {
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
	local := connInfo{
		SDP: offer, ICEs: ices,
	}
	// encode and broadcast
	data, err := json.Marshal(local)
	if err != nil {
		c.Close()
		return nil, err
	}
	answers, err := l.s.Broadcast(Signal{UID: l.uid, Data: data})
	if err != nil {
		c.Close()
		return nil, err
	}
	return &answerStream{self: l.uid, c: c, answers: answers, local: local}, nil
}

const primaryChan = "primary"

// connInfo combines SDP and ICE data of a specific peer.
type connInfo struct {
	SDP  js.Value   `json:"sdp"`
	ICEs []js.Value `json:"ices"`
}
