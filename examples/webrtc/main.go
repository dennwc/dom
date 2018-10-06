package main

import (
	"fmt"
	"io"
	"time"

	"github.com/dennwc/dom/net/webrtc"
)

func main() {
	sig := NewChannel()

	go Alice(sig)
	Bob(sig)
}

func Alice(sig webrtc.Signalling) {
	const name = "alice"
	p1 := webrtc.New(name, sig)

	fmt.Println(name + ": peer discovery started")
	peers, err := p1.Discover()
	if err != nil {
		panic(err)
	}
	defer peers.Close()

	fmt.Println(name + ": waiting for peers")
	info, err := peers.Accept()
	if err != nil {
		panic(err)
	}

	pname := info.UID()
	fmt.Printf(name+": dialing peer %q\n", pname)
	conn, err := info.Dial()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println(name + ": connected!")
	_, err = fmt.Fprintf(conn, "hello from %q\n", name)
	if err != nil {
		panic(err)
	}
	fmt.Println(name + ": sent data")

	buf := make([]byte, 128)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		fmt.Printf(name+": msg from %q: %q\n", pname, string(buf[:n]))
	}
}

func Bob(sig webrtc.Signalling) {
	const name = "bob"

	p2 := webrtc.New(name, sig)
	fmt.Println(name + ": listening for offers")
	peers, err := p2.Listen()
	if err != nil {
		panic(err)
	}
	defer peers.Close()

	info, err := peers.Accept()
	if err != nil {
		panic(err)
	}

	pname := info.UID()
	fmt.Printf(name+": dialing peer %q\n", pname)
	conn, err := info.Dial()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println(name + ": connected!")
	_, err = fmt.Fprintf(conn, "hello from %q\n", name)
	if err != nil {
		panic(err)
	}
	fmt.Println(name + ": sent data")

	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Printf(name+": msg from %q: %q\n", pname, string(buf[:n]))

	for {
		_, err = conn.Write([]byte(time.Now().String() + "\n"))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
	}
}

// NewChannel creates a new signalling channel. It expects exactly one call to Broadcast and exactly one call to Listen.
func NewChannel() webrtc.Signalling {
	return &signalChannel{
		broadcast: make(chan webrtc.Signal, 1),
		accept:    make(chan webrtc.Signal, 1),
	}
}

type signalChannel struct {
	broadcast chan webrtc.Signal
	accept    chan webrtc.Signal
}

func (b *signalChannel) Broadcast(s webrtc.Signal) (webrtc.AnswerStream, error) {
	b.broadcast <- s
	close(b.broadcast)
	return &answers{accept: b.accept}, nil
}

type answers struct {
	accept <-chan webrtc.Signal
}

func (a *answers) Next() (webrtc.Signal, error) {
	s, ok := <-a.accept
	if !ok {
		return webrtc.Signal{}, io.EOF
	}
	return s, nil
}

func (a *answers) Close() error {
	ch := make(chan webrtc.Signal)
	close(ch)
	a.accept = ch
	return nil
}

func (b *signalChannel) Listen(uid string) (webrtc.OfferStream, error) {
	return &offers{broadcast: b.broadcast, accept: b.accept}, nil
}

type offers struct {
	broadcast <-chan webrtc.Signal
	accept    chan<- webrtc.Signal
}

func (o *offers) Next() (webrtc.Offer, error) {
	s, ok := <-o.broadcast
	if !ok {
		return nil, io.EOF
	}
	return &offer{accept: o.accept, s: s}, nil
}

func (o *offers) Close() error {
	ch := make(chan webrtc.Signal)
	close(ch)
	o.broadcast = ch
	return nil
}

type offer struct {
	accept chan<- webrtc.Signal
	s      webrtc.Signal
}

func (o *offer) Answer(s webrtc.Signal) error {
	o.accept <- s
	close(o.accept)
	o.accept = nil
	return nil
}

func (o *offer) Info() webrtc.Signal {
	return o.s
}
