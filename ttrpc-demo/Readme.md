# ttrpc-demo

## install protoc-gen-gogottrpc
```sh
go install github.com/containerd/ttrpc/cmd/protoc-gen-gogottrpc
```

## generate rpc service code
```sh
make generate
```

## run server and client alone
```sh
make server
make client
```
and you will get the result:
```sh
$ make client
go run client/ttrpc/client.go
Hi,ttrpc client
```