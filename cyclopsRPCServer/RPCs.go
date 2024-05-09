package cyclopsRPCServer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/80asis/cyclops/cyclops"
	"github.com/80asis/cyclops/cyclopsMonitor"
)

type LocalCyclopsRpcSvcServer struct {
	cyclops.UnimplementedCyclopsRpcSvcServer
	Monitor *cyclopsMonitor.CyclopsMonitor
}

// Here we can add the logic for TriggerEntitySync
func (localServer *LocalCyclopsRpcSvcServer) TriggerEntitySync(ctx context.Context, request *cyclops.TriggerEntitySyncArg) (*cyclops.TriggerEntitySyncRet, error) {
	fmt.Println("Received TriggerEntitySync request: ", request)
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error serializing data")
	}
	localServer.Monitor.AddNotification(jsonData)
	return &cyclops.TriggerEntitySyncRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}

// Here we can add the logic for TriggerEntitySyncAZ
func (localServer LocalCyclopsRpcSvcServer) TriggerEntitySyncAZ(ctx context.Context, request *cyclops.TriggerEntitySyncAZArg) (*cyclops.TriggerEntitySyncAZRet, error) {
	fmt.Println("Received TriggerEntitySyncAZ request: ", request)
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error serializing data")
	}
	localServer.Monitor.AddNotification(jsonData)
	return &cyclops.TriggerEntitySyncAZRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}

// Here we can add the logic for ExecuteEntitySync
func (localServer LocalCyclopsRpcSvcServer) ExecuteEntitySync(ctx context.Context, request *cyclops.ExecuteEntitySyncArg) (*cyclops.ExecuteEntitySyncRet, error) {
	fmt.Println("Received ExecuteEntitySync request: ", request)
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error serializing data")
	}
	localServer.Monitor.AddNotification(jsonData)
	return &cyclops.ExecuteEntitySyncRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}

// Here we can add the logic for EnableEntitySync
func (localServer LocalCyclopsRpcSvcServer) EnableEntitySync(ctx context.Context, request *cyclops.EnableEntitySyncArg) (*cyclops.EnableEntitySyncRet, error) {
	fmt.Println("Received EnableEntitySync request: ", request)
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error serializing data")
	}
	localServer.Monitor.AddNotification(jsonData)
	return &cyclops.EnableEntitySyncRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}

// Here we can add the logic for DisableEntitySync
func (localServer LocalCyclopsRpcSvcServer) DisableEntitySync(ctx context.Context, request *cyclops.DisableEntitySyncArg) (*cyclops.DisableEntitySyncRet, error) {
	fmt.Println("Received EnableEntitySync request: ", request)
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error serializing data")
	}
	localServer.Monitor.AddNotification(jsonData)
	return &cyclops.DisableEntitySyncRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}
