package manager

import (
	"fmt"
	"time"
	"github.com/80asis/cyclops/tables"
	"github.com/80asis/cyclops/entity"
)

type AZPairingSubManager struct {
	Timestamp        time.Time
	Entities         []entity.Entity
	workflow         WorkflowType
	esm              *EntitySyncManager // Inheritance
	UtilsManagerInf  ManagerUtilsInf
	LocalTable       tables.LocalTableInterface
	RemoteTable      tables.RemoteTableInterface
}

func NewAZPairingSubManager(esm *EntitySyncManager) *AZPairingSubManager {
	return &AZPairingSubManager{
		esm:             esm,
		Timestamp:       esm.Timestamp,
		Entities:        esm.Entities,
		workflow:        esm.workflow,
		UtilsManagerInf: &Utils{},
	}
}

func (manager *AZPairingSubManager) InitializeTables() {
	manager.LocalTable = &tables.LocalTable{}
	manager.RemoteTable = &tables.RemoteTable{}
	fmt.Println("Local and remote tables initialized")
}

func (manager *AZPairingSubManager) FilterEntities() map[string][]entity.Entity {
	// Checksum check not required
	connectedAZs := manager.LocalTable.FetchConnectedAZs()
	entitiesToSync := manager.LocalTable.FetchAllEntities()

	entitiesToSyncByAZ := make(map[string][]entity.Entity)
	for _, az := range connectedAZs {
		entitiesToSyncByAZ[az] = entitiesToSync
	}
	return entitiesToSyncByAZ
}

func (manager *AZPairingSubManager) CreateErgonTasksForEntitySync(entitiesToSyncByAZ map[string][]entity.Entity) map[string]map[string]string {
	azToEntityTaskMap := make(map[string]map[string]string)
	for az, entities := range entitiesToSyncByAZ {

		parentTask := manager.UtilsManagerInf.CreateParentEntitySyncTask(entities, az, manager.workflow)
		fmt.Printf("Created %s", parentTask)

		entityToTaskMap := make(map[string]string)
		for _, entity := range entities {
			task := manager.UtilsManagerInf.CreateTriggerEntitySyncAZTask(entity, az, manager.workflow)
			time.Sleep(1 * time.Second)
			fmt.Printf("ErgonTask: %s created for entity: %+v, AZ: %s\n", task, entity, az)
			entityToTaskMap[entity.EntityID] = task
		}
		azToEntityTaskMap[az] = entityToTaskMap
	}
	return azToEntityTaskMap
}

func (manager *AZPairingSubManager) Start() (map[string]map[string]string, error) {
	manager.InitializeTables()
	entitiesToSyncByAZ := manager.FilterEntities()
	azToEntityTaskMap := manager.CreateErgonTasksForEntitySync(entitiesToSyncByAZ)
	return azToEntityTaskMap, nil
}
