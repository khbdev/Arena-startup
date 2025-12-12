package model


type TestData struct {
    TeacherTelegramID int64 `json:"teacher_telegram_id"`
    TestID            string `json:"test_id"`
    Questions         []Question `json:"questions"`
}

type Question struct {
    ID       string   `json:"id"`
    Question string   `json:"question"`
    Options  []Option `json:"options"`
    Correct  string   `json:"correct"`
}

type Option struct {
    ID   string `json:"id"`
    Text string `json:"text"`
}
