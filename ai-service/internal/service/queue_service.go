package service

import (
	"ai-service/internal/model"
	openai "ai-service/internal/openAi"
	"encoding/json"
	"fmt"
)

func ProcessMessage(body []byte) {
	var req model.TestRequest
	if err := json.Unmarshal(body, &req); err != nil {
		fmt.Println("Xabar parse qilinishida xatolik:", err)
		return
	}

	// AI’ga yuborish uchun call
	questions, err := openai.ProcessPrompt(req.Prompt, req.Count)
	if err != nil {
		fmt.Println("AI bilan ishlashda xatolik:", err)
		return
	}

	// Test uchun ekranga chiqarish
	fmt.Println("\n===== QUESTIONS =====")
	fmt.Printf("TelegramID: %d\n", req.TelegramID)
	for i, q := range questions {
		fmt.Printf("%d) %s\n", i+1, q.Question)
		fmt.Printf("Options: %v\n", q.Options)
		fmt.Printf("Correct: %s\n\n", q.Correct)
	}
}
