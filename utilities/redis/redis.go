package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

var RedisConn *redis.Pool

// Setup Initialize the Redis instance
func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     viper.GetInt("redis_docker.max_idle"),
		MaxActive:   viper.GetInt("redis_docker.max_active"),
		IdleTimeout: viper.GetDuration("redis_docker.idle_timeout"),
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", viper.GetString("redis_docker.host"), viper.GetInt("redis_docker.port")))
			if err != nil {
				return nil, err
			}

			if viper.GetString("redis_docker.password") != "" {
				authArgs := make([]interface{}, 0, 2)
				if viper.GetString("redis_docker.username") != "" {
					authArgs = append(authArgs, viper.GetString("redis_docker.username"))
				}
				authArgs = append(authArgs, viper.GetString("redis_docker.password"))
				if _, err := c.Do("AUTH", authArgs...); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

// Set a key/value
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		// helper.Log("error")("log connection func set", err)
		return err
	}
	// helper.Log("debug")("log connection func set", value)

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}

	return nil
}

// Exists check a key
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		// helper.Log("error")("log connection func exists", err)

		return false
	}
	// helper.Log("debug")("log connection func exists", exists)

	return exists
}

// Get get a key
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		// helper.Log("error")("log connection func get", err)

		return nil, err
	}
	// helper.Log("debug")("log connection func get", reply)

	return reply, nil
}

// Delete delete a kye
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	// helper.Log("debug")("log connection func delete", conn)
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes batch delete
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		// helper.Log("error")("log connection func LikeDeletes", err)

		return err
	}
	// helper.Log("debug")("log connection func LikeDeletes", keys)

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
