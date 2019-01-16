package consumer

import (
	"context"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/nsqio/go-nsq"
	"github.com/opentracing/opentracing-go"
	"log"
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

	json, err := jsoniter.Marshal(context.Background())
	if err != nil {
		panic(err)
	}

	var ctx = context.Background()
	log.Println(ctx.Value(""))

	log.Println(fmt.Sprintf("json context: %s", json))
	//debug(fmt.Sprintf("RECEIV topic %s: %s", topic, message.Body))
	/*err := jsoniter.Unmarshal(message.Body, &options)
	if err != nil {
		contextlib.PrintErrorCtx(ctx, err, "LogIris")
		return false
	}*/

	return false
}
