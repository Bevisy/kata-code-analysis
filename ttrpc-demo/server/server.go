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
	defer s.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	hello.RegisterHelloServiceService(s, &helloService{})
	if err := s.Serve(context.Background(), lis); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

type helloService struct{}

func (s helloService) HelloWorld(ctx context.Context, r *hello.HelloRequest) (*hello.HelloResponse, error) {
	if r.Msg == "" {
		return nil, errors.New("ErrNoInputMsgGiven")
	}
	return &hello.HelloResponse{Response: "Hi,ttrpc client"}, nil
}
