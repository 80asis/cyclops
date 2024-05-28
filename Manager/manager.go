package manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/80asis/cyclops/entity"
	"github.com/80asis/cyclops/tables"
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
func NewEntitySyncManager(entities []entity.Entity, targetAZ []string, workflow WorkflowType, forceSync bool) *EntitySyncManager {
	return &EntitySyncManager{
        Timestamp:       time.Now(),
        Entities:        entities,
        targetAZ:        targetAZ,
        workflow:        workflow,
        forceSync:       forceSync,
        UtilsManagerInf: &Utils{},
    }
}

func (manager *EntitySyncManager) InitializeTables() {
	manager.LocalTable = &tables.LocalTable{}
	manager.RemoteTable = &tables.RemoteTable{}
	fmt.Println("Local and remote tables initialized")
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
