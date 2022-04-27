package application

import (
	"connection-microservice/model"
	"connection-microservice/startup/config"
	"context"
	"fmt"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
)

type ConnectionService struct {
	store      model.ConnectionStore
	userClient userService.UserServiceClient
	config     *config.Config
}

func NewConnectionService(store model.ConnectionStore, c *config.Config) *ConnectionService {
	return &ConnectionService{
		store:      store,
		userClient: services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort))}
}

func (service *ConnectionService) CreateConnection(ctx context.Context, connection *model.Connection) error {
	isPrivate, err := service.userClient.IsUserPrivateRequest(ctx, &userService.UserIdRequest{UserId: connection.ConnectedUserId})

	if err != nil {
		return err
	}

	if isPrivate.IsPrivate {
		connection.IsConnected = false
		connection.PendingConnection = true
		connection.IsApprovedConnection = false
	} else {
		connection.IsConnected = true
		connection.PendingConnection = false
		connection.IsApprovedConnection = false
	}

	return service.store.CreateConnection(ctx, connection)
}

func (service *ConnectionService) ApproveConnection(ctx context.Context, userId string, connectedUserId string) error {
	connection, err := service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
	if err != nil {
		return err
	}
	connection.IsApprovedConnection = true
	connection.IsConnected = true
	err = service.store.UpdateConnection(ctx, connection)
	if err != nil {
		return err
	}
	return nil
}
