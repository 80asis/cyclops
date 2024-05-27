package manager

import (
	"fmt"
	"time"
	"github.com/80asis/cyclops/tables"
	"github.com/80asis/cyclops/entity"
)

type GenericSyncSubManager struct {
	Timestamp   time.Time
	Entities    []entity.Entity
	targetAZ    []string
	workflow    WorkflowType
	forceSync   bool
	esm         *EntitySyncManager // Inheritance
	UtilsManagerInf ManagerUtilsInf
	LocalTable      tables.LocalTableInterface
	RemoteTable     tables.RemoteTableInterface
}

func NewGenericSyncSubManager(esm *EntitySyncManager, forceSync bool) *GenericSyncSubManager {
	return &GenericSyncSubManager{
		esm:        esm,
		Timestamp:  esm.Timestamp,
		Entities:   esm.Entities,
		targetAZ:   esm.targetAZ,
		workflow:   esm.workflow,
		forceSync:  forceSync,
		UtilsManagerInf: &Utils{},
	}
}

func (manager *GenericSyncSubManager) InitializeTables() {
	manager.LocalTable = &tables.LocalTable{}
	manager.RemoteTable = &tables.RemoteTable{}
	fmt.Println("Local and remote tables initialized")
}

func (manager *GenericSyncSubManager) FilterEntities() map[string][]entity.Entity {
	connectedAZs := manager.LocalTable.FetchConnectedAZs()
	if len(manager.targetAZ) > 0{
		connectedAZs = manager.targetAZ
	}
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

func (manager *GenericSyncSubManager) CreateErgonTasksForEntitySync(entitiesToSyncByAZ map[string][]entity.Entity) map[string]map[string]string {
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

func (manager *GenericSyncSubManager) Start() (map[string]map[string]string, error) {
	manager.InitializeTables()
	entitiesToSyncByAZ := manager.FilterEntities()
	azToEntityTaskMap := manager.CreateErgonTasksForEntitySync(entitiesToSyncByAZ)
	return azToEntityTaskMap, nil
}
