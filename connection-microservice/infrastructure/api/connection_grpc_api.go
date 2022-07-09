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
	service      *application.ConnectionService
	blockService *application.BlockService
}

func NewConnectionHandler(service *application.ConnectionService, blockService *application.BlockService) *ConnectionHandler {
	return &ConnectionHandler{service: service,
		blockService: blockService}
}

func (handler *ConnectionHandler) NewUserConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "NewUserConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := handler.service.CreateConnection(ctx, &model.Connection{UserId: in.Connection.UserId, ConnectedUserId: in.Connection.ConnectedUserId})
	if err != nil {
		return nil, err
	}
	return &connectionService.UserConnectionResponse{Connection: mapConnection(connection)}, nil
}

func (handler *ConnectionHandler) ApproveConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "ApproveConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := handler.service.ApproveConnection(ctx, in.Connection.UserId, in.Connection.ConnectedUserId)
	if err != nil {
		return nil, err
	}

	return &connectionService.UserConnectionResponse{Connection: mapConnection(connection)}, nil
}

func (handler *ConnectionHandler) ApproveAllConnection(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.EmptyRequest, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "ApproveAllConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := handler.service.ApproveAllConnection(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &connectionService.EmptyRequest{}, nil
}

func (handler *ConnectionHandler) RejectConnection(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "RejectConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := handler.service.RejectConnection(ctx, in.Connection.UserId, in.Connection.ConnectedUserId)
	if err != nil {
		return nil, err
	}

	return &connectionService.UserConnectionResponse{Connection: nil}, nil
}

func (handler *ConnectionHandler) DeleteConnection(ctx context.Context, in *connectionService.Connection) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "DeleteConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := handler.service.DeleteConnection(ctx, in.UserId, in.ConnectedUserId)
	if err != nil {
		return nil, err
	}

	return &connectionService.UserConnectionResponse{Connection: nil}, nil
}

func (handler *ConnectionHandler) GetConnection(ctx context.Context, in *connectionService.Connection) (*connectionService.Connection, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllConnections")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := handler.service.GetConnection(ctx, in.UserId, in.ConnectedUserId)
	if err != nil {
		return nil, err
	}
	return mapConnection(connection), nil
}

func (handler *ConnectionHandler) GetAllConnections(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllConnections")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connections, err := handler.service.GetAllConnectionsByUserId(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	response := &connectionService.AllConnectionResponse{
		Connections: []*connectionService.Connection{},
	}
	for _, conn := range connections {
		current := mapConnection(conn)
		response.Connections = append(response.Connections, current)
	}
	return response, nil
}

func (handler *ConnectionHandler) GetFollowings(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetFollowings")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connections, err := handler.service.GetFollowings(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	response := &connectionService.AllConnectionResponse{
		Connections: []*connectionService.Connection{},
	}
	for _, conn := range connections {
		current := mapConnection(conn)
		response.Connections = append(response.Connections, current)
	}
	return response, nil
}

func (handler *ConnectionHandler) GetFollowers(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetFollowers")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connections, err := handler.service.GetFollowers(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	response := &connectionService.AllConnectionResponse{
		Connections: []*connectionService.Connection{},
	}
	for _, conn := range connections {
		current := mapConnection(conn)
		response.Connections = append(response.Connections, current)
	}
	return response, nil
}

