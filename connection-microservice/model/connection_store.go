package model

import "context"

type ConnectionStore interface {
	CreateConnection(ctx context.Context, connection *Connection) error
	UpdateConnection(ctx context.Context, connection *Connection) error
	GetConnectionByUsersId(ctx context.Context, userId string, connectedUserId string) (*Connection, error)
}
