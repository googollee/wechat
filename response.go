package wechat

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"time"
)

type Response interface {
	ReplyText(content string) error
}

type defaultResponse struct {
	w      http.ResponseWriter
	header MessageHeader
}

func (resp *defaultResponse) ReplyText(content string) error {
	msg := Message{
		MessageHeader: resp.header,
		Content:       content,
	}
	msg.CreateTime = time.Now().Unix()
	msg.MsgType = MsgText

	buf := bytes.NewBuffer(nil)
	buf.WriteString("<xml>")
	encoder := xml.NewEncoder(buf)
	if err := encoder.Encode(msg); err != nil {
		return err
	}
	buf.WriteString("</xml>")
	resp.w.Header().Set("Content-Type", "application/xml")
	if _, err := resp.w.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}