func (handler *ConnectionHandler) GetAllRequestConnectionsByUserId(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllRequestConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connections, err := handler.service.GetAllRequestConnectionsByUserId(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	response := &connectionService.AllConnectionResponse{
		Connections: []*connectionService.Connection{},
	}
	for _, conn := range connections {
		current := mapConnection(conn)
		response.Connections = append(response.Connections, current)
	}
	return response, nil
}

func (handler *ConnectionHandler) GetAllPendingConnectionsByUserId(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.AllConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllPendingConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connections, err := handler.service.GetAllPendingConnectionsByUserId(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	response := &connectionService.AllConnectionResponse{
		Connections: []*connectionService.Connection{},
	}
	for _, conn := range connections {
		current := mapConnection(conn)
		response.Connections = append(response.Connections, current)
	}
	return response, nil
}

func (handler *ConnectionHandler) BlockUser(ctx context.Context, in *connectionService.BlockUserRequest) (*connectionService.EmptyRequest, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "BlockUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := handler.blockService.BlockUser(ctx, in.Block.UserId, in.Block.BlockUserId)
	if err != nil {
		return nil, err
	}
	return &connectionService.EmptyRequest{}, nil
}

func (handler *ConnectionHandler) UnblockUser(ctx context.Context, in *connectionService.BlockUserRequest) (*connectionService.EmptyRequest, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "UnblockUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := handler.blockService.UnblockUser(ctx, in.Block.UserId, in.Block.BlockUserId)
	if err != nil {
		return nil, err
	}
	return &connectionService.EmptyRequest{}, nil
}

func (handler *ConnectionHandler) IsBlocked(ctx context.Context, in *connectionService.Block) (*connectionService.IsBlockedResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "IsBlocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	blocked, err := handler.blockService.IsBlocked(ctx, in.UserId, in.BlockUserId)
	if err != nil {
		return nil, err
	}
	return &connectionService.IsBlockedResponse{Blocked: blocked}, nil
}

func (handler *ConnectionHandler) IsBlockedAny(ctx context.Context, in *connectionService.Block) (*connectionService.IsBlockedResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "IsBlockedAny")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	blocked, err := handler.blockService.IsBlockedAny(ctx, in.UserId, in.BlockUserId)
	if err != nil {
		return &connectionService.IsBlockedResponse{Blocked: false}, err
	}
	return &connectionService.IsBlockedResponse{Blocked: blocked}, nil
}

func (handler *ConnectionHandler) Blocked(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.BlockedResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "Blocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	blocked, err := handler.blockService.GetBlocked(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &connectionService.BlockedResponse{UsersId: blocked}, nil
}

func (handler *ConnectionHandler) BlockedBy(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.BlockedResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "Blocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	blocked, err := handler.blockService.GetBlockedBy(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &connectionService.BlockedResponse{UsersId: blocked}, nil
}

func (handler *ConnectionHandler) BlockedAny(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.BlockedResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "Blocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	blocked, err := handler.blockService.GetBlockedAny(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	return &connectionService.BlockedResponse{UsersId: blocked}, nil
}

func (handler *ConnectionHandler) ChangeMessageNotification(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "ChangeMessageNotification")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := handler.service.ChangeMessageNotification(ctx, in.Connection.UserId, in.Connection.ConnectedUserId)
	if err != nil {
		return nil, err
	}

	return &connectionService.UserConnectionResponse{Connection: mapConnection(connection)}, nil

}

func (handler *ConnectionHandler) ChangePostNotification(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "ChangePostNotification")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := handler.service.ChangePostNotification(ctx, in.Connection.UserId, in.Connection.ConnectedUserId)
	if err != nil {
		return nil, err
	}

	return &connectionService.UserConnectionResponse{Connection: mapConnection(connection)}, nil

}

func (handler *ConnectionHandler) ChangeCommentNotification(ctx context.Context, in *connectionService.UserConnectionRequest) (*connectionService.UserConnectionResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "ChangeCommentNotification")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	connection, err := handler.service.ChangeCommentNotification(ctx, in.Connection.UserId, in.Connection.ConnectedUserId)
	if err != nil {
		return nil, err
	}

	return &connectionService.UserConnectionResponse{Connection: mapConnection(connection)}, nil

}

func (handler *ConnectionHandler) GetAllSuggestionsByUserId(ctx context.Context, in *connectionService.UserIdRequest) (*connectionService.SuggestionsResponse, error) {
	span := tracer.StartSpanFromContextMetadata(ctx, "GetAllSuggestionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	suggestions, err := handler.service.GetAllSuggestionsByUserId(ctx, in.UserId)

	if err != nil {
		return nil, err
	}

	response := &connectionService.SuggestionsResponse{
		SuggestionUserIds: []*connectionService.UserIdRequest{},
	}
	for _, userId := range suggestions {
		response.SuggestionUserIds = append(response.SuggestionUserIds, &connectionService.UserIdRequest{UserId: userId})
	}
	return response, nil
}
