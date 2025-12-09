package pkg
import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func getUserCacheKey(telegramID int64) string {
	return "user:" + fmt.Sprint(telegramID)
}

func getCacheTTL(envKey string, defaultSec int) time.Duration {
	val := os.Getenv(envKey)
	if val == "" {
		return time.Duration(defaultSec) * time.Second
	}
	sec, err := strconv.Atoi(val)
	if err != nil {
		return time.Duration(defaultSec) * time.Second
	}
	return time.Duration(sec) * time.Second
}
