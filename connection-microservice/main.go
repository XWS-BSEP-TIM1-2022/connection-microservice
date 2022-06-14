package main

import (
	"connection-microservice/application"
	"connection-microservice/startup"
	cfg "connection-microservice/startup/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var log = logrus.New()

func main() {
	configuringLog()
	log.Info("Server starting...")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
	defer server.Stop()
}

func configuringLog() {
	application.Log = log
	log.Out = os.Stdout
	path := "connection-microservice.log"
	/*Log rotation correlation function
	`Withlinkname 'establishes a soft connection for the latest logs
	`Withrotationtime 'sets the time of log splitting, and how often to split
	Only one of withmaxage and withrotationcount can be set
	 `Withmaxage 'sets the maximum save time before cleaning the file
	 `Withrotationcount 'sets the maximum number of files to be saved before cleaning
	*/
	//The following configuration logs rotate a new file every 1 minute, keep the log files of the last 3 minutes, and automatically clean up the surplus.
	writer, err := rotatelogs.New(
		path+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(8760)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)

	if err == nil {
		log.SetOutput(writer)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}
