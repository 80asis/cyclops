package cyclopsRPCServer

import (
	"net"

	"github.com/80asis/cyclops/cyclops"
	"github.com/80asis/cyclops/cyclopsMonitor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start() {

	listener, error := net.Listen("tcp", ":8089")
	if error != nil {
		log.Fatalf("Failed to create a listner. Error: %s", error)
	}
	localMonitor := cyclopsMonitor.GetNewMonitor()
	// starting go rpc server thread
	serverRgistrar := grpc.NewServer()
	service := &LocalCyclopsRpcSvcServer{
		Monitor: localMonitor,
	}
	cyclops.RegisterCyclopsRpcSvcServer(serverRgistrar, service)
	error = serverRgistrar.Serve(listener)
	if error != nil {
		log.Fatalf("Failed to create a server. Error: %s", error)
	}
}
