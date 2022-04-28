package model

import "context"

type ConnectionStore interface {
	CreateConnection(ctx context.Context, connection *Connection) (*Connection, error)
	UpdateConnection(ctx context.Context, connection *Connection) (*Connection, error)
	DeleteConnection(ctx context.Context, userId string, connectedUserId string) error
	GetAllConnectionsByUserId(ctx context.Context, userId string) ([]*Connection, error)
	GetConnectionByUsersId(ctx context.Context, userId string, connectedUserId string) (*Connection, error)
	GetFollowings(ctx context.Context, userId string) ([]*Connection, error)
	GetFollowers(ctx context.Context, connectedUserId string) ([]*Connection, error)
	GetAllRequestConnectionsByUserId(ctx context.Context, userId string) ([]*Connection, error)
	GetAllPendingConnectionsByUserId(ctx context.Context, userId string) ([]*Connection, error)
}
