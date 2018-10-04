package webrtc

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/js"
)

const debug = false

var lastID uint32

func newPeerConnection() *peerConnection {
	p := &peerConnection{
		id: atomic.AddUint32(&lastID, 1),
		v:  js.New("RTCPeerConnection"),

		errc:    make(chan error, 1),
		newChan: make(chan string, 1),
		chans:   make(map[string]*peerChannel),
	}
	p.registerStateHandlers()
	return p
}

type peerConnection struct {
	id uint32
	v  js.Value
	g  *js.CallbackGroup

	errc    chan error
	newChan chan string
	cmu     sync.Mutex
	chans   map[string]*peerChannel
}

func (c *peerConnection) Channel(name string) *peerChannel {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.chans[name]
}

func (c *peerConnection) WaitChannel(name string) (*peerChannel, error) {
	ch := c.Channel(name)
	if ch == nil {
		// wait for the channel to appear
		select {
		case err := <-c.errc:
			c.Close()
			return nil, err
		case cname := <-c.newChan:
			if cname != name {
				c.Close()
				return nil, fmt.Errorf("unexpected channel: %q", cname)
			}
			ch = c.Channel(name)
		}
	}
	select {
	case <-ch.done:
		return nil, errors.New("webrtc: channel closed")
	case <-ch.ready:
		return ch, nil
	}
}

func (c *peerConnection) newChannel(name string, ch js.Value, ready bool) {
	peer := &peerChannel{
		c: c, name: name, v: ch,
		ready: make(chan struct{}),
		done:  make(chan struct{}),
		read:  make(chan struct{}),
	}
	if ready {
		close(peer.ready)
	}
	c.cmu.Lock()
	defer c.cmu.Unlock()
	c.chans[name] = peer
	// keep only the last notification on the channel
	select {
	case c.newChan <- name:
	default:
		// remove value from the channel and retry
		select {
		case <-c.newChan:
		default:
			// holding mutex, we are the only who can send
			c.newChan <- name
		}
	}
}

type eventType int

const (
	eventError   = eventType(0)
	eventNew     = eventType(1)
	eventOpened  = eventType(2)
	eventClosed  = eventType(3)
	eventMessage = eventType(4)
)

type chanEvent struct {
	Name string
	Type eventType
	Data js.Value
}

func (c *peerConnection) registerStateHandlers() {
	c.g = c.v.NewCallbackGroup()
	// always register an error and connection state event handlers
	c.g.Set("onerror", func(v []js.Value) {
		dom.ConsoleLog("error:", c.id, v[0])
		select {
		case c.errc <- js.Error{Value: v[0]}:
		default:
		}
	})
	if debug {
		c.g.Set("oniceconnectionstatechange", func(v []js.Value) {
			dom.ConsoleLog("state:", c.id, c.v.Get("iceConnectionState"))
		})
		c.g.Set("onsignalingstatechange", func(v []js.Value) {
			dom.ConsoleLog("sig state:", c.id, c.v.Get("signalingState"))
		})
		c.g.Set("onicegatheringstatechange", func(v []js.Value) {
			dom.ConsoleLog("gather state:", c.id, c.v.Get("iceGatheringState"))
		})
	}

	// handle incoming data channels
	jfnc := js.NewFunction("fnc", `
return function(ce) {
	const ch = ce.channel;
	const name = ch.label;

	fnc(name, 1, ch);
	ch.onerror = (e) => {
		fnc(name, 0, e);
	}
	ch.onopen = (e) => {
		fnc(name, 2, e);
	}
	ch.onclose = (e) => {
		fnc(name, 3, e);
	}
	ch.onmessage = (e) => {
		fnc(name, 4, e.data);
	}
}
`)
	cb := c.newChanEventCallback()
	c.v.Set("ondatachannel", jfnc.Invoke(cb))
}

func (c *peerConnection) Close() error {
	c.v.Call("close")
	c.g.Release()
	return nil
}

func (c *peerConnection) onICECandidate(fnc func(cand js.Value)) js.Callback {
	cb := js.NewCallback(func(v []js.Value) {
		cand := v[0].Get("candidate")
		if debug {
			dom.ConsoleLog("candidate:", c.id, cand)
		}
		fnc(cand)
	})
	c.v.Set("onicecandidate", cb)
	return cb
}

func (c *peerConnection) OnICECandidate(fnc func(cand js.Value)) {
	cb := c.onICECandidate(fnc)
	c.g.Add(cb)
}

func (c *peerConnection) newChanEventCallback() js.Callback {
	cb := js.NewCallback(func(v []js.Value) {
		e := chanEvent{
			Name: v[0].String(),
			Type: eventType(v[1].Int()),
			Data: v[2],
		}

		if debug {
			dom.ConsoleLog("chan:", c.id, e.Name, e.Type, e.Data)
		}
		switch e.Type {
		case eventNew:
			c.newChannel(e.Name, e.Data, false)
		default:
			ch := c.Channel(e.Name)
			if ch != nil {
				ch.handleEvent(e)
			}
		}
	})
	c.g.Add(cb)
	return cb
}

func (c *peerConnection) NewDataChannel(name string) {
	// TODO: it can be the same callback
	cb := c.newChanEventCallback()
	// handle initiated data channels
	ch := js.NewFunction("v", "name", "fnc", `
	const ch = v.createDataChannel(name);
	ch.onerror = (e) => {
		fnc(name, 0, e);
	}
	ch.onopen = (e) => {
		fnc(name, 2, e);
	}
	ch.onclose = (e) => {
		fnc(name, 3, e);
	}
	ch.onmessage = (e) => {
		fnc(name, 4, e.data);
	}
	return ch;
`).Invoke(c.v, name, cb)

	// new channels are never reported as "opened"
	// but if they were added before connection was established, they do
	c.newChannel(name, ch, false)
}

func (c *peerConnection) AddICECandidate(v js.Value) error {
	if debug {
		dom.ConsoleLog("add candidate:", c.id, v)
	}
	_, err := c.v.Call("addIceCandidate", v).Await()
	return err
}

func (c *peerConnection) SetLocalDescription(d js.Value) error {
	if debug {
		dom.ConsoleLog("set local desc:", c.id, d)
	}
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
	if debug {
		dom.ConsoleLog("set remote desc:", c.id, v)
	}
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
