.PHONY: generate server client

generate:
	protoc -I pb/ --gogottrpc_out=plugins=ttrpc:pb pb/*.proto

server:
	go run server/server.go
client:
	go run client/ttrpc/client.go