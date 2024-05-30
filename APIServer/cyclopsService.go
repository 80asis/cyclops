package APIServer

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type EnableEntitySyncArg struct {
	EntityUuid   string `json:"entity_uuid,omitempty"`
	EntityType   string `json:"entity_type,omitempty"`
	RemoteAzUuid string `json:"remote_az_uuid,omitempty"`
	ForceSync    bool   `json:"force_sync,omitempty"`
}

func (handler *Handler) CreateEntityConfig(w http.ResponseWriter, r *http.Request) {
	log.Info("Received request for create entity sync")
	var entityConfigArgs EnableEntitySyncArg
	if err := json.NewDecoder(r.Body).Decode(&entityConfigArgs); err != nil {
		log.Error("Unable to parse the body of the request", err)
		return
	}
	response, err := handler.Service.CreateEntityConfig(r.Context(), entityConfigArgs)
	if err != nil {
		log.Error(err)
		return
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}
func (handler *Handler) DeleteEntityConfig(w http.ResponseWriter, r *http.Request) {
	handler.Service.DeleteEntityConfig(r.Context())
}
func (handler *Handler) GetEntityPolicy(w http.ResponseWriter, r *http.Request) {
	handler.Service.GetEntityPolicy(r.Context())
}
func (handler *Handler) CreateForceSync(w http.ResponseWriter, r *http.Request) {
	handler.Service.CreateForceSync(r.Context())
}
func (handler *Handler) GetEntityPolicies(w http.ResponseWriter, r *http.Request) {
	handler.Service.GetEntityPolicies(r.Context())
}
