package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Wechat struct {
	token      string
	appId      string
	appSecret  string
	ses        Session
	log        *log.Logger
	textNodes  []textNode
	eventNodes []eventNode
}

func New(token, appId, appSecret string) *Wechat {
	ret := &Wechat{
		token:     token,
		appId:     appId,
		appSecret: appSecret,
	}
	ret.SetLogger(nil)
	return ret
}

func (wc *Wechat) SetLogger(l *log.Logger) {
	if l == nil {
		wc.log = log.New(os.Stdout, "[wechat]", log.LstdFlags)
		return
	}
	wc.log = l
}

func (wc *Wechat) SetSession(ses Session) {
	wc.ses = ses
}

func (wc *Wechat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	switch r.Method {
	case "GET":
		wc.checkSignature(w, r)
		return
	case "POST":
		msg, err := ParseMessage(r.Body)
		if err != nil {
			wc.log.Printf("decode body error: %s", err)
			return
		}
		req, resp := wc.getRequestResponse(w, msg)
		switch msg.MsgType {
		case MsgText:
			wc.handleTextMessage(resp, req)
		case MsgEvent:
			wc.handleEventMessage(resp, req)
		}
	}

}

func (wc *Wechat) getRequestResponse(w http.ResponseWriter, msg Message) (*Request, Response) {
	req := &Request{
		Message: msg,
		ses:     wc.ses,
	}
	resp := &defaultResponse{
		w:      w,
		header: msg.MessageHeader,
	}
	resp.header.ToUserName, resp.header.FromUserName = resp.header.FromUserName, resp.header.ToUserName
	return req, resp
}

func (wc *Wechat) checkSignature(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	sign := query.Get("signature")
	timestamp := query.Get("timestamp")
	nonce := query.Get("nonce")
	echostr := query.Get("echostr")
	arr := []string{wc.token, timestamp, nonce}
	sort.Strings(arr)
	certStr := strings.Join(arr, "")
	h := sha1.New()
	h.Write([]byte(certStr))
	sum := hex.EncodeToString(h.Sum(nil))
	if sign != sum {
		wc.log.Printf("signature wrong: %s != %s%+v", sign, certStr, arr)
		http.Error(w, "cert error", 403)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(echostr))
}
