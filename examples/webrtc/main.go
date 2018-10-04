package main

import (
	"fmt"
	"io"
	"time"

	"github.com/dennwc/dom/net/webrtc"
)

type discovery struct {
	send chan<- []byte
	recv <-chan []byte
}

func (d *discovery) Broadcast(data []byte) (webrtc.Listener, error) {
	d.send <- data
	close(d.send)
	return &discoveryLis{recv: d.recv}, nil
}

type discoveryLis struct {
	recv <-chan []byte
}

func (l *discoveryLis) Accept() ([]byte, error) {
	data, ok := <-l.recv
	if !ok {
		return nil, io.EOF
	}
	return data, nil
}

func (l *discoveryLis) Close() error {
	return nil
}

type offerLis struct {
	recv <-chan []byte
	send chan<- []byte
}

func (l *offerLis) Answer(data []byte) error {
	l.send <- data
	close(l.send)
	l.send = nil
	return nil
}

func (l *offerLis) Accept() ([]byte, error) {
	data, ok := <-l.recv
	if !ok {
		return nil, io.EOF
	}
	return data, nil
}

func (l *offerLis) Close() error {
	if l.send != nil {
		close(l.send)
		l.send = nil
	}
	return nil
}

func main() {
	ch1to2 := make(chan []byte, 1)
	ch2to1 := make(chan []byte, 1)

	go func() {
		fmt.Println("1: peer discovery started")
		peers, err := webrtc.Discover(&discovery{
			send: ch1to2, recv: ch2to1,
		})
		if err != nil {
			panic(err)
		}
		defer peers.Close()

		fmt.Println("1: waiting for peers")
		info, err := peers.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("1: dialing peer")
		conn, err := info.Dial()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		fmt.Println("1: connected!")
		_, err = conn.Write([]byte("hello from 1\n"))
		if err != nil {
			panic(err)
		}
		fmt.Println("1: sent data")

		buf := make([]byte, 128)
		for {
			n, err := conn.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			fmt.Println("1:", string(buf[:n]))
		}
	}()

	fmt.Println("2: waiting for offers")
	peers, err := webrtc.Listen(&offerLis{
		send: ch2to1, recv: ch1to2,
	})
	if err != nil {
		panic(err)
	}
	defer peers.Close()

	info, err := peers.Accept()
	if err != nil {
		panic(err)
	}

	fmt.Println("2: dialing peer")
	conn, err := info.Dial()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("2: connected!")
	_, err = conn.Write([]byte("hello from 2\n"))
	if err != nil {
		panic(err)
	}
	fmt.Println("2: sent data")

	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println("2:", string(buf[:n]))

	for {
		_, err = conn.Write([]byte(time.Now().String() + "\n"))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 5)
	}
}
