package api

import (
	"connection-microservice/application"
	"connection-microservice/model"
	"context"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
)

type ConnectionHandler struct {
	connectionService.UnimplementedConnectionServiceServer
	service *application.ConnectionService
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{service: service}
}

func (handler *ConnectionHandler) NewUserConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "NewUserConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := handler.service.CreateConnection(ctx, &model.Connection{UserId: in.Connection.UserId, ConnectedUserId: in.Connection.ConnectedUserId})
	if err != nil {
		return nil, err
	}
	return &connectionService.UserConnectionResponse{Ok: true}, nil
}

func (handler *ConnectionHandler) ApproveConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "ApproveConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := handler.service.ApproveConnection(ctx, in.Connection.UserId, in.Connection.ConnectedUserId)
	if err != nil {
		return nil, err
	}

	return &connectionService.UserConnectionResponse{Ok: true}, nil
}
