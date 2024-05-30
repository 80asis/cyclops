package APIServer

import (
	"context"

	"github.com/80asis/cyclops/cyclops"
	log "github.com/sirupsen/logrus"
)

func NewService() CyclopsService {
	log.Info("Initializing the service")
	return &Service{}
}

type CyclopsService interface {
	CreateEntityConfig(ctx context.Context, request EnableEntitySyncArg) (cyclops.EnableEntitySyncRet, error)
	DeleteEntityConfig(ctx context.Context)
	GetEntityPolicy(ctx context.Context)
	CreateForceSync(ctx context.Context)
	GetEntityPolicies(ctx context.Context)
}

type Service struct{}

func (s *Service) CreateEntityConfig(ctx context.Context, request EnableEntitySyncArg) (cyclops.EnableEntitySyncRet, error) {
	ret := cyclops.EnableEntitySyncRet{
		TaskUuid: []byte{
			40,
			72, 40,
			72, 40,
			72, 40,
			72, 40,
			72, 40,
			72, 40,
			72,
		},
	}
	log.Infof("Policy Enablement Task Created with Task UUID: %v", ret.TaskUuid)
	return ret, nil
}
func (s *Service) DeleteEntityConfig(ctx context.Context) {}
func (s *Service) GetEntityPolicy(ctx context.Context)    {}
func (s *Service) CreateForceSync(ctx context.Context)    {}
func (s *Service) GetEntityPolicies(ctx context.Context)  {}
