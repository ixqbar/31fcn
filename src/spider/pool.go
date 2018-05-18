package spider

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

func NewRedisPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				Logger.Print(err)
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				Logger.Print(err)
			}
			return err
		},
	}
}