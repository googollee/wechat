package wechat

import (
	"regexp"
)

type TextHandler func(w Response, r *Request)

type textNode struct {
	reg     *regexp.Regexp
	handler TextHandler
}

func (wc *Wechat) Text(pattern string, handler TextHandler) error {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	wc.textNodes = append(wc.textNodes, textNode{reg, handler})
	return nil
}

func (wc *Wechat) handleTextMessage(resp Response, req *Request) {
	for _, node := range wc.textNodes {
		if node.reg.MatchString(req.Message.Content) {
			node.handler(resp, req)
			return
		}
	}
}
