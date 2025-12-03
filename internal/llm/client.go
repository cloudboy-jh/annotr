package llm

import (
	"context"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type CompletionResponse struct {
	Content string
	Model   string
	Usage   Usage
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Client interface {
	Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	Provider() string
}

func NewClient(provider, apiKey, model string) Client {
	switch provider {
	case "ollama":
		return NewOllamaClient(model)
	case "anthropic":
		return NewAnthropicClient(apiKey, model)
	case "openai":
		return NewOpenAIClient(apiKey, model)
	case "groq":
		return NewGroqClient(apiKey, model)
	default:
		return NewOllamaClient(model)
	}
}
