package wechat

import (
	"io"
	"net/url"
)

func (wc *Wechat) GetMediaUrl(mediaId string) (string, error) {
	token, err := wc.getAccessToken()
	if err != nil {
		return "", err
	}
	query := make(url.Values)
	query.Set("access_token", token)
	query.Set("media_id", mediaId)
	u := "http://file.api.weixin.qq.com/cgi-bin/media/get?" + query.Encode()
	return u, nil
}
