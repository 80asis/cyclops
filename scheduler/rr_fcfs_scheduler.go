package scheduler

import (
	"sync"

	t "github.com/80asis/cyclops/Task"
	log "github.com/sirupsen/logrus"
)

// AzTasks contains all the queued tasks
type AzTasks struct {
	AzPairTasks           t.Tasks
	PolicyEnablementTasks t.Tasks
	IndividualTasks       t.Tasks
}

type Quota struct {
	FreeSlots     int
	OccupiedSlots int
	TotalSlots    int
}

type scheduler interface {
	gethQuota()
	gethTasks()
	buildQuota()
	buildTasks()
	start()
	runTask()
}

type SchedulerStruct struct {
	/*
		quota = {
			"AZ-1": {
				"freeSlots": 20
				"occupiedSlots": 30
				"totalSlots": 50
			}
		}
	*/
	Quota map[string]Quota
	/*
		tasks = {
			"AZ-1": {
				"azPair": ["azPairTask-1", "azPairTask-2"]
				"policyEnablement": ["PETask-1","PETask-2"]
				"individualTasks": ["forceSync-1", "indvidual-update-2"]
			}
		}
	*/
	Tasks map[string]AzTasks
}

var wg *sync.WaitGroup

func init() {
	schd := SchedulerStruct{}
	schd.Quota = schd.buildQuota()
	schd.Tasks = schd.buildTasks()
	schd.start()
}

func (s *SchedulerStruct) getQuota() map[string]Quota {
	return s.Quota
}
func (s *SchedulerStruct) getTasks() map[string]AzTasks {
	return s.Tasks
}

func (s *SchedulerStruct) buildTasks() map[string]AzTasks {
	// Fetches all the running and queued tasks registered in Ergon and create a list in the memory
	// these task shall be put to running ASAP.
	// first we fetch all the running tasks and enqueue them
	// second it featch all the queued tasks and enqueue them
	// This comes active only on a restart
	return map[string]AzTasks{}
}

func (s *SchedulerStruct) buildQuota() map[string]Quota {
	// Build Quota will check the AZ table in IDF to identify how many AZ do we have
	// Then it will also consider the localAZ which shall have equal share of slots
	return map[string]Quota{}
}

func (s *SchedulerStruct) start() {
	// POOL OF GO ROUTINES WILL BE ALLOWED TO RUN FROM HERE
	// THIS WILL BE THE MINI-SCHEDULER THREAD RESPONSIBLE FOR EACH AZ
	// EACH AZ'S RUN() WILL RUN THE TASKS DEPENDING UPON THE Quota PROVIDED IN THE ROUND ROBIN FASHION
	/*
		1. Fetch the Quota of the AZ and Tasks (Resource)
		3. Allocate the tasks in a round robin fashion in a loop from i=0 to i=quota.FreeSlots (Identify demand and allocate)
		4. Updates s.Quota and s.Tasks
	*/
	var quota map[string]Quota = s.getQuota()
	for AZName, _ := range quota {
		go s.runAZTasks(AZName)
	}
}

func (s *SchedulerStruct) runAZTasks(AZName string) {
	var q Quota = s.getQuota()[AZName]
	log.Infof("Running tasks for PC: %v", AZName)
	for {
		if q.FreeSlots < 1 {
			// wait
		} else {
			var currentTask AzTasks = s.getTasks()[AZName]
			var currentTaskType string = "AZ-Pair-Task"
			var task t.Task
			for i := 0; i < q.FreeSlots; i++ {
				switch currentTaskType {
				case "AZ-Pair-Task":
					currentTaskType = "Policy-Enablement-Task"
					task = currentTask.AzPairTasks.Dequeue()
				case "Policy-Enablement-Task":
					currentTaskType = "Individual-Task"
					task = currentTask.PolicyEnablementTasks.Dequeue()
				case "Individual-Task":
					currentTaskType = "AZ-Pair-Task"
					task = currentTask.IndividualTasks.Dequeue()
				}
				if task.Uuid != nil {
					s.runTask(task)
				}
			}
		}
	}

}

func (s *SchedulerStruct) addTask(task t.Task) {
	// add the tasks to the queue
}
func (s *SchedulerStruct) runTask(task t.Task) {
	// Runs the tasks
}
