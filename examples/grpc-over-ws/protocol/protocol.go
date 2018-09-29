package protocol

import (
	"context"

	xcontext "golang.org/x/net/context"

	"google.golang.org/grpc"
)

//go:generate protoc --proto_path=$GOPATH/src:. --gogo_out=plugins=grpc:. ./hello.proto

type Service interface {
	Hello(ctx context.Context, name string) (string, error)
}

func AsService(cc *grpc.ClientConn) Service {
	cli := NewHelloServiceClient(cc)
	return implClient{cli}
}

func RegisterService(srv *grpc.Server, s Service) {
	RegisterHelloServiceServer(srv, implServer{s})
}

type implClient struct {
	cli HelloServiceClient
}

func (c implClient) Hello(ctx context.Context, name string) (string, error) {
	resp, err := c.cli.Hello(ctx, &HelloReq{Name: name})
	if err != nil {
		return "", err
	}
	return resp.Text, nil
}

type implServer struct {
	srv Service
}

func (s implServer) Hello(ctx xcontext.Context, req *HelloReq) (*HelloResp, error) {
	txt, err := s.srv.Hello(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &HelloResp{Text: txt}, nil
}
