package wechat

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

func (wc *Wechat) getAccessToken() (string, error) {
	wc.accessTokenLocker.Lock()
	defer wc.accessTokenLocker.Unlock()

	now := time.Now().Unix()
	if wc.accessExpiresAt > now {
		return wc.accessToken, nil
	}
	query := make(url.Values)
	query.Set("grant_type", "client_credential")
	query.Set("appid", wc.appId)
	query.Set("secret", wc.appSecret)

	u := "https://api.weixin.qq.com/cgi-bin/token?" + query.Encode()
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		var e Error
		if err := decoder.Decode(&e); err != nil {
			return "", err
		}
		return "", e
	}
	var token AccessToken
	if err := decoder.Decode(&token); err != nil {
		return "", err
	}
	wc.accessToken = token.Token
	wc.accessExpiresAt = now + int64(token.ExpiresIn)
	return wc.accessToken, nil
}

type Error struct {
	code int    `json:"errcode"`
	msg  string `json:"errmsg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d)%s", e.code, e.msg)
}

func (wc *Wechat) HttpRequest(method, urlStr string, query url.Values, body io.Reader, resp interface{}) error {
	token, err := wc.getAccessToken()
	if err != nil {
		return err
	}
	if query == nil {
		query = make(url.Values)
	}
	query.Set("access_token", token)
	req, err := http.NewRequest(method, urlStr+"?"+query.Encode(), body)
	if err != nil {
		return err
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if r.StatusCode != 200 {
		var e Error
		if err := decoder.Decode(&e); err != nil {
			return err
		}
		return e
	}
	if err := decoder.Decode(resp); err != nil {
		return err
	}
	return nil
}
