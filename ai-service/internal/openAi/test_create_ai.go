package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Question struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Correct  string   `json:"correct"`
}

func ProcessPrompt(prompt string, count int) ([]Question, error) {
	fullPrompt := fmt.Sprintf(`
Generate %d %s question in JSON:
{"question":"","options":["A","B","C","D"],"correct":""}
Return ONLY JSON array.
`, count, prompt)

	body, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-5-nano",
		"input": fullPrompt,
		"store": true,
	})

	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY missing")
	}

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/responses", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var parsed map[string]interface{}
	if err := json.Unmarshal(resBody, &parsed); err != nil {
		return nil, err
	}

	outputArr := parsed["output"].([]interface{})
	content := outputArr[1].(map[string]interface{})["content"].([]interface{})
	text := content[0].(map[string]interface{})["text"].(string)

	var questions []Question
	if err := json.Unmarshal([]byte(text), &questions); err != nil {
		return nil, err
	}

	return questions, nil
}
