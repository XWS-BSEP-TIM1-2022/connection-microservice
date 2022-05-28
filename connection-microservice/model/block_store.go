package model

import "context"

type BlockStore interface {
	BlockUser(ctx context.Context, block Block) error
	UnblockUser(ctx context.Context, block Block) error
	IsBlocked(ctx context.Context, block Block) (bool, error)
}
