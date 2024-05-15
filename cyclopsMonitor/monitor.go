package cyclopsMonitor

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

/*
Responsibilities of a monitor Monitor is a daemon running in the background.
Monitor Invokes Manager when required.

1. Takes in a api request from API Handler
2. Takes in a rpc request from RPC Handler --> done
These API and RPC Requests can be as below
	1. Force Sync Request --> done
	2. Policy Enablement Request
	3. Execute Cyclops Request --> done
	3. Non IDF Based Changes. <---This will be implemented on phase 2

These requests are then desearlized and forwarded to Manager for task creation.
Here we need to define what will be the payload but we can do that later once the flow is finalized by me.

Monitor will interact with CPDBclientAdapter to listen any changes on IDF.
These changes includes.
	1. AZ Pair change --> done
	2. Individual change --> done
Internal Cascading requests --> done. Will add request when on ExecuteCyclops Workflow.
*/

type Monitor interface {
	register()       // Registering Entity
	registerClient() // Registering Client
	unRegister()     // Unregisters the entity
	startWatch()     // start a watch on IDF
	stopWatch()      // stop a watch on IDF
	createEntityCb() // callbacks for Entity Creation
	updateEntityCb() // callbacks for Entity Updation
	deleteEntityCb() // callbacks for Entity Deletion
}

var Client *CyclopsMonitor

var wg sync.WaitGroup

type CyclopsMonitor struct {
	// specific client level details for IDF
	// maintains connection
	clientIp string
}

func GetNewMonitor() *CyclopsMonitor { //-> interface{
	// make sures the monitor is singleton and creates one if not created.
	log.Info("Getting monitor")
	if Client != nil {
		log.Info("Returning existing monitor")
		return Client
	}
	log.Info("No monitor found. Creating a new monitor")
	return &CyclopsMonitor{
		clientIp: "0.0.0.0:9090",
	}
}
func panicRecover() {
	// generic utility method to capture painc
	if err := recover(); err != nil {
		fmt.Println("Panik Panik!! Error: ", err)
	}
}

// Starting the Monitor Thread
func Start() {
	// starts the monitor thread and register and unregister threads
	defer panicRecover()
	log.Info("Initiating threads")
	wg.Add(2)
	client := GetNewMonitor()
	log.Info(client)
	go client.register()
	go client.unRegister()
	wg.Wait()
}

// For adding entitis to Manager
func (c *CyclopsMonitor) AddNotification(jsonData []byte) {
	// this submits entity data to manager
	defer panicRecover()
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Info("Error desearlizing data")
		return
	}
	log.Info("Adding Entity to Manager for task creation. Data: ", data)
}

// Registering Entity
func (c *CyclopsMonitor) register() {
	// Registers the entity watch to IDF
	defer wg.Done()
	defer panicRecover()
	c.registerClient()
	log.Info("Registering New Entities")
	for true {

		time.Sleep(5 * time.Second)
	}
}

// Register Client
func (c *CyclopsMonitor) registerClient() {
	log.Info("Registering client to IDF")
}

// Unregister Entity
func (c *CyclopsMonitor) unRegister() {
	// unresgiter entity watch from IDF

	defer wg.Done()
	defer panicRecover()
	log.Info("Removing entity from IDF as its inactive")
	for true {
		time.Sleep(3 * time.Second)
	}
}

// watch cycle
func (c *CyclopsMonitor) startWatch() {
	// this basically starts the watch on IDF of registered entities.
	log.Info("Start Watch")
}
func (c *CyclopsMonitor) stopWatch() {
	// This stops the watch on IDF of registered entity
	log.Info("Stop Watch")
}

// callbacks
func (c *CyclopsMonitor) createEntityCb() {
	// This callback should be triggred in case of creation on registered IDF entity
	log.Info("Creation on IDF detected.")
}
func (c *CyclopsMonitor) updateEntityCb() {
	// This callback should be triggred in case of updation on registered IDF entity
	log.Info("Update on IDF detected.")
}

func (c *CyclopsMonitor) deleteEntityCb() {
	// This callback should be triggred in case of deletion on registered IDF entity
	log.Info("Deletion on IDF detected.")
}

// Go magneto people are using to return interface instead of struct
// 2. you can use struct as class, use init() and also use method sets
// so when you return interface type example, Monitor interface but actually return the CyclopsMonitor struct
// you have the access to all the methodsets.
// follow to create the grpc, then depending methods for grpc A. Interface B. Independent threads C. Workflows.
