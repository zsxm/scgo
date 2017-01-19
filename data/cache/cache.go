package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/zsxm/scgo/data"
	"github.com/zsxm/scgo/log"
)

//redis配置
var Conf *Config

type Config struct {
	Address     string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var pool *redis.Pool

func Init(conf *Config) {
	if Conf == nil {
		Conf = conf
	}
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     conf.MaxIdle,
			MaxActive:   conf.MaxActive,
			IdleTimeout: conf.IdleTimeout,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", conf.Address)
				if err != nil {
					log.Error(err.Error())
					return nil, err
				}
				if conf.Password != "" {
					if _, err := conn.Do("AUTH", conf.Password); err != nil {
						conn.Close()
						log.Error(err.Error())
						return nil, err
					}
				}
				return conn, err
			},
			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				_, err := conn.Do("PING")
				if err != nil {
					log.Error(err.Error())
				} else {
					log.Info("Redis Connection ,IP Address:", Conf.Address, "\t[OK]")
				}
				return err
			},
		}
		pool.TestOnBorrow(pool.Get(), time.Now())
	}
}

func send(cmd string, args ...interface{}) error {
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		log.Error(err.Error())
		return err
	}
	err := conn.Send(cmd, args...)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return conn.Flush()
}

func HSet(key, field, value string) error {
	err := send("HSET", key, field, value)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func HSetMap(key string, value map[string]string) error {
	for k, v := range value {
		if err := HSet(key, k, v); err != nil {
			return err
		}
	}
	return nil
}

func HSetEntity(key string, entity data.EntityInterface) error {
	fields := entity.FieldNames()
	names := fields.Names()
	for _, v := range names {
		tmpField := entity.Field(v)
		if tmpField != nil {
			if err := HSet(key, v, tmpField.Value()); err != nil {
				return err
			}
		}
	}
	return nil
}

func HGet(key, field string) (string, error) {
	v, err := redis.String(do("HGET", key, field))
	if err != nil {
		log.Error(err.Error())
	}
	return v, err
}

func HGetMap(key string) (data.Map, error) {
	cvs, err := do("HGETALL", key)
	result := make(data.Map)

	if err != nil {
		log.Error(err.Error())
		return result, err
	}
	if vs, ok := cvs.([]interface{}); ok && cvs != nil {
		size := len(vs)
		for i := 0; i < size; i++ {
			result[string(vs[i].([]byte))] = string(vs[i+1].([]byte))
			i++
		}
	}
	return result, nil
}

func HGetEntity(key string, entity data.EntityInterface) error {
	cvs, err := do("HGETALL", key)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if vs, ok := cvs.([]interface{}); ok && cvs != nil {
		size := len(vs)
		for i := 0; i < size; i++ {
			tmpField := entity.Field(string(vs[i].([]byte)))
			if tmpField != nil {
				tmpField.SetValue(string(vs[i+1].([]byte)))
			}
			i++
		}
	}
	return nil
}

func HDelete(key, field string) error {
	err := send("HDEL", key, field)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func HExists(key, field string) bool {
	v, err := redis.Bool(do("HEXISTS", key, field))
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return v
}

func HLen(key string) (int, error) {
	v, err := redis.Int(do("HLEN", key))
	if err != nil {
		log.Error(err.Error())
	}
	return v, err
}

func Delete(key string) error {
	err := send("DEL", key)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func Exists(key string) bool {
	v, err := redis.Bool(do("EXISTS", key))
	if err != nil {
		log.Error(err.Error())
		return false
	}
	return v
}

func Expire(key string, second int) error {
	err := send("EXPIRE", key, second)
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func Set(key, val string) error {
	_, err := do("SET", key, val)
	return err
}

func Get(key string) (string, error) {
	v, err := do("GET", key)
	if v == nil {
		return "", err
	}
	return string(v.([]byte)), err
}

func TTL(key string) (int, error) {
	v, err := do("TTL", key)
	if err != nil {
		return -1, err
	}
	r := int(v.(int64))
	return r, nil
}

func do(cmd string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		log.Error(err.Error())
		return nil, err
	}
	v, err := conn.Do(cmd, args...)
	if err != nil {
		log.Error(err.Error())
	}
	return v, err
}
