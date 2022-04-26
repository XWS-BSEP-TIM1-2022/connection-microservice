package persistance

import (
	"connection-microservice/model"
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM1-2022/dislinkt/util/tracer"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type ConnectionNeo4jStore struct {
	driver neo4j.Driver
}

func NewConnectionNeo4jStore(driver neo4j.Driver) model.ConnectionStore {
	return &ConnectionNeo4jStore{
		driver: driver,
	}
}

func (store *ConnectionNeo4jStore) CreateConnection(ctx context.Context, connection *model.Connection) error {
	span := tracer.StartSpanFromContext(ctx, "CreateConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MERGE (user:User {userId:$userId}) "+
			"MERGE (connectedUser:User {userId:$connectedUserId}) "+
			"MERGE (user)-[:CONNECT {isConnected:$isConnected, isApprovedConnection:$isApprovedConnection, pendingConnection:$pendingConnection}]->(connectedUser) RETURN user",
			map[string]interface{}{
				"userId":               connection.UserId,
				"connectedUserId":      connection.ConnectedUserId,
				"isConnected":          connection.IsConnected,
				"isApprovedConnection": connection.IsApprovedConnection,
				"pendingConnection":    connection.PendingConnection,
			})
		if err != nil {
			return nil, err
		}

		if res.Next() {
			return res.Record().Values[0], nil
		}
		return nil, res.Err()

	})

	fmt.Println(result)
	if err != nil {
		return err
	}

	return nil
}
