package pkg
import (
	"encoding/json"
	
	"log"

	

	"user-service/internal/config"
	domain "user-service/internal/domen"
)

// WriteThroughUser - User uchun write-through pattern
func WriteThroughUser(user *domain.User, writeFunc func(*domain.User) (*domain.User, error)) (*domain.User, error) {
	// 1️⃣ DB yoki usecase ga yozish
	updatedUser, err := writeFunc(user)
	if err != nil {
		return nil, err
	}

	// 2️⃣ Redis'ga set qilish
	key := getUserCacheKey(updatedUser.TelegramID)

	// TTL .env orqali
	ttl := getCacheTTL("USER_CACHE_TTL", 300) // default 5min

	bytes, _ := json.Marshal(updatedUser)
	err = config.RedisClient.Set(config.Ctx, key, bytes, ttl).Err()
	if err != nil {
		log.Println("Redis set error:", err)
	} else {
		log.Println("Cache updated (write-through):", key)
	}

	return updatedUser, nil
}

