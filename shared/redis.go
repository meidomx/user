package shared

import "github.com/go-redis/redis/v8"

var RedisClient *redis.Client

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		//FIXME should use configuration
		Addr:     "192.168.31.231:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	RedisClient = rdb
}
