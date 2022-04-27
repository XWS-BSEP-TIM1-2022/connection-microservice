package persistance

import (
	"connection-microservice/model"
	"context"
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

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MERGE (user:User {userId:$userId}) "+
			"MERGE (connectedUser:User {userId:$connectedUserId}) "+
			"MERGE (user)-[c:CONNECT {isConnected:$isConnected, isApprovedConnection:$isApprovedConnection, pendingConnection:$pendingConnection}]->(connectedUser) RETURN c",
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

	if err != nil {
		return err
	}

	return nil
}

func (store *ConnectionNeo4jStore) UpdateConnection(ctx context.Context, connection *model.Connection) error {
	span := tracer.StartSpanFromContext(ctx, "CreateConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MATCH (user {userId:$userId})-[c]->(connectedUser {userId:$connectedUserId}) "+
			"SET c.isConnected=$isConnected, c.isApprovedConnection=$isApprovedConnection, c.pendingConnection=$pendingConnection "+
			"RETURN c",
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

	if err != nil {
		return err
	}

	return nil
}

func (store *ConnectionNeo4jStore) GetConnectionByUsersId(ctx context.Context, userId string, connectedUserId string) (*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var connection = model.Connection{}
	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MATCH (user {userId:$userId})-[c]->(connectedUser {userId:$connectedUserId}) "+
			"RETURN c.isConnected, c.isApprovedConnection, c.pendingConnection",
			map[string]interface{}{
				"userId":          userId,
				"connectedUserId": connectedUserId,
			})
		if err != nil {
			return nil, err
		}

		if res.Next() {
			connection = model.Connection{
				UserId:               userId,
				ConnectedUserId:      connectedUserId,
				IsConnected:          res.Record().Values[0].(bool),
				IsApprovedConnection: res.Record().Values[1].(bool),
				PendingConnection:    res.Record().Values[2].(bool),
			}
			return nil, nil
		}
		return nil, res.Err()

	})

	if err != nil {
		return nil, err
	}
	return &connection, nil
}
