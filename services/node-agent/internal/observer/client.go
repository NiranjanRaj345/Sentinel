package observer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Client struct {
	baseURL string
	client  *http.Client
	log     *logger.Logger
}

func NewClient(baseURL string, log *logger.Logger) *Client {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
		log:     log,
	}
}

func (c *Client) Status(ctx context.Context) (StatusResponse, error) {
	var status StatusResponse
	if err := c.request(ctx, "GET", "/status", nil, &status); err != nil {
		return status, err
	}
	return status, nil
}

func (c *Client) Environment(ctx context.Context) (EnvironmentResponse, error) {
	var env EnvironmentResponse
	if err := c.request(ctx, "GET", "/environment", nil, &env); err != nil {
		return env, err
	}
	return env, nil
}

func (c *Client) Health(ctx context.Context) error {
	return c.request(ctx, "GET", "/health", nil, nil)
}

func (c *Client) request(ctx context.Context, method, path string, body, out interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal observer request: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create observer request: %w", err)
	}
	if len(payload) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("observer request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("observer request failed: status=%d path=%s", resp.StatusCode, path)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decode observer response: %w", err)
		}
	}

	return nil
}
