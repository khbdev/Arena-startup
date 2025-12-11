package model

type TestRequest struct {
	TelegramID int64  
	Prompt     string 
	Count      int   
}

type TestData struct {
	TeacherTelegramID int64       `json:"teacher_telegram_id"`
	TestID            string      `json:"test_id"`
	Questions         interface{} `json:"questions"` 


}

type NotificationEvent struct {
	TelegramID int64  `json:"telegram_id"`
	TestID     string `json:"test_id"`
}
