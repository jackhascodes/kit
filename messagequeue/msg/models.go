package msg

type Log interface {
	Error(v ...interface{})
}

type Message struct {
	Id         string
	Topic      string
	Body       []byte
	ReplyTopic string
}

type Config struct {
	Host     string
	ClientId string
}

type Handler func(*Message)
