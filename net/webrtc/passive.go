package webrtc

import (
	"encoding/json"
	"net"
)

// offerStream is an implementation of peer discovery when listening (passive).
type offerStream struct {
	self string // user id
	c    *peerConnection

	local  connInfo
	offers OfferStream
}

func (p *offerStream) Close() error {
	p.offers.Close()
	if p.c != nil {
		return p.c.Close()
	}
	return nil
}

func (p *offerStream) Accept() (Peer, error) {
	// we are listening for connections, so we take a next offer and present it to user to decide if he wants to dial
	offer, err := p.offers.Next()
	if err != nil {
		return nil, err
	}
	s := offer.Info()
	var info connInfo
	if err = json.Unmarshal(s.Data, &info); err != nil {
		return nil, err
	}
	return &peerOffer{s: p, uid: s.UID, info: info, offer: offer}, nil
}

type peerOffer struct {
	s     *offerStream
	uid   string
	info  connInfo
	offer Offer
}

func (p *peerOffer) UID() string {
	return p.uid
}

func (p *peerOffer) Dial() (net.Conn, error) {
	// we are listening for connections, so we need to collect our local ICEs
	// and send an answer with our info
	c := p.s.c

	// prepare to collect local ICEs
	collectICEs := c.CollectICEs()

	// switch to this peer and start dialing it (he might reject)
	err := c.SetRemoteDescription(p.info.SDP)
	if err != nil {
		c.Close()
		return nil, err
	}

	// set remote candidates
	err = c.SetICECandidates(p.info.ICEs)
	if err != nil {
		c.Close()
		return nil, err
	}

	// we are ready to answer
	answer, err := c.CreateAnswer()
	if err != nil {
		c.Close()
		return nil, err
	}

	// switch to the config that we propose
	err = c.SetLocalDescription(answer)
	if err != nil {
		c.Close()
		return nil, err
	}

	// this allows us to gather local ICEs
	ices, err := collectICEs()
	if err != nil {
		c.Close()
		return nil, err
	}

	// now we know our own parameters
	local := connInfo{
		SDP: answer, ICEs: ices,
	}

	// send our information to the peer
	data, err := json.Marshal(local)
	if err != nil {
		return nil, err
	}
	err = p.offer.Answer(Signal{UID: p.s.self, Data: data})
	if err != nil {
		return nil, err
	}
	// take ownership of the connection
	p.s.c = nil

	// now we should only wait for a state change to "connected"
	// but instead we will wait for a data stream to come online
	ch, err := c.WaitChannel(primaryChan)
	if err != nil {
		return nil, err
	}
	return ch, nil
}
