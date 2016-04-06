package session

import (
	"net/http"
	"strings"

	"github.com/zsxm/scgo/data"
	"github.com/zsxm/scgo/data/cache"
	"github.com/zsxm/scgo/log"
	"github.com/zsxm/scgo/tools/uuid"
)

func New(w http.ResponseWriter, r *http.Request, o *Options) Interface {
	// ignore error -> http: named cookie not present
	var cookie *http.Cookie
	//cookie, _ := r.Cookie(cookieKey)
	cookies := r.Cookies()
	if len(cookies) > 0 {
		cookie = cookies[0]
	}
	if cookie == nil {
		sid := strings.Replace(uuid.NewV4().String(), "-", "", -1)
		cookie = &http.Cookie{
			Name:     cookieKey,
			Value:    sid,
			Path:     o.Path,
			Domain:   o.Domain,
			Secure:   o.Secure,
			HttpOnly: o.HttpOnly,
		}
		http.SetCookie(w, cookie)
	}
	s := &session{
		id:      cookie.Value,
		key:     sessionPrefix + ":" + cookie.Value,
		options: o,
		request: r,
	}
	return s
}

type session struct {
	id      string
	key     string
	options *Options
	request *http.Request
}

func (this *session) Id() string {
	return this.id
}

func (this *session) GetMap() (data.Map, error) {
	r, err := cache.HGetMap(this.key)
	if err != nil {
		log.Error(err)
		return r, err
	}
	this.expire(this.options.MaxAge)
	return r, nil
}

func (this *session) Get(key string) (string, error) {
	v, err := cache.HGet(this.key, key)
	if err != nil {
		log.Error(err)
		return "", err
	}
	this.expire(this.options.MaxAge)
	return v, nil
}

func (this *session) SetEntity(entity data.EntityInterface) error {
	err := cache.HSetEntity(this.key, entity)
	if err != nil {
		log.Error(err)
		return err
	}
	this.expire(this.options.MaxAge)
	return nil
}

func (this *session) SetMap(value map[string]string) error {
	err := cache.HSetMap(this.key, value)
	if err != nil {
		log.Error(err)
		return err
	}
	this.expire(this.options.MaxAge)
	return nil
}

func (this *session) Set(key string, val string) error {
	err := cache.HSet(this.key, key, val)
	if err != nil {
		log.Error(err)
		return err
	}
	this.expire(this.options.MaxAge)
	return nil
}

func (this *session) Delete(key string) error {
	err := cache.HDelete(this.key, key)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (this *session) Clear() error {
	err := cache.Delete(this.key)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (this *session) Options(options *Options) {
	this.options = options
}

func (this *session) exists(key string) bool {
	v := cache.HExists(this.key, key)
	return v
}

func (this *session) expire(second int) error {
	err := cache.Expire(this.key, second)
	if err != nil {
		return err
	}
	return nil
}
