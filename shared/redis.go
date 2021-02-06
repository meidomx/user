package shared

import "github.com/go-redis/redis/v8"

var RedisClient *redis.Client

func InitRedis(addr string, pwd string) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})

	RedisClient = rdb
}
