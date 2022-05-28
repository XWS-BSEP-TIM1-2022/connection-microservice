package application

import (
	"connection-microservice/model"
	"connection-microservice/startup/config"
	"context"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
)

type BlockService struct {
	store           model.BlockStore
	connectionStore model.ConnectionStore
	config          *config.Config
}

func NewBlockService(store model.BlockStore, connectionStore model.ConnectionStore, c *config.Config) *BlockService {
	return &BlockService{
		store:           store,
		connectionStore: connectionStore,
		config:          c,
	}
}

func (service *BlockService) BlockUser(ctx context.Context, userId string, blockedUserId string) error {
	span := tracer.StartSpanFromContext(ctx, "BlockUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := service.store.BlockUser(ctx, model.Block{UserId: userId, BlockedUserId: blockedUserId})
	if err != nil {
		return err
	}

	err = service.connectionStore.DeleteConnection(ctx, userId, blockedUserId)
	if err != nil {
		return err
	}

	err = service.connectionStore.DeleteConnection(ctx, blockedUserId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (service *BlockService) UnblockUser(ctx context.Context, userId string, blockedUserId string) error {
	span := tracer.StartSpanFromContext(ctx, "BlockUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := service.store.UnblockUser(ctx, model.Block{UserId: userId, BlockedUserId: blockedUserId})
	if err != nil {
		return err
	}
	return nil
}

func (service *BlockService) IsBlocked(ctx context.Context, userId string, blockedUserId string) (bool, error) {
	span := tracer.StartSpanFromContext(ctx, "IsBlocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.IsBlocked(ctx, model.Block{
		UserId:        userId,
		BlockedUserId: blockedUserId,
	})
}

func (service *BlockService) IsBlockedAny(ctx context.Context, userId string, blockedUserId string) (bool, error) {
	span := tracer.StartSpanFromContext(ctx, "IsBlockedAny")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	blocked, err := service.store.IsBlocked(ctx, model.Block{
		UserId:        userId,
		BlockedUserId: blockedUserId,
	})

	if err != nil {
		return false, err
	}
	if !blocked {
		blocked, err = service.store.IsBlocked(ctx, model.Block{
			UserId:        blockedUserId,
			BlockedUserId: userId,
		})

		if err != nil {
			return false, err
		}
	}

	return blocked, nil
}
