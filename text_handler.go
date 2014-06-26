package wechat

import (
	"bytes"
	"encoding/json"
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

func (wc *Wechat) SendTextMessage(to, text string) error {
	type Message struct {
		ToUser  string `json:"touser"`
		MsgType MsgType
		Text    struct {
			Content string `json:"content"`
		} `json:text"`
	}
	msg := Message{
		ToUser:  to,
		MsgType: MsgText,
	}
	msg.Text.Content = text
	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	if err := e.Encode(msg); err != nil {
		return err
	}
	var resp Error
	if err := wc.HttpRequest("POST", "https://api.weixin.qq.com/cgi-bin/message/custom/send", nil, buf, &resp); err != nil {
		return err
	}
	if resp.Code != 0 {
		return resp
	}
	return nil
}
