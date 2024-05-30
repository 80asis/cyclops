package Monitor

import (
	"time"

	log "github.com/sirupsen/logrus"
)

/*
Responsibilities of a monitor Monitor is a daemon running in the background.
Monitor Invokes Manager when required.

1. Takes in a api request from API Handler
2. Takes in a rpc request from RPC Handler --> done
These API and RPC Requests can be as below
	1. Force Sync Request --> done/API
	2. Policy Enablement Request --> done/API
	3. Execute Cyclops Request --> done/RPC
	3. Non IDF Based Changes. <---This will be implemented on phase 2

These requests are then desearlized and forwarded to Manager for task creation.
Here we need to define what will be the payload but we can do that later once the flow is finalized by me.

Monitor will interact with CPDBclientAdapter to listen any changes on IDF.
These changes includes.
	1. AZ Pair change --> done/Internal via IDF Callback directly comes as callbacks back to monitor
	2. Individual change --> done/Internal via IDF Callback directly comes as callbacks back to monitor
    3. Internal Cascading requests --> done. Will add request when on ExecuteCyclops Workflow./ Internal via IDF Callback directly comes as callbacks back to monitor
*/
type InsightAdapter interface {
	Register()
	Unregister()
	RegisterClient()
	ErrorCb()
	RegisterAllWatch()
	StartNewWatch(entity_type, callback func())
	StartWatch()
	StopWatch()
	NewEntity()
	UpdateEntityCb()
	DeleteEntityCb()
	Reset()
}

var InsightClient InsightAdapter

var RESET int = 1
var READY int = 2

type InsightClientAdapter struct {
	State int
	// DBclient *insights_interface.InsightsService
	// {str: [cb]}
	// Dictionary mapping entity_type_name to a list of registered callbacks.
	// This is used to keep all the entity_type whose callbacks returns just return one version of change
	EntityTypeWatches map[string]interface{}
	// {str: [cb]}
	// Dictionary mapping entity_type_name to a list of registered callbacks.
	// This dictionary tracks update watches registered on entity type.
	// This is used to keep all the entity_type whose callbacks returns just return two version of change
	// 1 for old and other for new version
	EntityTypeUpdateWatches map[string]interface{}
	EntityTypeDeleteWatches map[string]interface{}
}

// Starting the InsightClientAdapter Thread
func InsightClientAdapaterInit() {
	// starts the InsightClientAdapter thread and register and unregister threads
	defer panicRecover()
	log.Info("Initiating Insight Client Adapter threads")
	client := GetClientAdapter()
	log.Info(client)
	go client.Register()
	go client.Unregister()
}

// Gets the InsightClientAdapter
func GetClientAdapter() InsightAdapter { //-> interface{
	// This is a method to get InsightClientAdapter type
	// we are keep InsightClientAdapter type as singleton type

	log.Info("Getting Adapter")
	if InsightClient != nil {
		log.Info("Returning existing InsightClientAdapter")
		return InsightClient
	}
	log.Info("No InsightClientAdapter found. Creating...")
	InsightClient := &InsightClientAdapter{
		State:                   RESET,
		EntityTypeWatches:       make(map[string]interface{}),
		EntityTypeUpdateWatches: make(map[string]interface{}),
		EntityTypeDeleteWatches: make(map[string]interface{}),
	}
	return InsightClient
}

// Registers all the entity watch to the IDF in case of reset or restart
func (c *InsightClientAdapter) Register() {
	defer wg.Done()
	defer panicRecover()
	c.RegisterClient()
	log.Info("Registering New Entities")
	log.Info("Registering client to IDF")

	for true {
		if c.State == READY {
			// keep on waiting for callbacks
			time.Sleep(5 * time.Second)
		}
		c.RegisterClient()
		c.StopWatch()
		c.RegisterAllWatch()
		c.StartWatch()
	}

}

func (c *InsightClientAdapter) Unregister() {
	defer wg.Done()
	defer panicRecover()
	log.Info("Unregistering Inactive Entities")
	for true {
		// removes all the inactive watches
	}

}

func (c *InsightClientAdapter) RegisterClient() {
	// Connects to IDF Client
	// c.DBclient = &insights_interface.InsightsService("client_id", c.ErrorCb)
}

// Callback invoked by InsightsWatchClient to reset watch client.
func (c *InsightClientAdapter) ErrorCb() {
	c.Reset()
}

// Register all the exisiting watches in the IDF.
func (c *InsightClientAdapter) RegisterAllWatch() {
	// log.Info("Regsitering watches to the IDF")
	// c.EntityTypeWatches
	// c.EntityTypeUpdateWatches
	// c.EntityTypeDeleteWatches
}

// registers a new watch with callback
func (c *InsightClientAdapter) StartNewWatch(entity_type, callback func()) {}

// Starts the watch on IDF of registered entities.
func (c *InsightClientAdapter) StartWatch() {
	// log.Info("Start Watch")
}

// This stops the watch on IDF of registered entity
func (c *InsightClientAdapter) StopWatch() {

	// log.Info("Stop Watch")
}

// callbacks
// we need to find out the protos used here
func (c *InsightClientAdapter) NewEntity() {
	// This callback should be triggred in case of any table on registered IDF entity
	log.Info("Creation on IDF detected.")
}

// This callback should be triggred in case of updation on registered IDF entity
func (c *InsightClientAdapter) UpdateEntityCb() {
	log.Info("Update on IDF detected.")
}

// This callback should be triggred in case of deletion on registered IDF entity
func (c *InsightClientAdapter) DeleteEntityCb() {
	log.Info("Deletion on IDF detected.")
}

// This method should stop all the watches and resgiter the watches back
// this is invoked only when there is a error is received from IDF
func (c *InsightClientAdapter) Reset() {
	log.Error("Error reported by IDF")
	c.State = RESET
}

/* Notes:

Go magneto people are using to return interface instead of struct
2. you can use struct as class, use init() and also use method sets
so when you return interface type example, Monitor interface but actually return the CyclopsMonitor struct
you have the access to all the methodsets.
follow to create the grpc, then depending methods for grpc A. Interface B. Independent threads C. Workflows.

*/
