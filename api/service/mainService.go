package service

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/nsqio/go-nsq"
	"net/http"
	"strconv"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/gnsq"
	"tokopedia.se.training/Project1/usermanagesys/api/response"
)

var nsqPublisher *nsq.Producer
var nsqModule *gnsq.GNSQModule

type MainService struct {
	configuration *configuration.Configuration
}

func NewMainService(config *configuration.Configuration) *MainService {
	return &MainService{
		configuration: config,
	}
}

func (service *MainService) InitNsqPublisher(nsqPub *nsq.Producer, nsqMod *gnsq.GNSQModule) {
	nsqPublisher = nsqPub
	nsqModule = nsqMod
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

func (s *MainService) GetVisitorCount(w http.ResponseWriter, r *http.Request) {

	resWriter := response.New(w)

	defer func() {
		if err := recover(); err != nil {
			resWriter.WriteError(fmt.Sprintf("%v", err))
			return
		}
	}()

	strCounter := nsqModule.GetRedisNSQ(gnsq.TOPIC_VISITOR_COUNTER)
	counter, err := strconv.Atoi(strCounter)
	if err != nil {
		counter = 0;
	}

	resCounter := struct {
		VisitorCount int `json:"visitor_count"`
	}{VisitorCount: counter}

	resWriter.WriteSuccess(resCounter)
	return
}
