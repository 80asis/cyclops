package main

import (
	// "github.com/80asis/cyclops/cyclopsAPIServer"
	// _ "github.com/80asis/cyclops/cyclopsMonitor"
	// "github.com/80asis/cyclops/cyclopsRPCServer"
	// log "github.com/sirupsen/logrus"
	"fmt"
	"sync"
	// "time"
	"github.com/80asis/cyclops/manager"
	"github.com/80asis/cyclops/entity"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	Entities := []entity.Entity{
        {EntityID: "entity1", EntityKind: "protection_rule", OpType: "update"},
        {EntityID: "entity2", EntityKind: "recovery_plan", OpType: "delete"},
		{EntityID: "entity3", EntityKind: "category", OpType: "update"},
    }
	targetAZ := []string{}
    Workflow := manager.GenericUpdates
    ForceSync := false

    // Create an instance of EntitySyncManager with custom values
    entitySyncManager := manager.NewEntitySyncManager(Entities, targetAZ, Workflow, ForceSync)

	go func() {
		defer wg.Done()

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
