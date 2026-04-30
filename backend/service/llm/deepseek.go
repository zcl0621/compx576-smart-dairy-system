package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	summaryMaxLen = 200
	noteMaxLen    = 500
)

// LLMOutput holds summary and note parsed from a chat completion response.
type LLMOutput struct {
	Summary string `json:"summary"`
	Note    string `json:"note"`
}

// Client calls a DeepSeek (or OpenAI-compatible) chat completion endpoint.
// Returns errors directly; caller handles retries.
// HTTP must be set by the caller; use &http.Client{Timeout: Timeout} for production.
type Client struct {
	BaseURL string
	APIKey  string
	Model   string
	Timeout time.Duration
	HTTP    *http.Client
}

type chatRequest struct {
	Model          string    `json:"model"`
	Messages       []chatMsg `json:"messages"`
	ResponseFormat respFmt   `json:"response_format"`
}

type chatMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type respFmt struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Generate sends system + user messages and returns parsed output.
func (c *Client) Generate(ctx context.Context, system, user string) (LLMOutput, error) {
	body, err := json.Marshal(chatRequest{
		Model: c.Model,
		Messages: []chatMsg{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		ResponseFormat: respFmt{Type: "json_object"},
	})
	if err != nil {
		return LLMOutput{}, fmt.Errorf("deepseek marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return LLMOutput{}, fmt.Errorf("deepseek build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return LLMOutput{}, fmt.Errorf("deepseek http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return LLMOutput{}, fmt.Errorf("deepseek status %d: %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var parsed chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return LLMOutput{}, fmt.Errorf("deepseek decode envelope: %w", err)
	}
	if len(parsed.Choices) == 0 || parsed.Choices[0].Message.Content == "" {
		return LLMOutput{}, errors.New("deepseek empty content")
	}

	var out LLMOutput
	if err := json.Unmarshal([]byte(parsed.Choices[0].Message.Content), &out); err != nil {
		return LLMOutput{}, fmt.Errorf("deepseek parse content: %w", err)
	}
	if strings.TrimSpace(out.Summary) == "" {
		return LLMOutput{}, errors.New("deepseek missing summary")
	}
	if strings.TrimSpace(out.Note) == "" {
		return LLMOutput{}, errors.New("deepseek missing note")
	}
	out.Summary = truncate(out.Summary, summaryMaxLen)
	out.Note = truncate(out.Note, noteMaxLen)
	return out, nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}
