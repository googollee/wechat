package wechat

type VoiceHandler func(w Response, r *Request)

type voiceNode struct {
	handler VoiceHandler
}

func (wc *Wechat) Voice(handler VoiceHandler) error {
	wc.voiceNode = voiceNode{handler}
	return nil
}

func (wc *Wechat) handleVoiceMessage(resp Response, req *Request) {
	wc.voiceNode.handler(resp, req)
}
