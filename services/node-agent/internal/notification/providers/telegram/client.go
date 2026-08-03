package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

const telegramAPIBase = "https://api.telegram.org/bot%s/%s"

type Client struct {
	token string
	http  *http.Client
	log   *logger.Logger
}

func NewClient(token string, log *logger.Logger) *Client {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Client{
		token: token,
		http:  &http.Client{Timeout: 10 * time.Second},
		log:   log,
	}
}

func (c *Client) SendMessage(ctx context.Context, chatID, text string) error {
	if c.token == "" {
		return fmt.Errorf("telegram bot token is empty")
	}
	if chatID == "" {
		return fmt.Errorf("telegram chat id is empty")
	}

	reqBody := SendMessageRequest{
		ChatID:    chatID,
		Text:      text,
		ParseMode: "Markdown",
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal telegram request: %w", err)
	}

	url := fmt.Sprintf(telegramAPIBase, c.token, "sendMessage")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("telegram request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram request failed: status=%d", resp.StatusCode)
	}

	var apiResp SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return fmt.Errorf("decode telegram response: %w", err)
	}

	if !apiResp.OK {
		return fmt.Errorf("telegram api error: %s", apiResp.Description)
	}

	return nil
}
