package manager

import (
	"errors"
	"sync"
	"time"

	"github.com/80asis/cyclops/entity"
	"github.com/80asis/cyclops/tables"
	log "github.com/sirupsen/logrus"
)

type WorkflowType string

const (
	GenericUpdates   WorkflowType = "GenericUpdates"
	AZPairing        WorkflowType = "AZPairing"
	PolicyEnablement WorkflowType = "PolicyEnablement"
	Cascading        WorkflowType = "Cascading"
)

type EntitySyncManagerInterface interface {
	FilterEntities() map[string][]entity.Entity
	CreateErgonTasksForEntitySync(entitiesToSyncByAZ map[string][]entity.Entity) map[string]map[string]string
	Start() (map[string]map[string]string, error)
}

type EntitySyncManager struct {
	Timestamp       time.Time
	Entities        []entity.Entity
	targetAZ        []string
	workflow        WorkflowType
	forceSync       bool
	UtilsManagerInf ManagerUtilsInf
	LocalTable      tables.LocalTableInterface
	RemoteTable     tables.RemoteTableInterface
}

// NewEntitySyncManager constructor function with default values
func NewEntitySyncManager() *EntitySyncManager {
	return &EntitySyncManager{
		Timestamp:       time.Now(),
		UtilsManagerInf: &Utils{},
	}
}

func (manager *EntitySyncManager) InitializeTables() {
	manager.LocalTable = &tables.LocalTable{}
	manager.RemoteTable = &tables.RemoteTable{}
	log.Info("Local and remote tables initialized")
}

func (manager *EntitySyncManager) Process() (map[string]map[string]string, error) {
	var subManager EntitySyncManagerInterface
	manager.InitializeTables()

	switch manager.workflow {
	case GenericUpdates:
		subManager = NewGenericSyncSubManager(manager)
	case AZPairing:
		subManager = NewAZPairingSubManager(manager)
	case PolicyEnablement:
		subManager = NewPolicyEnablementSubManager(manager)
	case Cascading:
		subManager = NewCascadingSubManager(manager)
	default:
		return nil, errors.New("unknown workflow type")
	}

	if subManager == nil {
		return nil, errors.New("failed to initialize submanager")
	}

	return subManager.Start()
}

func (manager *EntitySyncManager) AddEntity() {
	//TODO: We need to test this in case of multiple calls on ESM parallely.
	// I doubt that wg.Wait() can cause issue calling the AddEntity of the same instance

	// TODO: We need to collect the data from Monitor about entity uuid and type
	// do initial processing to indentify entity details and call the submanager

	// Dummy data: adds the entity to manager
	var wg sync.WaitGroup
	Entities := []entity.Entity{
		{EntityID: "entity1", EntityKind: "protection_rule", OpType: "update"},
		{EntityID: "entity2", EntityKind: "recovery_plan", OpType: "delete"},
		{EntityID: "entity3", EntityKind: "category", OpType: "update"},
	}
	targetAZ := []string{}
	Workflow := GenericUpdates
	ForceSync := false

	// Create an instance of EntitySyncManager with custom values
	manager.Entities = Entities
	manager.targetAZ = targetAZ
	manager.workflow = Workflow
	manager.forceSync = ForceSync
	wg.Add(1)
	go func() {
		defer wg.Done()

		result, err := manager.Process()
		if err != nil {
			log.Info("Error processing payload: %v\n", err)
			return
		}

		// Handle the result
		log.Info(result)
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Print completion message
	log.Info("Entity synchronization process completed successfully")
}
