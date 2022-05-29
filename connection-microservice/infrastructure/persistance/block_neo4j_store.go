package persistance

import (
	"connection-microservice/model"
	"context"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type BlockNeo4jStore struct {
	driver neo4j.Driver
}

func NewBlockNeo4jStore(driver neo4j.Driver) model.BlockStore {
	return &BlockNeo4jStore{
		driver: driver,
	}
}

func (store *BlockNeo4jStore) BlockUser(ctx context.Context, block model.Block) error {
	span := tracer.StartSpanFromContext(ctx, "BlockUser")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MERGE (user:User {userId:$userId}) "+
			"MERGE (blockedUser:User {userId:$blockedUserId}) "+
			"MERGE (user)-[b:BLOCK]->(blockedUser) RETURN b",
			map[string]interface{}{
				"userId":        block.UserId,
				"blockedUserId": block.BlockedUserId,
			})
		if err != nil {
			return nil, err
		}

		if res.Next() {
			return res.Record().Values[0], nil
		}
		return nil, res.Err()

	})

	if err != nil {
		return err
	}

	return nil
}

func (store *BlockNeo4jStore) UnblockUser(ctx context.Context, block model.Block) error {
	span := tracer.StartSpanFromContext(ctx, "DeleteConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run("MATCH (user {userId:$userId})-[b:BLOCK]->(blockedUser {userId:$blockedUserId}) "+
			"DELETE b",
			map[string]interface{}{
				"userId":        block.UserId,
				"blockedUserId": block.BlockedUserId,
			})
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (store *BlockNeo4jStore) IsBlocked(ctx context.Context, block model.Block) (bool, error) {
	span := tracer.StartSpanFromContext(ctx, "IsBlocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	blocked, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MATCH (user {userId:$userId})-[b:BLOCK]->(blockedUser {userId:$blockedUserId}) "+
			"RETURN b",
			map[string]interface{}{
				"userId":        block.UserId,
				"blockedUserId": block.BlockedUserId,
			})
		if err != nil {
			return false, err
		}

		if res.Next() {
			if res.Record() != nil {
				return true, nil
			}
		}
		return false, res.Err()

	})

	if blocked == true {
		return true, nil
	}

	if err != nil {
		return false, err
	}
	return false, nil
}

func (store *BlockNeo4jStore) GetBlocked(ctx context.Context, userId string) ([]string, error) {
	span := tracer.StartSpanFromContext(ctx, "GetBlocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var blockedUserIds []string
	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MATCH (user {userId:$userId})-[b:BLOCK]->(blockedUser) "+
			"RETURN blockedUser.userId",
			map[string]interface{}{
				"userId": userId,
			})
		if err != nil {
			return false, err
		}

		for res.Next() {
			blockedUserIds = append(blockedUserIds, res.Record().Values[0].(string))
		}
		return false, res.Err()

	})

	if err != nil {
		return nil, err
	}
	return blockedUserIds, nil

}

func (store *BlockNeo4jStore) GetBlockedBy(ctx context.Context, userId string) ([]string, error) {
	span := tracer.StartSpanFromContext(ctx, "GetBlocked")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var blockedUserIds []string
	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MATCH (user)-[b:BLOCK]->(blockedUser {userId:$userId}) "+
			"RETURN user.userId",
			map[string]interface{}{
				"userId": userId,
			})
		if err != nil {
			return false, err
		}

		for res.Next() {
			blockedUserIds = append(blockedUserIds, res.Record().Values[0].(string))
		}
		return false, res.Err()

	})

	if err != nil {
		return nil, err
	}
	return blockedUserIds, nil

}
