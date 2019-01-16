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
)

const (
	WIB string = "Asia/Jakarta"
)

type GNSQModule struct {
	Configuration         *dto.NSQDto
	Producer              map[string]*nsq.Producer
	ChEmailRefund         chan *nsq.Message
	ChReactivateVoucher   chan *nsq.Message
	ChReactivateMPVoucher chan *nsq.Message
	ChCreateVoucher       chan *nsq.Message
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
	redisPool = NewPool(uri, "", maxActive, maxIdle)
}

func NewPool(server, password string, maxActive, maxIdle int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 240 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
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
			duplicate := GetRedisNSQ(channel, string(message.ID[:]))

			if getTimestampDifferentMinute(message.Timestamp) > a.Configuration.NSQLookupd.TimeLimitRequeue {
				message.Finish()
			}

			if duplicate != "1" {
				log.Println(fmt.Sprintf("Got a message: %s", message.Body))

				isRequeue := serviceHandler(message, currTopic)

				if isRequeue == true {
					message.Requeue(-1)
				} else {
					SetRedisNSQ(channel, string(message.ID[:]), 600, "1")
					message.Finish()
				}

			} else {
				log.Println(channel, string(message.ID[:]), ", is duplicate")
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

func GetKeyRedisNSQ(types string, id string) string {
	key := fmt.Sprintf("nsq:%s:%s", strings.ToLower(types), strings.ToLower(id))
	return key
}

func DeleteRedisNSQ(types string, id string) {
	c := redisPool.Get()
	defer c.Close()

	key := GetKeyRedisNSQ(types, id)

	_, err := redis.Int(c.Do("DEL", key))
	if err != nil {
		log.Println(err)
	}
}

func GetRedisNSQ(types string, id string) string {
	c := redisPool.Get()
	defer c.Close()

	key := GetKeyRedisNSQ(types, id)

	value, err := redis.String(c.Do("GET", key))
	if err != nil && err != redis.ErrNil {
		return ""
	}

	return value
}

func SetRedisNSQ(types string, id string, expire int, value string) {
	c := redisPool.Get()
	defer c.Close()

	key := GetKeyRedisNSQ(types, id)

	var t int

	if expire > 0 {
		t = expire
	} else {
		t = int(time.Now().Unix()) + 24*60*60
	}

	redis.String(c.Do("SETEX", key, t, value))
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
