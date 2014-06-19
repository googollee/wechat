package wechat

import (
	"github.com/googollee/go-assert"
	"testing"
	"time"
)

func TestSession(t *testing.T) {
	ses := NewSession(time.Second)
	ses.Set("user", "key", 1)
	time.Sleep(time.Second * 2)
	value := ses.Get("user", "key")
	assert.Equal(t, value, nil)
}
