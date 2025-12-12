package repository

import (
	"test-section-serve/internal/config"

	"github.com/redis/go-redis/v9"
)


func GetJSONByKey(key string) (string, error) {
	val, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil 
		}
		return "", err 
	}
	return val, nil
}
