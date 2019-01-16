package gnsq

type GNSQTopic struct {
	Topic string
}

const (
	TOPIC_VISITOR_COUNTER string = "Topic.Visitor.Counter"
)

const (
	CHANNEL_VISITOR_COUNTER string = TOPIC_VISITOR_COUNTER + ".CH"
)
