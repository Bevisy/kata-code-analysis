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
	client := hello.NewHelloServiceClient(ttrpc.NewClient(conn))
	serverResponse, err := client.HelloWorld(context.Background(), &hello.HelloRequest{
		Msg: "Hello Server",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, serverResponse.Response)
}
