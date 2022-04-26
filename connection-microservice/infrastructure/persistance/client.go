package persistance

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func GetDriver(uri, username, password string) (neo4j.Driver, error) {
	auth := neo4j.BasicAuth(username, password, "")

	driver, err := neo4j.NewDriver(uri, auth)
	if err != nil {
		return nil, err
	}

	return driver, nil
}
