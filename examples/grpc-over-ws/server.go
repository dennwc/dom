package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/grpc"

	"github.com/dennwc/dom/examples/grpc-over-ws/protocol"
	"github.com/dennwc/dom/net/ws"
)

//go:generate GOOS=js GOARCH=wasm go build -o app.wasm ./client.go

func main() {
	s := server{}

	srv := grpc.NewServer()
	protocol.RegisterService(srv, s)

	const host = "localhost:8080"

	handler := http.FileServer(http.Dir("."))
	lis, err := ws.Listen("ws://"+host+"/ws", handler)
	if err != nil {
		panic(err)
	}
	defer lis.Close()

	log.Printf("listening on http://%s", host)
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}

type server struct{}

func (server) Hello(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello, %s!", name), nil
}
