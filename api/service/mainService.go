package service

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/nsqio/go-nsq"
	"net/http"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/gnsq"
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

	pubData := gnsq.GNSQData{}
	err := jsoniter.Unmarshal([]byte(data), &pubData)
	if err != nil {
		resWriter.WriteError("Unmarshal failed")
		return
	}

	//publish nsq
	nsqPublisher.Publish(pubData.Topic, []byte(pubData.Message))

	resWriter.WriteSuccess("Successfully publish")
	return
}
