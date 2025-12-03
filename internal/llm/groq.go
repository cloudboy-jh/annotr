package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GroqClient struct {
	apiKey   string
	model    string
	endpoint string
	client   *http.Client
}

func NewGroqClient(apiKey, model string) *GroqClient {
	return &GroqClient{
		apiKey:   apiKey,
		model:    model,
		endpoint: "https://api.groq.com/openai/v1/chat/completions",
		client:   &http.Client{},
	}
}

func (c *GroqClient) Provider() string {
	return "groq"
}

func (c *GroqClient) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	model := req.Model
	if model == "" {
		model = c.model
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1024
	}

	groqReq := openaiRequest{
		Model:       model,
		Messages:    req.Messages,
		MaxTokens:   maxTokens,
		Temperature: req.Temperature,
	}

	body, err := json.Marshal(groqReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("groq returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var groqResp openaiResponse
	if err := json.NewDecoder(resp.Body).Decode(&groqResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	content := ""
	if len(groqResp.Choices) > 0 {
		content = groqResp.Choices[0].Message.Content
	}

	return &CompletionResponse{
		Content: content,
		Model:   groqResp.Model,
		Usage: Usage{
			PromptTokens:     groqResp.Usage.PromptTokens,
			CompletionTokens: groqResp.Usage.CompletionTokens,
			TotalTokens:      groqResp.Usage.TotalTokens,
		},
	}, nil
}
