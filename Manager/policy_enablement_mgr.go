package manager

import (
	"fmt"
	"time"
	"github.com/80asis/cyclops/entity"
)

type PolicyEnablementSubManager struct {
	esm *EntitySyncManager
}

func NewPolicyEnablementSubManager(esm *EntitySyncManager) *PolicyEnablementSubManager {
	return &PolicyEnablementSubManager{
		esm: esm,
	}
}

func (manager *PolicyEnablementSubManager) FilterEntities() map[string][]entity.Entity {
	connectedAZs := manager.esm.LocalTable.FetchConnectedAZs()
	checksums := make(map[string]string)
	for _, entity := range manager.esm.Entities {
		checksum := manager.esm.UtilsManagerInf.CalculateChecksum(entity)
		checksums[entity.EntityID] = checksum
	}
	entitiesToSync := make(map[string]string)
	for _, entity := range manager.esm.Entities {
		if manager.esm.LocalTable.VerifyChecksum(entity.EntityID, checksums[entity.EntityID]) {
			continue
		}
		manager.esm.LocalTable.UpdateChecksum(entity.EntityID, checksums[entity.EntityID])
		manager.esm.LocalTable.MarkOutOfSync(entity.EntityID)
		entitiesToSync[entity.EntityID] = checksums[entity.EntityID]
	}
	entitiesToSyncByAZ := make(map[string][]entity.Entity)
	for _, az := range connectedAZs {
		entitiesToSyncInAZ := make([]entity.Entity, 0)
		for entityID, checksum := range entitiesToSync {
			if manager.esm.RemoteTable.VerifyChecksumInAZ(entityID, checksum, az) {
				continue
			}
			for _, entity := range manager.esm.Entities {
				if entity.EntityID == entityID {
					entitiesToSyncInAZ = append(entitiesToSyncInAZ, entity)
					break
				}
			}
		}
		entitiesToSyncByAZ[az] = entitiesToSyncInAZ
	}
	return entitiesToSyncByAZ
}

func (manager *PolicyEnablementSubManager) CreateErgonTasksForEntitySync(entitiesToSyncByAZ map[string][]entity.Entity) map[string]map[string]string {
	azToEntityTaskMap := make(map[string]map[string]string)
	for az, entities := range entitiesToSyncByAZ {
		entityToTaskMap := make(map[string]string)
		for _, entity := range entities {
			task := manager.esm.UtilsManagerInf.CreateTriggerEntitySyncAZTask(entity, az, manager.esm.workflow)
			time.Sleep(1 * time.Second)
			fmt.Printf("\nErgonTask: %s created for entity: %+v, AZ: %s\n", task, entity, az)
			entityToTaskMap[entity.EntityID] = task
		}
		azToEntityTaskMap[az] = entityToTaskMap
	}
	return azToEntityTaskMap
}

func (manager *PolicyEnablementSubManager) Start() (map[string]map[string]string, error) {
	entitiesToSyncByAZ := manager.FilterEntities()
	azToEntityTaskMap := manager.CreateErgonTasksForEntitySync(entitiesToSyncByAZ)
	return azToEntityTaskMap, nil
}
