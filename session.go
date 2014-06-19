package wechat

import (
	"time"
)

type Session interface {
	Set(user, key string, data interface{})
	Get(user, key string) interface{}
	Remove(user, key string)
}

type sessionStorage struct {
	lastModified int64
	data         map[string]interface{}
}

type session struct {
	duration time.Duration
	storages map[string]sessionStorage
}

func NewSession(duration time.Duration) Session {
	return &session{
		duration: duration,
		storages: make(map[string]sessionStorage),
	}
}

func (ses *session) Set(user, key string, data interface{}) {
	storage, ok := ses.storages[user]
	if !ok {
		storage = sessionStorage{
			data: make(map[string]interface{}),
		}
	}
	storage.lastModified = time.Now().Unix()
	storage.data[key] = data
	ses.storages[user] = storage
}

func (ses *session) Get(user, key string) interface{} {
	storage, ok := ses.storages[user]
	if !ok {
		return nil
	}
	if (storage.lastModified + int64(ses.duration/time.Second)) < time.Now().Unix() {
		delete(ses.storages, user)
		return nil
	}
	return storage.data[key]
}

func (ses *session) Remove(user, key string) {
	storage, ok := ses.storages[user]
	if !ok {
		return
	}
	if (storage.lastModified + int64(ses.duration)) < time.Now().Unix() {
		delete(ses.storages, user)
		return
	}
	delete(storage.data, key)
	if len(storage.data) == 0 {
		delete(ses.storages, user)
	}
}
