package main

import (
	"sync"

	"github.com/80asis/cyclops/APIServer"
	_ "github.com/80asis/cyclops/Monitor"
	"github.com/80asis/cyclops/RPCServer"
	log "github.com/sirupsen/logrus"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	go RPCServer.Start()

	// starting API Server
	cyclopsService := APIServer.NewService()
	handler := APIServer.NewHandler(cyclopsService)
	if err := handler.Serve(); err != nil {
		log.Error("Failed to serve the application")
	}

}
