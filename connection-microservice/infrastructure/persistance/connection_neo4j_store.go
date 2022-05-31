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

func (store *ConnectionNeo4jStore) CreateConnection(ctx context.Context, connection *model.Connection) (*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MERGE (user:User {userId:$userId}) "+
			"MERGE (connectedUser:User {userId:$connectedUserId}) "+
			"MERGE (user)-[c:CONNECT {isConnected:$isConnected, pendingConnection:$pendingConnection, isMessageNotificationEnabled:$isMessageNotificationEnabled, isPostNotificationEnabled:$isPostNotificationEnabled, isCommentNotificationEnabled:$isCommentNotificationEnabled}]->(connectedUser) RETURN c",
			map[string]interface{}{
				"userId":                       connection.UserId,
				"connectedUserId":              connection.ConnectedUserId,
				"isConnected":                  connection.IsConnected,
				"pendingConnection":            connection.PendingConnection,
				"isMessageNotificationEnabled": true,
				"isPostNotificationEnabled":    true,
				"isCommentNotificationEnabled": true,
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
		return nil, err
	}

	return connection, nil
}

func (store *ConnectionNeo4jStore) UpdateConnection(ctx context.Context, connection *model.Connection) (*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MATCH (user {userId:$userId})-[c:CONNECT]->(connectedUser {userId:$connectedUserId}) "+
			"SET c.isConnected=$isConnected, c.pendingConnection=$pendingConnection , c.isMessageNotificationEnabled=$isMessageNotificationEnabled , c.isPostNotificationEnabled=$isPostNotificationEnabled , c.isCommentNotificationEnabled=$isCommentNotificationEnabled "+
			"RETURN c",
			map[string]interface{}{
				"userId":                       connection.UserId,
				"connectedUserId":              connection.ConnectedUserId,
				"isConnected":                  connection.IsConnected,
				"pendingConnection":            connection.PendingConnection,
				"isMessageNotificationEnabled": connection.IsMessageNotificationEnabled,
				"isPostNotificationEnabled":    connection.IsPostNotificationEnabled,
				"isCommentNotificationEnabled": connection.IsCommentNotificationEnabled,
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
		return nil, err
	}

	return connection, nil
}

func (store *ConnectionNeo4jStore) DeleteConnection(ctx context.Context, userId string, connectedUserId string) error {
	span := tracer.StartSpanFromContext(ctx, "DeleteConnection")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		_, err := transaction.Run("MATCH (user {userId:$userId})-[c:CONNECT]->(connectedUser {userId:$connectedUserId}) "+
			"DELETE c",
			map[string]interface{}{
				"userId":          userId,
				"connectedUserId": connectedUserId,
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

func (store *ConnectionNeo4jStore) GetConnectionByUsersId(ctx context.Context, userId string, connectedUserId string) (*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnectionByUsersId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var connection = model.Connection{}
	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run("MATCH (user {userId:$userId})-[c:CONNECT]->(connectedUser {userId:$connectedUserId}) "+
			"RETURN c.isConnected, c.pendingConnection, c.isMessageNotificationEnabled, c.isPostNotificationEnabled, c.isCommentNotificationEnabled",
			map[string]interface{}{
				"userId":          userId,
				"connectedUserId": connectedUserId,
			})
		if err != nil {
			return nil, err
		}

		if res.Next() {
			connection = model.Connection{
				UserId:                       userId,
				ConnectedUserId:              connectedUserId,
				IsConnected:                  res.Record().Values[0].(bool),
				PendingConnection:            res.Record().Values[1].(bool),
				IsMessageNotificationEnabled: res.Record().Values[2].(bool),
				IsPostNotificationEnabled:    res.Record().Values[3].(bool),
				IsCommentNotificationEnabled: res.Record().Values[4].(bool),
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

func (store *ConnectionNeo4jStore) GetAllConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cypher := "MATCH (user {userId:$userId})-[c:CONNECT]->(connectedUser) " +
		"RETURN user.userId, connectedUser.userId, c.isConnected, c.pendingConnection"

	params := map[string]interface{}{
		"userId": userId,
	}

	connections, err := store.GetConnections(ctx, cypher, params)

	if err != nil {
		return nil, err
	}

	cypher = "MATCH (user)-[c:CONNECT]->(connectedUser {userId:$userId}) " +
		"RETURN user.userId, connectedUser.userId, c.isConnected, c.pendingConnection"

	params = map[string]interface{}{
		"userId": userId,
	}

	newConnections, err := store.GetConnections(ctx, cypher, params)

	if err != nil {
		return nil, err
	}

	connections = append(connections, newConnections...)

	return connections, nil
}

func (store *ConnectionNeo4jStore) GetFollowings(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetFollowings")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cypher := "MATCH (user {userId:$userId})-[c:CONNECT {isConnected:true}]->(connectedUser) " +
		"RETURN user.userId, connectedUser.userId, c.isConnected, c.pendingConnection"

	params := map[string]interface{}{
		"userId": userId,
	}

	connections, err := store.GetConnections(ctx, cypher, params)

	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (store *ConnectionNeo4jStore) GetFollowers(ctx context.Context, connectedUserId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetFollowers")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cypher := "MATCH (user)-[c:CONNECT {isConnected:true}]->(connectedUser {userId:$connectedUserId}) " +
		"RETURN user.userId, connectedUser.userId, c.isConnected, c.pendingConnection, c.isMessageNotificationEnabled, c.isPostNotificationEnabled, c.isCommentNotificationEnabled"

	params := map[string]interface{}{
		"connectedUserId": connectedUserId,
	}

	connections, err := store.GetConnections(ctx, cypher, params)

	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (store *ConnectionNeo4jStore) GetAllRequestConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllRequestConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cypher := "MATCH (user)-[c:CONNECT {isConnected:false, pendingConnection:true}]->(connectedUser {userId:$connectedUserId}) " +
		"RETURN user.userId, connectedUser.userId, c.isConnected, c.pendingConnection"

	params := map[string]interface{}{
		"connectedUserId": userId,
	}

	connections, err := store.GetConnections(ctx, cypher, params)

	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (store *ConnectionNeo4jStore) GetAllPendingConnectionsByUserId(ctx context.Context, userId string) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllPendingConnectionsByUserId")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	cypher := "MATCH (user {userId:$userId})-[c:CONNECT {isConnected:false, pendingConnection:true}]->(connectedUser) " +
		"RETURN user.userId, connectedUser.userId, c.isConnected, c.pendingConnection"

	params := map[string]interface{}{
		"userId": userId,
	}

	connections, err := store.GetConnections(ctx, cypher, params)

	if err != nil {
		return nil, err
	}

	return connections, nil
}

func (store *ConnectionNeo4jStore) GetConnections(ctx context.Context, cypher string, params map[string]interface{}) ([]*model.Connection, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnections")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	session := store.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	var connection []*model.Connection
	_, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		res, err := transaction.Run(cypher, params)
		if err != nil {
			return nil, err
		}

		for res.Next() {
			connection = append(connection, &model.Connection{
				UserId:                       res.Record().Values[0].(string),
				ConnectedUserId:              res.Record().Values[1].(string),
				IsConnected:                  res.Record().Values[2].(bool),
				PendingConnection:            res.Record().Values[3].(bool),
				IsMessageNotificationEnabled: res.Record().Values[4].(bool),
				IsPostNotificationEnabled:    res.Record().Values[5].(bool),
				IsCommentNotificationEnabled: res.Record().Values[6].(bool),
			})
		}
		return nil, res.Err()

	})

	if err != nil {
		return nil, err
	}
	return connection, nil
}
