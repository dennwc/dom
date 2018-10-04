package webrtc

import (
	"encoding/json"
	"net"
)

type Peers struct {
	c   *peerConnection
	lis Listener

	local peerDesc
	ans   OfferListener
}

func (p *Peers) Close() error {
	p.lis.Close()
	if p.c != nil {
		return p.c.Close()
	}
	return nil
}

func (p *Peers) Accept() (*PeerInfo, error) {
	// listen for incoming answers, but don't use any of them yet
	data, err := p.lis.Accept()
	if err != nil {
		p.c.Close()
		return nil, err
	}
	var info peerDesc
	if err = json.Unmarshal(data, &info); err != nil {
		p.c.Close()
		return nil, err
	}
	return &PeerInfo{peers: p, info: info}, nil
}

type PeerInfo struct {
	peers *Peers
	info  peerDesc
}

func (p *PeerInfo) Dial() (net.Conn, error) {
	if p.peers.ans == nil {
		return p.dialActive()
	}
	return p.dialPassive()
}

func (p *PeerInfo) dial() (net.Conn, error) {
	// take ownership of the connection
	c := p.peers.c
	p.peers.c = nil

	// now we should only wait for a state change to "connected"
	// but instead we will wait for a data stream to come online
	ch, err := c.WaitChannel(primaryChan)
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func (p *PeerInfo) dialActive() (net.Conn, error) {
	// if we are initiating a connection, we have just received an info from peer
	// and we are ready to apply its configuration and start dialing
	c := p.peers.c

	// switch to this peer and try to dial it
	err := c.SetRemoteDescription(p.info.Desc)
	if err != nil {
		c.Close()
		return nil, err
	}

	err = c.SetICECandidates(p.info.ICEs)
	if err != nil {
		c.Close()
		return nil, err
	}
	return p.dial()
}

func (p *PeerInfo) dialPassive() (net.Conn, error) {
	// if we are listening for connections, so we need to collect our local ICEs
	// and send an answer with our info
	c := p.peers.c

	// prepare to collect local ICEs
	collectICEs := c.CollectICEs()

	// switch to this peer and start dialing it (he might reject)
	err := c.SetRemoteDescription(p.info.Desc)
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
	p.peers.local = peerDesc{
		Desc: answer, ICEs: ices,
	}

	// send our information to the peer
	data, err := json.Marshal(p.peers.local)
	if err != nil {
		return nil, err
	}
	err = p.peers.ans.Answer(data)
	if err != nil {
		return nil, err
	}
	return p.dial()
}
