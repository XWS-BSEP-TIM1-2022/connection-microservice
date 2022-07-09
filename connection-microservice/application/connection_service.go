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
	"github.com/sirupsen/logrus"
)

type ConnectionService struct {
	store        model.ConnectionStore
	userClient   userService.UserServiceClient
	config       *config.Config
	blockService *BlockService
}

var Log = logrus.New()

func NewConnectionService(store model.ConnectionStore, c *config.Config, blockService *BlockService) *ConnectionService {
	return &ConnectionService{
		store:        store,
		blockService: blockService,
		config:       c,
		userClient:   services.NewUserClient(fmt.Sprintf("%s:%s", c.UserServiceHost, c.UserServicePort))}
}

func (service *ConnectionService) CreateConnection(ctx context.Context, connection *model.Connection) (*model.Connection, error) {
	Log.Info("Creating new connection by user with id: " + connection.UserId + " , with user with id: " + connection.ConnectedUserId)

	span := tracer.StartSpanFromContext(ctx, "CreateConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isBlocked, _ := service.blockService.IsBlockedAny(ctx, connection.UserId, connection.ConnectedUserId)

	if isBlocked {
		Log.Warn("Cant create connection, user with id: " + connection.UserId + " is blocked.")
		return nil, errors.New("user is blocked")
	}

	isPrivate, err := service.userClient.IsUserPrivateRequest(ctx, &userService.UserIdRequest{UserId: connection.ConnectedUserId})

	if err != nil {
		Log.Error("Error while creating connection. Error: " + err.Error())
		return nil, err
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

func (service *ConnectionService) ApproveConnection(ctx context.Context, userId string, connectedUserId string) (*model.Connection, error) {
	Log.Info("Approving connection by user with id: " + connectedUserId + " , for connection request from user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "ApproveConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isBlocked, _ := service.blockService.IsBlockedAny(ctx, userId, connectedUserId)

	if isBlocked {
		Log.Warn("Cant approve connection, user with id: " + connectedUserId + " is blocked.")
		return nil, errors.New("user is blocked")
	}

	connection, err := service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
	if err != nil {
		Log.Error("Error while approving connection. Error: " + err.Error())
		return nil, err
	}
	if connection.PendingConnection {
		connection.IsConnected = true
		connection.PendingConnection = false
	} else {
		return nil, errors.New("not pending connection")
	}
	conn, err := service.store.UpdateConnection(ctx, connection)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (service *ConnectionService) RejectConnection(ctx context.Context, userId string, connectedUserId string) error {
	Log.Info("Rejecting connection of users with id1: " + userId + " , id2: " + connectedUserId)

	span := tracer.StartSpanFromContext(ctx, "RejectConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isBlocked, _ := service.blockService.IsBlockedAny(ctx, userId, connectedUserId)

	if isBlocked {
		return errors.New("user is blocked")
	}

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
	Log.Info("Deleting connection of users with id1: " + userId + " , id2: " + connectedUserId)

	span := tracer.StartSpanFromContext(ctx, "DeleteConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isBlocked, _ := service.blockService.IsBlockedAny(ctx, userId, connectedUserId)

	if isBlocked {
		return errors.New("user is blocked")
	}

	return service.store.DeleteConnection(ctx, userId, connectedUserId)
}

// GetAllConnectionsByUserId isConnected = true || false
func (service *ConnectionService) GetAllConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	Log.Info("Get all connections of user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "GetAllConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllConnectionsByUserId(ctx, userId)
}

// GetFollowings isConnected = true
func (service *ConnectionService) GetFollowings(ctx context.Context, userId string) ([]*model.Connection, error) {
	Log.Info("Get followings of user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "GetConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetFollowings(ctx, userId)
}

// GetFollowers isConnected = true
func (service *ConnectionService) GetFollowers(ctx context.Context, connectedUserId string) ([]*model.Connection, error) {
	Log.Info("Get followers of user with id: " + connectedUserId)

	span := tracer.StartSpanFromContext(ctx, "GetConnectionsByConnectedUserid")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetFollowers(ctx, connectedUserId)
}

func (service *ConnectionService) GetAllRequestConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	Log.Info("Get all request connections of user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "GetAllRequestConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllRequestConnectionsByUserId(ctx, userId)
}

func (service *ConnectionService) GetAllPendingConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	Log.Info("Get all pending connections of user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "GetAllPendingConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllPendingConnectionsByUserId(ctx, userId)
}

func (service *ConnectionService) ApproveAllConnection(ctx context.Context, userId string) error {
	Log.Info("Approve all connections of user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "ApproveAllConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	pendingConnections, err := service.GetAllRequestConnectionsByUserId(ctx, userId)
	if err != nil {
		return err
	}

	for _, connection := range pendingConnections {
		_, err := service.ApproveConnection(ctx, connection.UserId, connection.ConnectedUserId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *ConnectionService) ChangeMessageNotification(ctx context.Context, userId string, connectedUserId string) (*model.Connection, error) {
	Log.Info("Change message notification")

	span := tracer.StartSpanFromContextMetadata(ctx, "ChangeMessageNotification")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isBlocked, _ := service.blockService.IsBlockedAny(ctx, userId, connectedUserId)

	if isBlocked {
		return nil, errors.New("user is blocked")
	}

	connection, err := service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
	if err != nil {
		return nil, err
	}

	connection.IsMessageNotificationEnabled = !connection.IsMessageNotificationEnabled

	conn, err := service.store.UpdateConnection(ctx, connection)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (service *ConnectionService) ChangePostNotification(ctx context.Context, userId string, connectedUserId string) (*model.Connection, error) {
	Log.Info("Change post notification")

	span := tracer.StartSpanFromContextMetadata(ctx, "ChangePostNotification")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isBlocked, _ := service.blockService.IsBlockedAny(ctx, userId, connectedUserId)

	if isBlocked {
		return nil, errors.New("user is blocked")
	}

	connection, err := service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
	if err != nil {
		return nil, err
	}

	connection.IsPostNotificationEnabled = !connection.IsPostNotificationEnabled

	conn, err := service.store.UpdateConnection(ctx, connection)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (service *ConnectionService) ChangeCommentNotification(ctx context.Context, userId string, connectedUserId string) (*model.Connection, error) {
	Log.Info("Change comment notification")

	span := tracer.StartSpanFromContextMetadata(ctx, "ChangeCommentNotification")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	isBlocked, _ := service.blockService.IsBlockedAny(ctx, userId, connectedUserId)

	if isBlocked {
		return nil, errors.New("user is blocked")
	}

	connection, err := service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
	if err != nil {
		return nil, err
	}

	connection.IsCommentNotificationEnabled = !connection.IsCommentNotificationEnabled

	conn, err := service.store.UpdateConnection(ctx, connection)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (service *ConnectionService) GetConnection(ctx context.Context, userId string, connectedUserId string) (*model.Connection, error) {
	Log.Info("Get connection of users with id1: " + userId + ", id2: " + connectedUserId)

	span := tracer.StartSpanFromContextMetadata(ctx, "GetConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetConnectionByUsersId(ctx, userId, connectedUserId)
}

func (service *ConnectionService) GetAllSuggestionsByUserId(ctx context.Context, userId string) ([]string, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllSuggestionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	potentialUsers := StringSet{set: map[string]bool{}}

	followings, err := service.store.GetFollowings(ctx, userId)
	if err != nil {
		return nil, err
	}

	// friends of my friends
	for _, connection := range followings {
		users, err := service.store.GetFollowingsOfMyFollowings(ctx, connection.ConnectedUserId, userId)
		if err == nil {
			for _, user := range users {
				potentialUsers.Add(user)
			}
		}
	}

	users, err := service.store.GetRandom(ctx, userId, 15-potentialUsers.length())
	if err == nil {
		for _, user := range users {
			potentialUsers.Add(user)
		}
	}

	var retVal []string

	for i, _ := range potentialUsers.set {
		blockedAny, _ := service.blockService.IsBlockedAny(ctx, userId, i)
		if i != userId && !blockedAny {
			retVal = append(retVal, i)
		}
	}

	return retVal, nil
}

type StringSet struct {
	set map[string]bool
}

func (set *StringSet) Add(i string) bool {
	_, found := set.set[i]
	set.set[i] = true
	return !found //False if it existed already
}

func (set *StringSet) length() int {
	return len(set.set)
}
