package wechat

type Request struct {
	Message Message
	ses     Session
}

func (r *Request) Header() MessageHeader {
	return r.Message.MessageHeader
}

func (r *Request) Set(key string, data interface{}) {
	r.ses.Set(r.Message.FromUserName, key, data)
}

func (r *Request) Get(key string) interface{} {
	r.ses.Get(r.Message.FromUserName, key)
}

func (r *Request) Remove(key string) {
	r.ses.Remove(r.Message.FromUserName, key)
}
