package main

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/examples/grpc-over-ws/protocol"
	"github.com/dennwc/dom/net/ws"
)

func dialer(s string, dt time.Duration) (net.Conn, error) {
	return ws.Dial(s)
}

func main() {
	p1 := dom.Doc.CreateElement("p")
	dom.Body.AppendChild(p1)

	inp := dom.Doc.NewInput("text")
	p1.AppendChild(inp)

	btn := dom.Doc.NewButton("Go!")
	p1.AppendChild(btn)

	ch := make(chan string, 1)
	btn.OnClick(func(_ dom.Event) {
		ch <- inp.Value()
	})

	conn, err := grpc.Dial("ws://localhost:8080/ws", grpc.WithDialer(dialer), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cli := protocol.AsService(conn)

	printMsg := func(s string) {
		p := dom.Doc.CreateElement("p")
		p.SetTextContent(s)
		dom.Body.AppendChild(p)
	}

	ctx := context.Background()
	for {
		name := <-ch
		printMsg("say hello to: " + name)

		txt, err := cli.Hello(ctx, name)
		if err != nil {
			panic(err)
		}
		printMsg(txt)
	}

	dom.Loop()
}
