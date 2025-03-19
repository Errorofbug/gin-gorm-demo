package gredis

import (
	"encoding/json"
	"fmt"
	"gin-gorm-demo/conf"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RedisConn *redis.Pool

// InitRedis 初始化redis连接池
func InitRedis() {
	RedisConn = &redis.Pool{
		MaxIdle:     conf.Settings.Redis.MaxIdle,
		MaxActive:   conf.Settings.Redis.MaxActive,
		IdleTimeout: conf.Settings.Redis.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Settings.Redis.Host)
			if err != nil {
				log.Fatalf(err.Error())
			}
			if conf.Settings.Redis.Password != "" {
				if _, err := c.Do("AUTH", conf.Settings.Redis.Password); err != nil {
					c.Close()
					log.Fatalf(err.Error())
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			log.Fatalf(err.Error())
			return err
		},
	}
}

// GenerateKey 生成带有前缀的key
func GenerateKey(prefix string, key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}

// Set 设置string
func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	var value []byte
	var err error

	switch v := data.(type) {
	case string:
		value = []byte(v)
	default:
		value, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}

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

// Exists 检查键值是否存在
func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

// TTL 查看键值TTL
func TTL(key string) (int, error) {
	conn := RedisConn.Get()
	defer conn.Close()
	ttl, err := redis.Int(conn.Do("TTL", key))
	if err != nil {
		return 0, err
	}

	return ttl, nil
}

// Expire 设置键值超时时间
func Expire(key string, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()
	_, err := conn.Do("EXPIRE", key, time)
	if err != nil {
		return err
	}
	return nil
}

// Get 获取string
func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Delete 删除键值
func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

// LikeDeletes 删除指定匹配模式的键值
func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	var cursor int = 0
	for {
		res, err := redis.Values(conn.Do("SCAN", cursor, "MATCH", "*"+key+"*"))
		if err != nil {
			return err
		}
		cursor, _ = redis.Int(res[0], nil)
		keys, _ := redis.Strings(res[1], nil)
		for _, key := range keys {
			_, err = Delete(key)
			if err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}

	return nil
}
