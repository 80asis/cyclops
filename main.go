package main

import (
	"sync"

	"github.com/80asis/cyclops/cyclopsAPIServer"
	_ "github.com/80asis/cyclops/cyclopsMonitor"
	"github.com/80asis/cyclops/cyclopsRPCServer"
	log "github.com/sirupsen/logrus"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go cyclopsRPCServer.Start()

	// starting API Server
	cyclopsService := cyclopsAPIServer.NewService()
	handler := cyclopsAPIServer.NewHandler(cyclopsService)
	if err := handler.Serve(); err != nil {
		log.Error("Failed to serve the application")
	}

}
