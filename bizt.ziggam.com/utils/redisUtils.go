package utils

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	pool *redis.Pool
	// redisServer = flag.String("redisServer", "localhost:6379", "wlrrka1")
)

type RedisController struct {
	pool *redis.Pool
}

var RPool RedisController

func (redigo *RedisController) NewPool(server string, password string) *redis.Pool {
	redigo.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			if _, err := c.Do("SELECT", 0); err != nil {
				c.Close()
				return nil, err
			}

			return c, nil
		},
	}

	return redigo.pool
}

func (redigo *RedisController) Ping() error {
	conn := redigo.pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("PING"))
	if err != nil {
		return err
	}

	fmt.Printf("Redis Connect Ok!!!!!!!!   PING Response = %s\n", s)
	return nil
}

func (redigo *RedisController) HSet(key string, hashkey string, value string) *redis.Pool {
	conn := redigo.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("HSET", key, hashkey, value); err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func (redigo *RedisController) Expire(key string, timeout int) *redis.Pool {
	conn := redigo.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("EXPIRE", key, timeout); err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func (redigo *RedisController) HGet(key string, hashkey string) (string, error) {
	conn := redigo.pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("HGET", key, hashkey))
	if err != nil {
		return s, err
	}

	return s, err
}
