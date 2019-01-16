package service

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"net/http"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/response"
)

var nsqPublisher *nsq.Producer

type MainService struct {
	configuration *configuration.Configuration
}

func NewMainService(config *configuration.Configuration) *MainService {
	return &MainService{
		configuration: config,
	}
}

func (service *MainService) InitNsqPublisher(nsqPub *nsq.Producer) {
	nsqPublisher = nsqPub
}

func (service *MainService) PublishNSQ(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resWriter := response.New(w)

	defer func() {
		if err := recover(); err != nil {
			resWriter.WriteError(fmt.Sprintf("%v", err))
			return
		}
	}()

	queryValues := r.URL.Query()
	data := queryValues.Get("data")
	if data == "" {
		resWriter.WriteError("Data is required")
		return
	}
	topic := queryValues.Get("topic")
	if data == "" {
		resWriter.WriteError("Topic is required")
		return
	}

	//publish nsq
	nsqPublisher.Publish(topic, []byte(data))

	resWriter.WriteSuccess("Successfully publish")
	return
}
