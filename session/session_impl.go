package session

import (
	"errors"
	"net/http"
	"strings"

	"github.com/zsxm/scgo/data/cache"
	"github.com/zsxm/scgo/log"
	"github.com/zsxm/scgo/tools/uuid"
)

func New(w http.ResponseWriter, r *http.Request, o *Options) Interface {
	// ignore error -> http: named cookie not present
	cookie, _ := r.Cookie(cookieKey)
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

func (this *session) Get(key string) (string, error) {
	v, err := cache.HGet(this.key, key)
	if err != nil {
		log.Error(err)
		return "", err
	}
	this.expire(this.options.MaxAge)
	return v, nil
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

func (this *session) AddFlash(value string, vars ...string) {
}

func (this *session) Flashes(vars ...string) []string {
	return []string{}
}

func (this *session) Options(options *Options) {
	this.options = options
}

func (this *session) IsLogin() bool {
	if !this.exists(principalId) {
		return false
	}
	if id, err := this.Get(principalId); err == nil {
		log.Info("----------------------", id)
		if id != "" {
			return true
		}
	}
	return false
}

func (this *session) Principal() (Principal, error) {
	p := Principal{}
	if !this.exists(principalId) {
		return p, errors.New("Not logged in")
	}
	id, err := this.Get(principalId)
	if err != nil {
		return p, err
	}
	name, err := this.Get(principalName)
	if err != nil {
		return p, err
	}
	loginName, err := this.Get(principalLoginName)
	if err != nil {
		return p, err
	}
	loginTime, err := this.Get(principalLoginTime)
	if err != nil {
		return p, err
	}
	permissions, err := this.Get(principalPermissions)
	if err != nil {
		return p, err
	}
	p.Id = id
	p.Name = name
	p.LoginName = loginName
	p.LoginTime = loginTime
	p.Permissions = permissions
	return p, nil
}

func (this *session) SetPrincipal(value Principal) error {
	if err := this.Set(principalId, value.Id); err != nil {
		return err
	}
	if err := this.Set(principalName, value.Name); err != nil {
		return err
	}
	if err := this.Set(principalLoginName, value.LoginName); err != nil {
		return err
	}
	if err := this.Set(principalLoginTime, value.LoginTime); err != nil {
		return err
	}
	if err := this.Set(principalPermissions, value.Permissions); err != nil {
		return err
	}
	return nil
}

func (this *session) ResetPrincipal() error {
	if err := this.Delete(principalId); err != nil {
		return err
	}
	if err := this.Delete(principalName); err != nil {
		return err
	}
	if err := this.Delete(principalLoginName); err != nil {
		return err
	}
	if err := this.Delete(principalLoginTime); err != nil {
		return err
	}
	if err := this.Delete(principalPermissions); err != nil {
		return err
	}
	return nil
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
