package manager

import (
	"fmt"
	"time"
	"github.com/80asis/cyclops/tables"
	"github.com/80asis/cyclops/entity"
)

type CascadingSubManager struct {
	Timestamp       time.Time
	Entities        []entity.Entity
	workflow        WorkflowType
	forceSync       bool
	esm             *EntitySyncManager // Inheritance
	UtilsManagerInf ManagerUtilsInf
	LocalTable      tables.LocalTableInterface
	RemoteTable     tables.RemoteTableInterface
}

func NewCascadingSubManager(esm *EntitySyncManager) *CascadingSubManager {
	return &CascadingSubManager{
		esm:             esm,
		Timestamp:       esm.Timestamp,
		Entities:        esm.Entities,
		workflow:        esm.workflow,
		UtilsManagerInf: &Utils{},
	}
}

func (manager *CascadingSubManager) InitializeTables() {
	manager.LocalTable = &tables.LocalTable{}
	manager.RemoteTable = &tables.RemoteTable{}
	fmt.Println("Local and remote tables initialized")
}

func (manager *CascadingSubManager) FilterEntities() map[string][]entity.Entity {
	connectedAZs := manager.LocalTable.FetchConnectedAZs()
	checksums := make(map[string]string)
	for _, entity := range manager.Entities {
		checksum := manager.UtilsManagerInf.CalculateChecksum(entity)
		checksums[entity.EntityID] = checksum
	}
	entitiesToSync := make(map[string]string)
	for _, entity := range manager.Entities {
		if manager.LocalTable.VerifyChecksum(entity.EntityID, checksums[entity.EntityID]) {
			continue
		}
		manager.LocalTable.UpdateChecksum(entity.EntityID, checksums[entity.EntityID])
		manager.LocalTable.MarkOutOfSync(entity.EntityID)
		entitiesToSync[entity.EntityID] = checksums[entity.EntityID]
	}
	entitiesToSyncByAZ := make(map[string][]entity.Entity)
	for _, az := range connectedAZs {
		entitiesToSyncInAZ := make([]entity.Entity, 0)
		for entityID, checksum := range entitiesToSync {
			if manager.RemoteTable.VerifyChecksumInAZ(entityID, checksum, az) {
				continue
			}
			for _, entity := range manager.Entities {
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

func (manager *CascadingSubManager) CreateErgonTasksForEntitySync(entitiesToSyncByAZ map[string][]entity.Entity) map[string]map[string]string {
	azToEntityTaskMap := make(map[string]map[string]string)
	for az, entities := range entitiesToSyncByAZ {
		entityToTaskMap := make(map[string]string)
		for _, entity := range entities {
			task := manager.UtilsManagerInf.CreateTriggerEntitySyncAZTask(entity, az, manager.workflow)
			time.Sleep(1 * time.Second)
			fmt.Printf("\nErgonTask: %s created for entity: %+v, AZ: %s\n", task, entity, az)
			entityToTaskMap[entity.EntityID] = task
		}
		azToEntityTaskMap[az] = entityToTaskMap
	}
	return azToEntityTaskMap
}

func (manager *CascadingSubManager) Start() (map[string]map[string]string, error) {
	manager.InitializeTables()
	entitiesToSyncByAZ := manager.FilterEntities()
	azToEntityTaskMap := manager.CreateErgonTasksForEntitySync(entitiesToSyncByAZ)
	return azToEntityTaskMap, nil
}
