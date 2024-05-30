package Scheduler

import (
	"sync"

	t "github.com/80asis/cyclops/Task"
)

// AzTasks contains all the queued tasks

type PriorityQuota struct {
	FreeSlots        int `description:"Available slots for allocation"`
	OccupiedSlots    int `description:"Current Usage of the Quota"`
	SoftMaxCostQouta int `description:"The maximum resources a queue can use before it starts borrowing from other queues"`
	HighMaxCostQouta int `description:"The absolute limit of resources a queue can borrow, even if its a high priority queue."`
}

type MLFQQuota struct {
	LowPriorityTaskQuota  PriorityQuota
	MidPriorityTaskQuota  PriorityQuota
	HighPriorityTaskQuota PriorityQuota
}

type MlFQTasks struct {
	lowPriorityQueue  []t.Task // ["azPairTask"]
	midPriorityQueue  []t.Task // ["PolicyEnablementTask"]
	highPriorityQueue []t.Task // ["ForceSync", "indvidual-change-sync"]
}

type MLFQSchedulerInt interface {
	mlfq_gethQuota()
	mlfq_getTasks()
	mlfq_buildQuota()
	mlfq_buildTasks()
	mlfq_start()
	mlfq_runTask()
	mlfq_borrower()
}

type MLFQScheduler struct {
	/*
		quota = {
			"LowPriorityTaskQuota": {
				"freeSlots": 20
				"occupiedSlots": 30
				"SoftMaxCostQouta": 33% of total Cyclops Slots
				"HighMaxCostQouta": 20% of available quota from Mid/High PriorityQueue
			},
			"MidPriorityTaskQuota": {
				"freeSlots": 20
				"occupiedSlots": 30
				"SoftMaxCostQouta": 33% of total Cyclops Slots
				"HighMaxCostQouta": 20% of available quota from High/Low PriorityQueue
			},
			"HighPriorityTaskQuota": {
				"freeSlots": 20
				"occupiedSlots": 30
				"SoftMaxCostQouta": 33% of total Cyclops Slots
				"HighMaxCostQouta": 20% of available quota from Low/Mid PriorityQueue
			}
		}
	*/
	Quota MLFQQuota
	/*
		tasks = {
			"lowPriorityQueue": ["azPairTask"]
			"midPriorityQueue": ["PolicyEnablementTask"]
			"highPriorityQueue": ["ForceSync", "indvidual-change-sync"]
		}
	*/
	Tasks MlFQTasks
}

var MLFQwg *sync.WaitGroup

func init() {
	schd := MLFQScheduler{}
	schd.Quota = schd.mlfq_buildQuota()
	schd.Tasks = schd.mlfq_buildTasks()
	schd.mlfq_start()
}

func (s *MLFQScheduler) mlfq_getQuota() MLFQQuota {
	return s.Quota
}
func (s *MLFQScheduler) mlfq_getTasks() MlFQTasks {
	return s.Tasks
}

func (s *MLFQScheduler) mlfq_buildTasks() MlFQTasks {
	// Fetches all the running and queued tasks registered in Ergon and create a list in the memory
	// these task shall be put to running ASAP.
	// first we fetch all the running tasks and enqueue them
	// second it featch all the queued tasks and enqueue them
	// This comes active only on a restart
	return MlFQTasks{}
}

func (s *MLFQScheduler) mlfq_buildQuota() MLFQQuota {
	// First it gets the total quota it has from config.
	return MLFQQuota{}
}

func (s *MLFQScheduler) mlfq_start() {
	/*
		MLFQ works between 3 queues, high, mid and low priority queues
		low priority queue have az pairing tasks
		mid priority queue have policy enablement tasks
		high priority queue have inidvidual updates

		Each queue have there own quota which inturn has all the queued and running tasks only
		each queue has there respective go routine that runs concurrently depeneding upon how may free slots and queued task it has

		When an operation is pushed onto the queue, the Scheduler checks if there are sufficient resources to execute it immediately.
		If the resources are sufficient, the operation is executed right away (fastpath).
		If the resources are not sufficient, the operation is queued.
		The operation will wait in the queue until there are enough resources available for its execution.
	*/
	go s.runHighPriorityTasks()
	go s.runMidPriorityTasks()
	go s.runLowPriorityTasks()
	go s.mlfq_borrower()
}

func (s *MLFQScheduler) runHighPriorityTasks() {
	// This go routine runs the high priority tasks It works on FCFS basis.
	// Soft Max Cost Quota: If the cost of executing the operation is less than or equal to the soft maximum quota, the operation can be executed.
	// Hard Max Cost Quota: For high-priority operations, if the cost is within the hard maximum quota, the operation can still be excuted by borrowing resources (mlfq_borrower()).

	/*
		During Peek Operation:
			If a queue has exceeded its current resource quota. NO operation can be executed due to resource constraints,
			the run task may still admit an operation from the highest priority queue, borrowing resources from other queues if necessary.
	*/
	for {
		s.mlfq_runTask(t.Task{})
	}
}
func (s *MLFQScheduler) runMidPriorityTasks() {
	// This go routine runs the mid priority tasks It works on FCFS basis.
	// Soft Max Cost Quota: If the cost of executing the operation is less than or equal to the soft maximum quota, the operation can be executed.
	// Hard Max Cost Quota: For high-priority operations, if the cost is within the hard maximum quota, the operation can still be excuted by borrowing resources (mlfq_borrower()).

	/*
		During Peek Operation:
			If a queue has exceeded its current resource quota. NO operation can be executed due to resource constraints,
			the run task may still admit an operation from the highest priority queue, borrowing resources from other queues if necessary.
	*/
	for {
		s.mlfq_runTask(t.Task{})
	}

}
func (s *MLFQScheduler) runLowPriorityTasks() {
	// This go routine runs the low priority tasks It works on FCFS basis.
	// Soft Max Cost Quota: If the cost of executing the operation is less than or equal to the soft maximum quota, the operation can be executed.
	// Hard Max Cost Quota: For high-priority operations, if the cost is within the hard maximum quota, the operation can still be excuted by borrowing resources (mlfq_borrower()).

	/*
		During Peek Operation:
			If a queue has exceeded its current resource quota. NO operation can be executed due to resource constraints,
			the run task may still admit an operation from the highest priority queue, borrowing resources from other queues if necessary.
	*/
	for {
		s.mlfq_runTask(t.Task{})
	}
}
func (s *MLFQScheduler) mlfq_runTask(task t.Task) {
	// Runs the tasks
}

func (s *MLFQScheduler) mlfq_borrower() {
	/*
		.
		Borrowing ensures that even if a queue has exceeded its initial resource allocation, operations can still proceed by temporarily using resources allocated to other queues.
		This borrowing is bounded by the soft and hard max cost quotas to prevent excessive resource usage by any single queue.
		This method is to be called during peek time when the resources are already allocated and queue has tasks which are in queued state.
		If a queue has exceeded its current resource quota, it checks if the queue can borrow resources up to its soft or hard maximum quota.
		Soft Max Cost Quota: If the cost of executing the operation is less than or equal to the soft maximum quota, the operation can be executed.
		Hard Max Cost Quota: For high-priority operations, if the cost is within the hard maximum quota, the operation can still be executed.
	*/
}
