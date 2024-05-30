package RPCServer

import (
	"net"

	"github.com/80asis/cyclops/Monitor"
	"github.com/80asis/cyclops/cyclops"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start() {

	listener, error := net.Listen("tcp", ":8089")
	if error != nil {
		log.Fatalf("Failed to create a listner. Error: %s", error)
	}
	monitor := Monitor.GetMonitor()
	cyclopsMonitor := monitor.(*Monitor.CyclopsMonitor)

	// starting go rpc server thread
	serverRgistrar := grpc.NewServer()
	service := &LocalCyclopsRpcSvcServer{
		Monitor: cyclopsMonitor,
	}
	cyclops.RegisterCyclopsRpcSvcServer(serverRgistrar, service)
	error = serverRgistrar.Serve(listener)
	if error != nil {
		log.Fatalf("Failed to create a server. Error: %s", error)
	}
}
