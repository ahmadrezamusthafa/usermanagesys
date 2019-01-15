package main

import (
	"database/sql"
	"tokopedia.se.training/Project1/usermanagesys/api"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/gnsq"
	"tokopedia.se.training/Project1/usermanagesys/api/gnsq/consumer"
	"tokopedia.se.training/Project1/usermanagesys/api/repository"
	"tokopedia.se.training/Project1/usermanagesys/api/service"
)

var config *configuration.Configuration

type Application struct {
	Server api.Server
}

func (a *Application) connectDefaultDatabase(config *configuration.Configuration) (*sql.DB, error) {
	return sql.Open("postgres", config.ConnectionString.Default)
}

func (a *Application) serverInit() {
	configuration, err := configuration.NewConfiguration()
	if err != nil {
		panic(err.Error())
	}
	config = configuration

	db, err := a.connectDefaultDatabase(configuration)
	if err != nil {
		panic(err.Error())
	}

	userRepository := repository.NewTokopediaUserRepository(db)
	userService := service.NewUserService(configuration, userRepository)

	nsqModule := gnsq.NewNSQModule(configuration)
	if nsqModule.Configuration.NSQ.Enabled == true {
		nsqModule.InitNSQProducer()
	}

	server := api.NewServer(configuration, nsqModule, userService)
	a.Server = *server
}

func main() {
	a := Application{}
	a.serverInit()

	//NSQ check
	if a.Server.NsqModule.Configuration.NSQ.Enabled == true {
		//test
		/*err := api.SERVER.NsqModule.Producer["Server1"].Publish("TOPIC1", []byte("ini paket yang dikirim"))
		if err != nil {
			panic(err)
		}*/

		//add cunsomer
		consumer.InitConsumer()
	}

	a.Server.Run()
}
