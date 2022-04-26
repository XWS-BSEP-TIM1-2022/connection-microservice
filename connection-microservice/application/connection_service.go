package application

import (
	"connection-microservice/model"
	"context"
)

type ConnectionService struct {
	store model.ConnectionStore
}

func NewConnectionService(store model.ConnectionStore) *ConnectionService {
	return &ConnectionService{store: store}
}

func (service *ConnectionService) CreateConnection(ctx context.Context, connection *model.Connection) error {
	return service.store.CreateConnection(ctx, connection)
}
