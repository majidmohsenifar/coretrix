package platform

import "github.com/go-redis/redis/v8"

func NewRedisClient(c Configs) *redis.Client {
	dsn := c.GetString("redis.dsn")
	password := c.GetString("redis.password")
	db := c.GetInt("redis.db")
	client := redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: password,
		DB:       db,
	})
	return client
}
