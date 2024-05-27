package manager

import (
	"github.com/80asis/cyclops/entity"
)

type ManagerUtilsInf interface {
	CalculateChecksum(entityID entity.Entity) string
	CreateTriggerEntitySyncAZTask(entity entity.Entity, az string, workflowType WorkflowType) string
	CreateParentEntitySyncTask(childEntities []entity.Entity, az string, workflowType WorkflowType) string
}

type Utils struct{}

func (um *Utils) CalculateChecksum(entity entity.Entity) string {
	// Implement checksum calculation logic
	return entity.EntityID + "_dummy_checksum"
}

func (um *Utils) CreateTriggerEntitySyncAZTask(entity entity.Entity, az string, workflowType WorkflowType) string {
	// Implement logic to create ErgonTask
	return entity.EntityID + "_taskUUID"
}

func (um *Utils) CreateParentEntitySyncTask(childEntities []entity.Entity, az string, workflowType WorkflowType) string {
	return az + "_parent_taskUUID"
}

func (um *Utils) PollEntitySyncStatus(entityID entity.Entity) string {
	// Fetching task from DB that corresponds to this entity
	// Making call to ergon to check the task status
	status := "COMPLETE_SUCCESS"
	return status
}
