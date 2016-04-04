package session

import (
	"github.com/zsxm/scgo/data"
)

var OptionsConfig *Options

type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

// Interface stores the values and optional configuration for a session.
type Interface interface {
	Get(key string) (string, error)
	Set(key string, val string) error
	Delete(key string) error
	Clear() error
	Options(*Options)
	SetMap(value map[string]string) error
	GetMap() (data.Map, error)
	SetEntity(entity data.EntityInterface) error
	Id() string
}
