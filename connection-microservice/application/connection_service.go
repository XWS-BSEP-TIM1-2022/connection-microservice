package application

import (
	"connection-microservice/model"
	"connection-microservice/startup/config"
	"context"
	"errors"
	"fmt"
	userService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/user"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/services"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
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
	span := tracer.StartSpanFromContext(ctx, "CreateConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isPrivate, err := service.userClient.IsUserPrivateRequest(ctx, &userService.UserIdRequest{UserId: connection.ConnectedUserId})

	if err != nil {
		return err
	}

	if isPrivate.IsPrivate {
		connection.IsConnected = false
		connection.PendingConnection = true
	} else {
		connection.IsConnected = true
		connection.PendingConnection = false
	}

	return service.store.CreateConnection(ctx, connection)
}

func (service *ConnectionService) ApproveConnection(ctx context.Context, userId string, connectedUserId string) error {
	span := tracer.StartSpanFromContext(ctx, "ApproveConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
	if err != nil {
		return err
	}
	if connection.PendingConnection {
		connection.IsConnected = true
		connection.PendingConnection = false
	} else {
		return errors.New("not pending connection")
	}
	err = service.store.UpdateConnection(ctx, connection)
	if err != nil {
		return err
	}
	return nil
}

func (service *ConnectionService) RejectConnection(ctx context.Context, userId string, connectedUserId string) error {
	span := tracer.StartSpanFromContext(ctx, "RejectConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
	if err != nil {
		return err
	}
	if !connection.PendingConnection {
		return errors.New("not pending connection")
	}
	return service.store.DeleteConnection(ctx, userId, connectedUserId)
}

func (service *ConnectionService) DeleteConnection(ctx context.Context, userId string, connectedUserId string) error {
	span := tracer.StartSpanFromContext(ctx, "DeleteConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.DeleteConnection(ctx, userId, connectedUserId)
}

// GetAllConnectionsByUserId isConnected = true || false
func (service *ConnectionService) GetAllConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllConnectionsByUserId(ctx, userId)
}

// GetFollowings isConnected = true
func (service *ConnectionService) GetFollowings(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetFollowings(ctx, userId)
}

// GetFollowers isConnected = true
func (service *ConnectionService) GetFollowers(ctx context.Context, connectedUserId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnectionsByConnectedUserid")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetFollowers(ctx, connectedUserId)
}

func (service *ConnectionService) GetAllRequestConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllRequestConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllRequestConnectionsByUserId(ctx, userId)
}

func (service *ConnectionService) GetAllPendingConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllPendingConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllPendingConnectionsByUserId(ctx, userId)
}
