package dto

type ConnectionStringDto struct {
	Default       string
	AccessControl string
}

type AppDto struct {
	ServerAddress   string
	BackEndAddress  string
	FrontEndAddress string
	Debug           bool
}

type NSQDto struct {
	NSQ struct {
		ChannelSize int
		Enabled     bool
	}
	Redis struct {
		Url       string
		MaxActive int
		MaxIdle   int
	}
	NSQLookupd struct {
		Url               string
		TCPPort           string
		HTTPPort          string
		TimeLimitRequeue  int
		MaxAttemptRequeue uint16
	}
	NSQD map[string]*struct {
		NsqdIP       string
		NsqdTCPPort  string
		NsqdHTTPPort string
	}
}
