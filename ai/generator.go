package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"aiart-bot/config"
)

type GenerateResponse struct {
	Data struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

type ResultResponse struct {
	Data struct {
		Images []string `json:"images"`
	} `json:"data"`
}

func GenerateImage(prompt string) ([]byte, error) {
	url := "https://api.novita.ai/v1/task/create" // 🔁 правильный endpoint
	payload := map[string]interface{}{
		"prompt": prompt,
		"model":  "Anime-v2", // варианты: MeinaMix, Anything-v5, AbyssOrangeMix
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.NovitaAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Generation error: %s", string(respData))
	}

	var genResp GenerateResponse
	json.NewDecoder(resp.Body).Decode(&genResp)
	taskID := genResp.Data.TaskID
	log.Println("TaskID:", taskID)

	// ⏳ Подождём немного, пока картинка генерируется
	time.Sleep(8 * time.Second)

	// 🔁 Получаем результат
	fetchURL := fmt.Sprintf("https://api.novita.ai/v1/task/fetch?task_id=%s", taskID)
	req2, _ := http.NewRequest("GET", fetchURL, nil)
	req2.Header.Set("Authorization", config.NovitaAPIKey)

	resp2, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != 200 {
		respText, _ := io.ReadAll(resp2.Body)
		return nil, fmt.Errorf("Fetch error: %s", string(respText))
	}

	var result ResultResponse
	json.NewDecoder(resp2.Body).Decode(&result)

	if len(result.Data.Images) == 0 {
		return nil, fmt.Errorf("No image returned")
	}

	// Загружаем картинку по URL
	imgResp, err := http.Get(result.Data.Images[0])
	if err != nil {
		return nil, err
	}
	defer imgResp.Body.Close()

	return io.ReadAll(imgResp.Body)
}
