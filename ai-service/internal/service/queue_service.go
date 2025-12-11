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
		fmt.Println(" Xabar parse qilinishida xatolik:", err)
		return
	}


	questions, err := openai.ProcessPrompt(req.Prompt, req.Count)
	if err != nil {
		fmt.Println(" AI bilan ishlashda xatolik:", err)
		return
	}


	err = TestService(req.TelegramID, questions)
	if err != nil {
		fmt.Println(" TestService xatosi:", err)
		return
	}

	fmt.Println(" Test Redis ga saqlandi!")
}
