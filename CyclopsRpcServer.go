package main

import (
	"context"
	"fmt"

	"github.com/80asis/cyclops/cyclops"
)

type LocalCyclopsRpcSvcServer struct {
	cyclops.UnimplementedCyclopsRpcSvcServer
}

// Here we can add the logic for TriggerEntitySync. Create the TriggerEntitySync and TriggerEntitySyncAZ Task for AZ Pairing.
func (localServer *LocalCyclopsRpcSvcServer) TriggerEntitySync(ctx context.Context, request *cyclops.TriggerEntitySyncArg) (*cyclops.TriggerEntitySyncRet, error) {
	fmt.Println("Received TriggerEntitySyncArg: ", request)
	return &cyclops.TriggerEntitySyncRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}

// Here we can add the logic for TriggerEntitySyncAZ. Create the TriggerEntitySyncAZ Task for Independent and Force Sync Task.
func (localServer LocalCyclopsRpcSvcServer) TriggerEntitySyncAZ(ctx context.Context, request *cyclops.TriggerEntitySyncAZArg) (*cyclops.TriggerEntitySyncAZRet, error) {
	fmt.Println("Received TriggerEntitySyncAZArg: ", request)
	return &cyclops.TriggerEntitySyncAZRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}

// Here we can add the logic for ExecuteEntitySync. Create ExecuteEntitySync Task for updating the entity locally.
func (localServer LocalCyclopsRpcSvcServer) ExecuteEntitySync(ctx context.Context, request *cyclops.ExecuteEntitySyncArg) (*cyclops.ExecuteEntitySyncRet, error) {
	fmt.Println("Received ExecuteEntitySyncArg: ", request)
	return &cyclops.ExecuteEntitySyncRet{
		TaskUuid: []byte{
			72,
			101,
		},
	}, nil
}
