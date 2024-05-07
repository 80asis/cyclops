package main

import (
	"log"
	"net"

	"cyclopsWorkflows"

	"github.com/80asis/cyclops/cyclops"
	"google.golang.org/grpc"
)

func main() {
	listener, error := net.Listen("tcp", ":8089")
	if error != nil {
		log.Fatalf("Failed to create a listner. Error: %s", error)
	}
	cyclopsWorkflows.BaseCyclopsTask.Run()
	serverRgistrar := grpc.NewServer()
	service := &LocalCyclopsRpcSvcServer{}
	cyclops.RegisterCyclopsRpcSvcServer(serverRgistrar, service)
	error = serverRgistrar.Serve(listener)
	if error != nil {
		log.Fatalf("Failed to create a server. Error: %s", error)
	}

}
