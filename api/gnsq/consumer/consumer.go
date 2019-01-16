package consumer

import (
	"context"
	"fmt"
	"github.com/nsqio/go-nsq"
	"github.com/opentracing/opentracing-go"
	"log"
	"strconv"
	"tokopedia.se.training/Project1/usermanagesys/api"
	. "tokopedia.se.training/Project1/usermanagesys/api/gnsq"
)

var svr *api.Server

func InitConsumer() {
	svr = &api.SERVER

	//tambah consumer
	svr.NsqModule.AddConsumer([]string{TOPIC_VISITOR_COUNTER}, CHANNEL_VISITOR_COUNTER, listener)
}

func listener(message *nsq.Message, topic string) bool {
	span, _ := opentracing.StartSpanFromContext(context.Background(), "NSQConsumer."+TOPIC_VISITOR_COUNTER)
	defer span.Finish()

	finalMsg := string(message.Body)
	log.Println(fmt.Sprintf("data : %s", finalMsg))

	if topic == TOPIC_VISITOR_COUNTER && finalMsg == "OK" {
		strCounter := svr.NsqModule.GetRedisNSQ(TOPIC_VISITOR_COUNTER)
		counter, err := strconv.Atoi(strCounter)
		if err != nil {
			counter = 0;
		}
		counter++
		svr.NsqModule.SetRedisNSQ(TOPIC_VISITOR_COUNTER, 0, fmt.Sprintf("%d", counter))
	}

	//debug(fmt.Sprintf("RECEIV topic %s: %s", topic, message.Body))
	/*err := jsoniter.Unmarshal(message.Body, &options)
	if err != nil {
		contextlib.PrintErrorCtx(ctx, err, "LogIris")
		return false
	}*/

	return false
}
