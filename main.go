package main

import (
	"github.com/80asis/cyclops/cyclopsAPIServer"
	"github.com/80asis/cyclops/cyclopsMonitor"
	"github.com/80asis/cyclops/cyclopsRPCServer"
	log "github.com/sirupsen/logrus"
)

func main() {

	// starting go monitor thread
	go cyclopsMonitor.Start()
	// starting go RPC thread
	go cyclopsRPCServer.Start()
	// starting API Server
	cyclopsService := cyclopsAPIServer.NewService()
	handler := cyclopsAPIServer.NewHandler(cyclopsService)
	if err := handler.Serve(); err != nil {
		log.Error("Failed to serve the application")
	}
}
