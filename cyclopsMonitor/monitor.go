package cyclopsMonitor

import (
	"encoding/json"
	"sync"

	"github.com/80asis/cyclops/cyclops"
	config "github.com/80asis/cyclops/cyclopsConfig"
	log "github.com/sirupsen/logrus"
)

type Monitor interface {
	AddNotification(jsonData []byte) // calls the manager and adds the changes to manager queue.
	AddNotificationForForceSync(entity_uuid []byte)
	CreateUpdateCb(entity_proto, old_entity_proto []byte)
	DeleteCb(entity_proto, old_entity_proto []byte)
	availabilityZoneAddCb(entity_proto []byte)
	availabilityZoneDeleteCb(entity_proto []byte)
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
	InsightsAdapter InsightAdapter
	// specific client level details for IDF
	// maintains connection
	clientIp string
}

func GetNewMonitor() *CyclopsMonitor { //-> interface{
	// make sures the monitor is singleton and creates one if not created.
	fmt.Println("Getting monitor")
	if Client != nil {
		fmt.Println("Returning existing monitor")
		return Client
	}
	// TODO: use a package to create singleton monitor
	fmt.Println("No monitor found. Creating a new monitor")
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
func init() {
	// starts the monitor thread and register all the plugins to idf
	// Here we will addiing watch for az tables for delete update or add entity
	// We will also add all the entities/plugins to be watches by idf
	//
	defer panicRecover()
	log.Info("Initiating Monitor Thread")
	wg.Add(2)
	client := GetNewMonitor()
	// Cast the Monitor interface to *CyclopsMonitor
	cyclopsClient := client.(*CyclopsMonitor)
	cyclopsClient.InsightsAdapter = GetClientAdapter()
	log.Info(client)
	for entity_type := range config.RegisterPlugins {
		// fetches all the plugin and registeres them to IDF
		// Example: client.InsightsAdapter.startNewWatch(plugin, client.addKindCreateCb)
		log.Infof("Regsitering %v", config.Entity_type_str_map[cyclops.EntitySyncEntityType_Type(entity_type)])
	}

	// we call this AddNotification in case of any restart. AddNotification if no data is given triggeres sync for all the OOS (out-of-sync) entities.
	client.AddNotification([]byte{})
	fmt.Println(client)

	//TODO: Check how to properly register IDF watches
	// Registering IDF watch
	go client.register()
	// UnRegistering IDF watch
	go client.unRegister()
	wg.Wait()
}

func GetNewMonitor() Monitor { //-> interface{
	// This is a method to get monitor type
	// we are keep monitor type as singleton type
	log.Info("Getting monitor")
	if Client != nil {
		log.Info("Returning existing monitor")
		return Client
	}
	log.Info("No monitor found. Creating a new monitor")
	return &CyclopsMonitor{}
}

func (c *CyclopsMonitor) AddNotification(jsonData []byte) {
	// This submits every change to the manager. AddNotification can be called by
	// 1. RPCs/API - Policy Enablement
	// 2. By Execute for adding entities to be synced for Cascading
	// 3. IDF Callbacks and NonIDF Callbacks
	// 4. AZ Pairing Request
	defer panicRecover()
	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Info("Error desearlizing data")
		return
	}
	log.Info("Adding Entity to Manager for task creation. Data: ", data)
}

func (c *CyclopsMonitor) AddNotificationForForceSync(entity_uuid []byte) {
	// this submits every change to manager
	// AddNotificationForForceSync can be called by
	// 1. RPCs/API - ForceSync
	// 2. By ExecuteSync for adding entities to be forceSynced in case of Cascading

	// Args:
	//	entity []byte: This is uuid of that entity that need to be forsynced to other PCs
}

func (c *CyclopsMonitor) CreateUpdateCb(entity_proto, old_entity_proto []byte) {
	// this is a generic method which will receive requests from idf in form of callbacks
	// this will issue a trigger to addnotification for the specific entity
	// Args:
	//	entity_proto: This is the new updated proto that will be shared by IDF
	//					Entity Type of proto is mentioned here https://sourcegraph.ntnxdpro.com/ntnxdb-master/-/blob/ntnxdb_client/insights/insights_interface/insights_interface.proto
	//  old_entity_proto: This is the old value of the proto that will be shared by IDF
	//					Entity Type of proto is mentioned here https://sourcegraph.ntnxdpro.com/ntnxdb-master/-/blob/ntnxdb_client/insights/insights_interface/insights_interface.proto
// Register Client
func (c *CyclopsMonitor) registerClient() {
	// Registering your client to IDF service
	fmt.Println("Registering client to IDF")
}

func (c *CyclopsMonitor) DeleteCb(entity_proto, old_entity_proto []byte) {
	// this is a generic method which will receive requests from idf in form of callbacks
	// this will issue a trigger to addnotification for the specific entity
	// Args:
	//	entity_proto: This is the new updated proto that will be shared by IDF
	//					Entity Type of proto is mentioned here https://sourcegraph.ntnxdpro.com/ntnxdb-master/-/blob/ntnxdb_client/insights/insights_interface/insights_interface.proto
	//  old_entity_proto: This is the old value of the proto that will be shared by IDF
	//					Entity Type of proto is mentioned here https://sourcegraph.ntnxdpro.com/ntnxdb-master/-/blob/ntnxdb_client/insights/insights_interface/insights_interface.proto
}

func (c *CyclopsMonitor) availabilityZoneAddCb(entity_proto []byte) {
	//  Callback when a 'availability_zone_physical' entity is added in IDF.
	// Args:
	//	entity_proto: This is the new updated proto that will be shared by IDF
	//					Entity Type of proto is mentioned here https://sourcegraph.ntnxdpro.com/ntnxdb-master/-/blob/ntnxdb_client/insights/insights_interface/insights_interface.proto
}

func (c *CyclopsMonitor) availabilityZoneDeleteCb(entity_proto []byte) {
	//  Callback when a 'availability_zone_physical' entity is removed in IDF.
	// Args:
	//	entity_proto: This is the new updated proto that will be shared by IDF
	//					Entity Type of proto is mentioned here https://sourcegraph.ntnxdpro.com/ntnxdb-master/-/blob/ntnxdb_client/insights/insights_interface/insights_interface.proto
}
