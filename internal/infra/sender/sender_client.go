package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LevanPro/insider/internal/service"
)

type requestPayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type responsePayload struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

type Client struct {
	httpClient *http.Client
	url        string
	apiKey     string
}

func NewClient(url, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		url:    url,
		apiKey: apiKey,
	}
}

func (c *Client) Send(ctx context.Context, to, content string) (*service.SendResponse, error) {
	bodyBytes, err := json.Marshal(requestPayload{
		To:      to,
		Content: content,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-ins-auth-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rp responsePayload
	if err := json.NewDecoder(resp.Body).Decode(&rp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &service.SendResponse{
		MessageID: rp.MessageID,
	}, nil
}
