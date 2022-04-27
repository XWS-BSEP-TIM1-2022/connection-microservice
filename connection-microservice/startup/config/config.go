package config

import (
	"os"
	"time"
)

type Config struct {
	Port                  string
	ConnectionDBURI       string
	ConnectionDBUsername  string
	ConnectionDBPassword  string
	ConnectionServiceName string
	ExpiresIn             time.Duration
	UserServiceHost       string
	UserServicePort       string
}

func NewConfig() *Config {
	return &Config{
		Port:                  getEnv("CONNECTION_SERVICE_PORT", "8087"),
		ConnectionDBURI:       getEnv("CONNECTION_DB_URI", "neo4j+s://ac87e36d.databases.neo4j.io"),
		ConnectionDBUsername:  getEnv("CONNECTION_DB_USERNAME", "neo4j"),
		ConnectionDBPassword:  getEnv("CONNECTION_DB_PASSWORD", "I7InmmqDyQoT4BhAF5iXOCDB-EQ3wh-hcJn2-8QSobY"),
		ConnectionServiceName: getEnv("CONNECTION_SERVICE_NAME", "connection_service"),
		UserServiceHost:       getEnv("USER_SERVICE_HOST", "localhost"),
		UserServicePort:       getEnv("USER_SERVICE_PORT", "8085"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
