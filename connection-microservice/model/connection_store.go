package model

import "context"

type ConnectionStore interface {
	CreateConnection(ctx context.Context, connection *Connection) error
}
