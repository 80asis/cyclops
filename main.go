package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/80asis/cyclops/cyclops"
	"google.golang.org/grpc"
)

type LocalCyclopsServer struct {
	cyclops.UnimplementedCyclopsServer
}

func (localServer LocalCyclopsServer) SyncEntity(ctx context.Context, request *cyclops.CyclopsRequest) (*cyclops.CyclopsResponse, error) {
	fmt.Println("Got a request for ", request.EntityName)
	return &cyclops.CyclopsResponse{
		Status: true,
		Error:  "",
	}, nil
}

func main() {
	listener, error := net.Listen("tcp", ":8089")
	if error != nil {
		log.Fatalf("Failed to create a listner. Error: %s", error)
	}

	serverRgistrar := grpc.NewServer()
	service := &LocalCyclopsServer{}

	cyclops.RegisterCyclopsServer(serverRgistrar, service)
	error = serverRgistrar.Serve(listener)
	if error != nil {
		log.Fatalf("Failed to create a server. Error: %s", error)
	}

}
