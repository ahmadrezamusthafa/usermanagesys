package gnsq

type GNSQData struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

const (
	TOPIC_VISITOR_COUNTER string = "Topic.Visitor.Counter"
)

const (
	CHANNEL_VISITOR_COUNTER = TOPIC_VISITOR_COUNTER + ".CH"
)
