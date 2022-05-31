package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UserId primitive.ObjectID
}

type Connection struct {
	UserId                       string
	ConnectedUserId              string
	IsConnected                  bool
	PendingConnection            bool
	IsMessageNotificationEnabled bool
	IsPostNotificationEnabled    bool
	IsCommentNotificationEnabled bool
}
