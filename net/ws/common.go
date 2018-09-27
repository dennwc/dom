package ws

import "net"

var _ net.Addr = wsAddr{}

type wsAddr struct {
}

func (wsAddr) Network() string {
	return "ws"
}

func (wsAddr) String() string {
	// TODO: proper address, if possible
	return "ws://localhost"
}
