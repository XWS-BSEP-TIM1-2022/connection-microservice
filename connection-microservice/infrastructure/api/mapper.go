package api

import (
	"connection-microservice/model"
	connectionService "github.com/XWS-BSEP-TIM1-2022/dislinkt/util/proto/connection"
)

func mapConnection(connection *model.Connection) *connectionService.Connection {
	connectionPb := &connectionService.Connection{
		UserId:                       connection.UserId,
		ConnectedUserId:              connection.ConnectedUserId,
		IsConnected:                  connection.IsConnected,
		PendingConnection:            connection.PendingConnection,
		IsMessageNotificationEnabled: connection.IsMessageNotificationEnabled,
		IsPostNotificationEnabled:    connection.IsPostNotificationEnabled,
		IsCommentNotificationEnabled: connection.IsCommentNotificationEnabled,
	}
	return connectionPb
}
