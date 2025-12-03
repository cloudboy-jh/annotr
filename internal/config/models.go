package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Model struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ContextWindow int    `json:"contextWindow,omitempty"`
}

type Provider struct {
	APIKeyPattern string  `json:"apiKeyPattern,omitempty"`
	Endpoint      string  `json:"endpoint"`
	RequiresKey   *bool   `json:"requiresKey,omitempty"`
	Models        []Model `json:"models,omitempty"`
}

type ModelsManifest struct {
	Version     string              `json:"version"`
	LastUpdated string              `json:"lastUpdated"`
	Providers   map[string]Provider `json:"providers"`
}

func DefaultModelsManifest() *ModelsManifest {
	return &ModelsManifest{
		Version:     "1.0.0",
		LastUpdated: "2024-12-02",
		Providers: map[string]Provider{
			"anthropic": {
				APIKeyPattern: "sk-ant-",
				Endpoint:      "https://api.anthropic.com/v1/messages",
				Models: []Model{
					{ID: "claude-sonnet-4-20250514", Name: "Claude Sonnet 4", ContextWindow: 200000},
					{ID: "claude-3-5-sonnet-20241022", Name: "Claude 3.5 Sonnet", ContextWindow: 200000},
				},
			},
			"openai": {
				APIKeyPattern: "sk-",
				Endpoint:      "https://api.openai.com/v1/chat/completions",
				Models: []Model{
					{ID: "gpt-4o", Name: "GPT-4o", ContextWindow: 128000},
					{ID: "gpt-4o-mini", Name: "GPT-4o Mini", ContextWindow: 128000},
				},
			},
			"groq": {
				APIKeyPattern: "gsk_",
				Endpoint:      "https://api.groq.com/openai/v1/chat/completions",
				Models: []Model{
					{ID: "llama-3.3-70b-versatile", Name: "Llama 3.3 70B", ContextWindow: 32768},
				},
			},
			"ollama": {
				Endpoint:    "http://localhost:11434",
				RequiresKey: boolPtr(false),
			},
		},
	}
}

func boolPtr(b bool) *bool {
	return &b
}

func ModelsManifestPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "models.json"), nil
}

func LoadModelsManifest() (*ModelsManifest, error) {
	path, err := ModelsManifestPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultModelsManifest(), nil
		}
		return nil, err
	}

	var manifest ModelsManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}

	return &manifest, nil
}

func (m *ModelsManifest) Save() error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := ModelsManifestPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
