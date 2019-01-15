package main

import (
	"database/sql"
	"tokopedia.se.training/Project1/usermanagesys/api"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/gnsq"
	"tokopedia.se.training/Project1/usermanagesys/api/repository"
	"tokopedia.se.training/Project1/usermanagesys/api/service"
)

var config *configuration.Configuration


func connectDefaultDatabase(config *configuration.Configuration) (*sql.DB, error) {
	return sql.Open("postgres", config.ConnectionString.Default)
}

func serverInit() {
	configuration, err := configuration.NewConfiguration()
	if err != nil {
		panic(err.Error())
	}
	config = configuration

	db, err := connectDefaultDatabase(configuration)
	if err != nil {
		panic(err.Error())
	}

	userRepository := repository.NewTokopediaUserRepository(db)
	userService := service.NewUserService(configuration, userRepository)

	nsqModule := gnsq.NewNSQModule(configuration)
	if nsqModule.Configuration.NSQ.Enabled == true {
		nsqModule.InitNSQProducer()

		//add cunsomer
		//consumer.InitConsumer()
	}

	server := api.NewServer(configuration, nsqModule, userService)
	server.Run()
}

func main() {
	serverInit()
}
