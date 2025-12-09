package pkg

import (
	"encoding/json"

	"log"
	

	"user-service/internal/config"
	domain "user-service/internal/domen"

	"github.com/redis/go-redis/v9"
)

// ReadThroughUser - User uchun read-through pattern
func ReadThroughUser(telegramID int64, loader func() (*domain.User, error)) (*domain.User, error) {
	key := getUserCacheKey(telegramID)

	// 1️⃣ Redis'dan tekshirish
	val, err := config.RedisClient.Get(config.Ctx, key).Result()
	if err == nil {
		var u domain.User
		if err := json.Unmarshal([]byte(val), &u); err == nil {
			log.Println("Cache hit:", key)
			return &u, nil
		}
		log.Println("Cache corrupt, fallback to loader")
	} else if err != nil && err != redis.Nil {
		log.Println("Redis error:", err)
	}

	// 2️⃣ Redisda topilmasa → loader function chaqiriladi
	user, err := loader()
	if err != nil {
		return nil, err
	}

	// 3️⃣ TTLni .env dan o'qish
	ttl := getCacheTTL("USER_CACHE_TTL", 300) // default 300s = 5min

	// Redis'ga set qilish
	bytes, _ := json.Marshal(user)
	err = config.RedisClient.Set(config.Ctx, key, bytes, ttl).Err()
	if err != nil {
		log.Println("Redis set error:", err)
	} else {
		log.Println("Cache set:", key)
	}

	return user, nil
}


