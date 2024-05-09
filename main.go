package main

import (
	"github.com/80asis/cyclops/cyclopsMonitor"
	"github.com/80asis/cyclops/cyclopsRPCServer"
)

func main() {

	// starting go monitor thread
	go cyclopsMonitor.Start()
	// starting go RPC thread
	cyclopsRPCServer.Start()

}
