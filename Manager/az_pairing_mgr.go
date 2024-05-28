package manager

import (
	"fmt"
	"time"
	"github.com/80asis/cyclops/entity"
)

type AZPairingSubManager struct {
	esm *EntitySyncManager
}

func NewAZPairingSubManager(esm *EntitySyncManager) *AZPairingSubManager {
	return &AZPairingSubManager{
		esm: esm,
	}
}

func (manager *AZPairingSubManager) FilterEntities() map[string][]entity.Entity {
	// Checksum check not required
	connectedAZs := manager.esm.LocalTable.FetchConnectedAZs()
	if len(manager.esm.targetAZ) > 0 {
		connectedAZs = manager.esm.targetAZ
	}
	entitiesToSync := manager.esm.LocalTable.FetchAllEntities()

	entitiesToSyncByAZ := make(map[string][]entity.Entity)
	for _, az := range connectedAZs {
		entitiesToSyncByAZ[az] = entitiesToSync
	}
	return entitiesToSyncByAZ
}

func (manager *AZPairingSubManager) CreateErgonTasksForEntitySync(entitiesToSyncByAZ map[string][]entity.Entity) map[string]map[string]string {
	azToEntityTaskMap := make(map[string]map[string]string)
	for az, entities := range entitiesToSyncByAZ {

		parentTask := manager.esm.UtilsManagerInf.CreateParentEntitySyncTask(entities, az, manager.esm.workflow)
		fmt.Printf("Created %s", parentTask)

		entityToTaskMap := make(map[string]string)
		for _, entity := range entities {
			task := manager.esm.UtilsManagerInf.CreateTriggerEntitySyncAZTask(entity, az, manager.esm.workflow)
			time.Sleep(1 * time.Second)
			fmt.Printf("ErgonTask: %s created for entity: %+v, AZ: %s\n", task, entity, az)
			entityToTaskMap[entity.EntityID] = task
		}
		azToEntityTaskMap[az] = entityToTaskMap
	}
	return azToEntityTaskMap
}

func (manager *AZPairingSubManager) Start() (map[string]map[string]string, error) {
	entitiesToSyncByAZ := manager.FilterEntities()
	azToEntityTaskMap := manager.CreateErgonTasksForEntitySync(entitiesToSyncByAZ)
	return azToEntityTaskMap, nil
}
