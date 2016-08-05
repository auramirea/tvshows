package client

import (
	"fmt"
	"gopkg.in/redis.v4"
	"sync"
)

type RedisClient struct {
	client *redis.Client
}

var instance *RedisClient
var once sync.Once

func GetRedisClient() *RedisClient {
	once.Do(func() {
		instance = &RedisClient{client: newClient()}
	})
	return instance
}

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.10.13:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	client.FlushAll()
	fmt.Println(pong, err)
	return client
}

func (c *RedisClient) SetKey(key string, value string) {
	err := c.client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func (c *RedisClient) GetKey(key string) string {
	val, err := c.client.Get(key).Result()
	if err == redis.Nil {
		fmt.Println(key, "does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println(key, val)
	}
	return val
}
