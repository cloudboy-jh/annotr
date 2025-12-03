package config

import (
	"encoding/json"
	"net/http"
	"time"
)

type OllamaModel struct {
	Name       string `json:"name"`
	ModifiedAt string `json:"modified_at"`
	Size       int64  `json:"size"`
}

type OllamaTagsResponse struct {
	Models []OllamaModel `json:"models"`
}

func DetectOllama() (bool, []OllamaModel, error) {
	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false, nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, nil
	}

	var tagsResp OllamaTagsResponse
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return true, nil, nil
	}

	return true, tagsResp.Models, nil
}

func IsOllamaRunning() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
