package gnsq

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nsqio/go-nsq"
	"log"
	"strings"
	"time"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration"
	"tokopedia.se.training/Project1/usermanagesys/api/configuration/dto"
	. "tokopedia.se.training/Project1/usermanagesys/api/redis"
)

const (
	WIB string = "Asia/Jakarta"
)

type GNSQModule struct {
	Configuration *dto.NSQDto
	Producer      map[string]*nsq.Producer
}

var redisPool *redis.Pool

func NewNSQModule(config *configuration.Configuration) *GNSQModule {
	nsqConfig := &config.NSQ
	module := GNSQModule{
		Configuration: nsqConfig,
	}
	if nsqConfig.NSQ.Enabled {
		redisConnectNSQ(nsqConfig.Redis.Url, nsqConfig.Redis.MaxActive, nsqConfig.Redis.MaxIdle)
		log.Println("[NSQ]", nsqConfig)
	}

	return &module
}

func redisConnectNSQ(uri string, maxActive, maxIdle int) {
	redisPool = InitRedisPool(uri, "", maxActive, maxIdle)
}

func (a *GNSQModule) InitNSQProducer() error {
	var err error

	a.Producer = make(map[string]*nsq.Producer)

	for key := range a.Configuration.NSQD {
		a.Producer[key], err = nsq.NewProducer(
			fmt.Sprintf("%s:%s", a.Configuration.NSQD[key].NsqdIP, a.Configuration.NSQD[key].NsqdTCPPort),
			nsq.NewConfig(),
		)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

type ServiceHandler func(*nsq.Message, string) bool

func (a *GNSQModule) AddConsumer(topics []string, channel string, serviceHandler ServiceHandler) error {
	config := nsq.NewConfig()
	config.MaxInFlight = 5
	config.MaxAttempts = 10

	for i := range topics {
		q, err := nsq.NewConsumer(topics[i], channel, config)
		if err != nil {
			log.Println(err)
			return err
		}

		currTopic := topics[i]
		q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			if getTimestampDifferentMinute(message.Timestamp) > a.Configuration.NSQLookupd.TimeLimitRequeue {
				message.Finish()
			}

			log.Println(fmt.Sprintf("terima %s", message.Body))

			isRequeue := serviceHandler(message, currTopic)

			if isRequeue == true {
				message.Requeue(-1)
			} else {
				a.SetRedisNSQ(channel, 600, "1")
				message.Finish()
			}

			return nil
		}))

		err = q.ConnectToNSQLookupd(a.Configuration.NSQLookupd.Url + ":" + a.Configuration.NSQLookupd.HTTPPort)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func GetKeyRedisNSQ(channel string) string {
	key := fmt.Sprintf("nsq:%s", strings.ToLower(channel))
	return key
}

func (a *GNSQModule) DeleteRedisNSQ(channel string) {
	c := redisPool.Get()
	defer c.Close()

	key := GetKeyRedisNSQ(channel)

	_, err := redis.Int(c.Do("DEL", key))
	if err != nil {
		log.Println(err)
	}
}

func (a *GNSQModule) GetRedisNSQ(channel string) string {
	c := redisPool.Get()
	defer c.Close()

	key := GetKeyRedisNSQ(channel)

	value, err := redis.String(c.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return ""
	}

	fmt.Println("GET", key)
	return value
}

func (a *GNSQModule) SetRedisNSQ(channel string, expire int, value string) {
	c := redisPool.Get()
	defer c.Close()

	key := GetKeyRedisNSQ(channel)

	var t int

	if expire > 0 {
		t = expire
	} else {
		t = int(time.Now().Unix()) + 24*60*60
	}

	redis.String(c.Do("SETEX", key, t, value))
	fmt.Println("SETEX", key, value)
}

func getTimestampDifferentMinute(timeStamp int64) (minutes int) {
	if timeStamp > 0 {
		timeConvert := time.Unix(0, timeStamp).In(GetLocation())
		timeNow := GetTimeNow()
		diff := timeNow.Sub(timeConvert)
		minutes = int(diff.Minutes())
	}

	return minutes
}

func GetTimeNow() time.Time {
	return time.Now().In(GetLocation())
}

func GetTimeWithLocation(t time.Time) time.Time {
	return t.In(GetLocation())
}

func GetLocation() *time.Location {
	return time.FixedZone(WIB, 7*3600)
}
