package telegram

import (
	"context"
	"fmt"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

type Provider struct {
	client *Client
	chatID string
}

func NewProvider(client *Client, chatID string) *Provider {
	return &Provider{client: client, chatID: chatID}
}

func (p *Provider) Name() string {
	return "telegram"
}

func (p *Provider) Send(ctx context.Context, n notification.Notification) error {
	if p.client == nil {
		return fmt.Errorf("telegram client not configured")
	}
	if p.chatID == "" {
		return fmt.Errorf("telegram chat id not configured")
	}

	text := FormatNotification(n)
	return p.client.SendMessage(ctx, p.chatID, text)
}
