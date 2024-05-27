package manager

import (
	"errors"
	"time"
    "github.com/80asis/cyclops/entity"
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

// Base Struct
type EntitySyncManager struct {
	Timestamp time.Time
	Entities  []entity.Entity
	workflow  WorkflowType
}

func NewEntitySyncManager(timestamp time.Time, entities []entity.Entity, workflow WorkflowType) *EntitySyncManager {
	return &EntitySyncManager{
		Timestamp: timestamp,
		Entities:  entities,
		workflow:  workflow,
	}
}

func (manager *EntitySyncManager) Process() (map[string]map[string]string, error) {
	var subManager EntitySyncManagerInterface

	switch manager.workflow {
	case GenericUpdates:
		subManager = NewGenericSyncSubManager(manager, false)
	case AZPairing:
		subManager = NewAZPairingSubManager(manager)
	case PolicyEnablement:
		subManager = NewPolicyEnablementSubManager(manager, false)
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
