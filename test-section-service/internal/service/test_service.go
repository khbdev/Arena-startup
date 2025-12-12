package service

import (
	"encoding/json"
	"errors"
	"test-section-serve/internal/model"
	repository "test-section-serve/internal/repostroy"
)

type TestService struct{}

func NewTestService() *TestService {
	return &TestService{}
}

func (s *TestService) GetUserTestResult(telegramID int64, testID string) (string, error) {

	
	redisKey := testID 
	jsonStr, err := repository.GetJSONByKey(redisKey)
	if err != nil {
		return "", err
	}
	if jsonStr == "" {
		return "", errors.New("data not found in redis")
	}


	var data model.TestData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", err
	}


	if data.TeacherTelegramID != telegramID {
	
		for i := range data.Questions {
			data.Questions[i].Correct = ""
		}
	}

	filteredJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(filteredJSON), nil
}
