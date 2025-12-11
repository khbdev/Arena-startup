package service

import (
	"ai-service/internal/config"
	"ai-service/internal/event"
	"ai-service/internal/model"

	"ai-service/internal/util"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

func TestService(teacherTelegramID int64, questions interface{}) error {

	testID := util.GenerateTestID()

	
	data := model.TestData{
		TeacherTelegramID: teacherTelegramID,
		TestID:            testID,
		Questions:         questions,
	}

	
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json marshal xatosi: %v", err)
	}

	
	ttlHoursStr := os.Getenv("REDIS_TTL")
	ttlHours, err := strconv.Atoi(ttlHoursStr)
	if err != nil || ttlHours <= 0 {
		ttlHours = 2
	}

	err = config.RedisClient.Set(
		config.Ctx,
		testID,
		jsonData,
		time.Duration(ttlHours)*time.Hour,
	).Err()
	if err != nil {
		return fmt.Errorf("redisga yozishda xatolik: %v", err)
	}

	fmt.Println(" Redis saqlandi, key:", testID)

r := config.NewRabbitMQ() 
	if err := event.PublishNotification(r.Channel, teacherTelegramID, testID); err != nil {
		return fmt.Errorf("event yuborishda xatolik: %v", err)
	}
	
	return nil
}
