package webrtc

import (
	"encoding/json"
	"net"
)

// answerStream is an implementation of peer discovery when broadcasting (active).
type answerStream struct {
	self string // user id
	c    *peerConnection

	local   connInfo // local WebRTC info
	answers AnswerStream
}

func (p *answerStream) Close() error {
	p.answers.Close()
	if p.c != nil {
		return p.c.Close()
	}
	return nil
}

func (p *answerStream) Accept() (Peer, error) {
	// get the next answer, but don't use it yet
	ans, err := p.answers.Next()
	if err != nil {
		return nil, err
	}
	var info connInfo
	if err = json.Unmarshal(ans.Data, &info); err != nil {
		return nil, err
	}
	return &peerAnswer{s: p, uid: ans.UID, info: info}, nil
}

type peerAnswer struct {
	s    *answerStream
	uid  string
	info connInfo
}

func (p *peerAnswer) UID() string {
	return p.uid
}

func (p *peerAnswer) Dial() (net.Conn, error) {
	// if we are initiating a connection, we have just received an info from peer
	// and we are ready to apply its configuration and start dialing
	c := p.s.c

	// switch to this peer and try to dial it
	err := c.SetRemoteDescription(p.info.SDP)
	if err != nil {
		c.Close()
		return nil, err
	}

	err = c.SetICECandidates(p.info.ICEs)
	if err != nil {
		c.Close()
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
