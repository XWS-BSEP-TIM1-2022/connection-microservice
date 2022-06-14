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
	Log.Info("User with id: " + userId + " blocks user with id: " + blockedUserId)

	span := tracer.StartSpanFromContext(ctx, "BlockUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := service.store.BlockUser(ctx, model.Block{UserId: userId, BlockedUserId: blockedUserId})
	if err != nil {
		Log.Error("Error on blocking user. Error: " + err.Error())
		return err
	}

	err = service.connectionStore.DeleteConnection(ctx, userId, blockedUserId)
	if err != nil {
		Log.Error("Error deleting connection after blocking user. Error: " + err.Error())
		return err
	}

	err = service.connectionStore.DeleteConnection(ctx, blockedUserId, userId)
	if err != nil {
		Log.Error("Error deleting connection after blocking user. Error: " + err.Error())
		return err
	}

	return nil
}

func (service *BlockService) UnblockUser(ctx context.Context, userId string, blockedUserId string) error {
	Log.Info("User with id: " + userId + " unblocks user with id: " + blockedUserId)

	span := tracer.StartSpanFromContext(ctx, "BlockUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	err := service.store.UnblockUser(ctx, model.Block{UserId: userId, BlockedUserId: blockedUserId})
	if err != nil {
		Log.Error("Error on unblocking user. Error: " + err.Error())
		return err
	}
	return nil
}

func (service *BlockService) IsBlocked(ctx context.Context, userId string, blockedUserId string) (bool, error) {
	Log.Info("Is user with id: " + blockedUserId + " blocked by user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "IsBlocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.IsBlocked(ctx, model.Block{
		UserId:        userId,
		BlockedUserId: blockedUserId,
	})
}

func (service *BlockService) IsBlockedAny(ctx context.Context, userId string, blockedUserId string) (bool, error) {
	Log.Info("Are any of users with id1: " + userId + " , id2: " + blockedUserId + " blocked")

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

func (service *BlockService) GetBlocked(ctx context.Context, userId string) ([]string, error) {
	Log.Info("Get blocked of user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "GetBlocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetBlocked(ctx, userId)
}

func (service *BlockService) GetBlockedBy(ctx context.Context, userId string) ([]string, error) {
	Log.Info("Get users blocked by user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "GetBlockedBy")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetBlockedBy(ctx, userId)
}

func (service *BlockService) GetBlockedAny(ctx context.Context, userId string) ([]string, error) {
	Log.Info("Get blocked any of user with id: " + userId)

	span := tracer.StartSpanFromContext(ctx, "GetBlockedAny")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	blocked, err := service.store.GetBlocked(ctx, userId)
	if err != nil {
		return nil, err
	}
	blockedBy, err := service.store.GetBlockedBy(ctx, userId)
	if err != nil {
		return nil, err
	}

	blocked = append(blocked, blockedBy...)
	return blocked, nil
}
