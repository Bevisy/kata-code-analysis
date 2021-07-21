package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"

	hello "github.com/bevisy/kata-code-analysis/ttrpc-demo/pb"
	"github.com/containerd/ttrpc"
)

const port = ":9000"

func main() {

	s, err := ttrpc.NewServer()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer s.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// 注册 service 的实现
	hello.RegisterGreetingServiceService(s, &greetingService{})
	if err := s.Serve(context.Background(), lis); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

type greetingService struct{}

func (s greetingService) Greeting(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	if req.Msg == "" {
		return nil, errors.New("ErrNoInputMsgGiven")
	}

	fmt.Fprintln(os.Stdout, "received msg: ", req.Msg)
	return &hello.HelloResponse{Response: "Hi, client."}, nil
}

func (s greetingService) Bye(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	if req.Msg == "" {
		return nil, errors.New("ErrNoInputMsgGiven")
	}

	fmt.Fprintln(os.Stdout, "received msg: ", req.Msg)
	return &hello.HelloResponse{Response: "Bye, client."}, nil
}
