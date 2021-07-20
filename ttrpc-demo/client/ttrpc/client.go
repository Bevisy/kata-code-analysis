package main

import (
	"context"
	"fmt"
	"net"
	"os"

	hello "github.com/bevisy/kata-code-analysis/ttrpc-demo/pb"

	"github.com/containerd/ttrpc"
)

const port = ":9000"

func main() {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to dial: %v \n", err)
		os.Exit(1)
	}
	client := hello.NewGreetingServiceClient(ttrpc.NewClient(conn))
	serverResponse, err := client.Greeting(context.Background(), &hello.HelloRequest{
		Msg: "Hi greeting server",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, serverResponse.Response)
}
