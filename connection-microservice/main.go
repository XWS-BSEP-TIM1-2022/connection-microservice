package main

import (
	"connection-microservice/startup"
	cfg "connection-microservice/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
	defer server.Stop()
}
