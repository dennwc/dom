package webrtc

import (
	"fmt"
	"sync/atomic"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/js"
)

var lastID uint32

func newPeerConnection() *peerConnection {
	p := &peerConnection{
		id: atomic.AddUint32(&lastID, 1),
		v:  js.New("RTCPeerConnection"),
	}
	p.g = p.v.NewCallbackGroup()
	return p
}

type peerConnection struct {
	id uint32
	v  js.Value
	g  *js.CallbackGroup
}

func (c *peerConnection) Close() error {
	c.g.Release()
	// TODO: close channels?
	return nil
}

func (c *peerConnection) OnDataChannel(fnc func(ch js.Value)) {
	c.g.Set("ondatachannel", func(v []js.Value) {
		ch := v[0].Get("channel")
		dom.ConsoleLog("chan:", c.id, ch)
		fnc(ch)
	})
}

func (c *peerConnection) onICECandidate(fnc func(cand js.Value)) js.Callback {
	cb := js.NewCallback(func(v []js.Value) {
		cand := v[0].Get("candidate")
		dom.ConsoleLog("candidate:", c.id, cand)
		fnc(cand)
	})
	c.v.Set("onicecandidate", cb)
	return cb
}

func (c *peerConnection) OnICECandidate(fnc func(cand js.Value)) {
	cb := c.onICECandidate(fnc)
	c.g.Add(cb)
}

func (c *peerConnection) NewDataChannel(name string) js.Value {
	ch := c.v.Call("createDataChannel", name)
	g := ch.NewCallbackGroup()
	g.Set("onopen", func(v []js.Value) {
		dom.ConsoleLog("open", c.id, name)
		go func() {
			ch.Call("send", fmt.Sprintf("hello from peer %d", c.id))
		}()
	})
	g.Set("onclose", func(v []js.Value) {
		dom.ConsoleLog("close", c.id, name)
	})
	return ch
}

func (c *peerConnection) AddICECandidate(v js.Value) error {
	dom.ConsoleLog("add candidate:", c.id, v)
	_, err := c.v.Call("addIceCandidate", v).Await()
	return err
}

func (c *peerConnection) SetLocalDescription(d js.Value) error {
	dom.ConsoleLog("set local desc:", c.id, d)
	_, err := c.v.Call("setLocalDescription", d).Await()
	return err
}

func (c *peerConnection) CreateOffer() (js.Value, error) {
	vals, err := c.v.Call("createOffer").Await()
	if err != nil {
		return js.Value{}, err
	}
	offer := vals[0]
	return offer, nil
}

func (c *peerConnection) SetRemoteDescription(v js.Value) error {
	dom.ConsoleLog("set remote desc:", c.id, v)
	_, err := c.v.Call("setRemoteDescription", v).Await()
	return err
}

func (c *peerConnection) CreateAnswer() (js.Value, error) {
	vals, err := c.v.Call("createAnswer").Await()
	if err != nil {
		return js.Value{}, err
	}
	answer := vals[0]
	return answer, nil
}

type iceFunc func() ([]js.Value, error)

func (c *peerConnection) CollectICEs() iceFunc {
	var (
		cb   js.Callback
		done = make(chan struct{})
		ices []js.Value
	)
	cb = c.onICECandidate(func(cand js.Value) {
		if !cand.Valid() {
			close(done)
			cb.Release()
			return
		}
		ices = append(ices, cand)
	})
	return func() ([]js.Value, error) {
		<-done // TODO: listen on some kind of error channel
		if len(ices) == 0 {
			return nil, fmt.Errorf("no ICE candidates collected")
		}
		return ices, nil
	}
}

func (c *peerConnection) SetICECandidates(ices []js.Value) error {
	for _, ice := range ices {
		if err := c.AddICECandidate(ice); err != nil {
			return err
		}
	}
	// signal "no more ICEs"
	// TODO: docs says it's should be sent, but this call fails
	//if err := c.AddICECandidate(js.ValueOf("")); err != nil {
	//	return err
	//}
	return nil
}
