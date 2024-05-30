package APIServer

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Router  *mux.Router
	Service CyclopsService
	Server  *http.Server
}

type Response struct {
	Message string `json:"message"`
}

func NewHandler(service CyclopsService) *Handler {
	log.Info("Setting up the handler")
	handler := &Handler{
		Service: service,
	}
	handler.Router = mux.NewRouter()
	handler.Router.Use(JsonMiddleware)
	handler.Router.Use(LoggingMiddleware)
	handler.Router.Use(TimeoutMiddlerware)

	handler.mapRoutes()

	handler.Server = &http.Server{
		Addr: "0.0.0.0:8080",

		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler.Router,
	}
	return handler
}

func (handler *Handler) mapRoutes() {
	handler.Router.HandleFunc("/alive", handler.AliveCheck).Methods("GET")
	// CREATE Entity SYNC config
	handler.Router.HandleFunc("/api/data-policies/v4.0.b1/config/entity-sync-policies/", handler.CreateEntityConfig).Methods("POST")
	// DELETE Entity SYNC policy
	handler.Router.HandleFunc("/api/data-policies/v4.0.b1/mgmt/entity-sync-policies/{id}", handler.DeleteEntityConfig).Methods("DELETE")
	// GET Entity SYNC policy
	handler.Router.HandleFunc("/api/data-policies/v4.0.b1/mgmt/entity-sync-policies/{id}", handler.GetEntityPolicy).Methods("GET")
	// Force_entity_sync_policy
	handler.Router.HandleFunc("/api/data-policies/v4.0.b1/mgmt/entity-sync-policy/{id}/$actions/force-sync", handler.CreateForceSync).Methods("POST")
	// LIST Entity SYNC policy
	handler.Router.HandleFunc("/api/prism/v4.0.b1/mgmt/entity-sync-policies", handler.GetEntityPolicies).Methods("GET")
}

func (handler *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am alive"}); err != nil {
		panic(err)
	}
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for
	ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()
	h.Server.Shutdown(ctx)

	log.Println("Shutting down server gracefully")
	return nil
}
