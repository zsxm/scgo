package cache

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
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

func Init(conf Config) {
	if pool == nil {
		pool = &redis.Pool{
			MaxIdle:     conf.MaxIdle,
			MaxActive:   conf.MaxActive,
			IdleTimeout: conf.IdleTimeout,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", conf.Address)
				if err != nil {
					log.Println(err.Error())
					return nil, err
				}
				if conf.Password != "" {
					if _, err := conn.Do("AUTH", conf.Password); err != nil {
						conn.Close()
						log.Println(err.Error())
						return nil, err
					}
				}
				return conn, err
			},
			TestOnBorrow: func(conn redis.Conn, t time.Time) error {
				_, err := conn.Do("PING")
				if err != nil {
					log.Println(err.Error())
				}
				return err
			},
		}
	}
}

func send(cmd string, args ...interface{}) error {
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		log.Println(err.Error())
		return err
	}
	err := conn.Send(cmd, args...)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return conn.Flush()
}

func HSet(key, field, value string) error {
	err := send("HSET", key, field, value)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func HGet(key, field string) (string, error) {
	v, err := redis.String(do("HGET", key, field))
	if err != nil {
		log.Println(err.Error())
	}
	return v, err
}

func HDelete(key, field string) error {
	err := send("HDEL", key, field)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func HExists(key, field string) bool {
	v, err := redis.Bool(do("HEXISTS", key, field))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return v
}

func HLen(key string) (int, error) {
	v, err := redis.Int(do("HLEN", key))
	if err != nil {
		log.Println(err.Error())
	}
	return v, err
}

func Set(key, val string) error {
	_, err := do("SET", key, val)
	return err
}

func Get(key string) (string, error) {
	v, err := do("GET", key)
	return string(v.([]byte)), err
}

func do(cmd string, args ...interface{}) (interface{}, error) {
	conn := pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}
	v, err := conn.Do(cmd, args...)
	if err != nil {
		log.Println(err.Error())
	}
	return v, err
}
