package webrtc

import (
	"encoding/json"
)

func (s *Local) Discover(net Signaling) (*Peers, error) {
	c := newPeerConnection()
	for _, name := range s.chans {
		c.NewDataChannel(name)
	}
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
	return &Peers{l: s, c: c, lis: lis, local: local}, nil
}
