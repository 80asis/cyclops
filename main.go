package main

import (
	"github.com/80asis/cyclops/cyclopsAPIServer"
	_ "github.com/80asis/cyclops/cyclopsMonitor"
	"github.com/80asis/cyclops/cyclopsRPCServer"
	log "github.com/sirupsen/logrus"
	"fmt"
	"sync"
	"time"
	"github.com/80asis/cyclops/manager"
	"github.com/80asis/cyclops/entity"
)

func main() {

	// starting go RPC thread
	go cyclopsRPCServer.Start()

	// starting API Server
	cyclopsService := cyclopsAPIServer.NewService()
	handler := cyclopsAPIServer.NewHandler(cyclopsService)
	if err := handler.Serve(); err != nil {
		log.Error("Failed to serve the application")
	}
	// Use WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Increment the WaitGroup counter
	wg.Add(1)

	// Create an instance of EntitySyncManager
	entitySyncManager := manager.NewEntitySyncManager(
		time.Now(),
		[]entity.Entity{
			{EntityID: "entity1", EntityKind: "protection_rule", OpType: "update"},
			{EntityID: "entity2", EntityKind: "recovery_plan", OpType: "delete"},
			{EntityID: "entity3", EntityKind: "category", OpType: "update"},
		},
		manager.GenericUpdates,
	)

	// Start a goroutine to execute the Process function
	go func() {
		defer wg.Done()

		// Call the Process function
		result, err := entitySyncManager.Process()
		if err != nil {
			fmt.Printf("Error processing payload: %v\n", err)
			return
		}

		// Handle the result
		fmt.Println(result)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Print completion message
	fmt.Println("Entity synchronization process completed successfully")
}
