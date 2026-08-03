package telegram

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

func TestClient_SendMessage_EmptyToken_ReturnsError(t *testing.T) {
	client := NewClient("", logger.New(logger.Info, nil))
	err := client.SendMessage(context.Background(), "123", "test")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestClient_SendMessage_EmptyChatID_ReturnsError(t *testing.T) {
	client := NewClient("token", logger.New(logger.Info, nil))
	err := client.SendMessage(context.Background(), "", "test")
	if err == nil {
		t.Fatal("expected error for empty chat id")
	}
}

func TestClient_SendMessage_TelegramReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":false,"description":"bad request"}`))
	}))
	defer server.Close()

	client := NewClient("token", logger.New(logger.Info, nil))
	client.http = &http.Client{Timeout: 5 * time.Second}
	err := client.SendMessage(context.Background(), "123", "test")
	if err == nil {
		t.Fatal("expected error from telegram api")
	}
}

func TestProvider_ImplementsInterface(t *testing.T) {
	var _ notification.Provider = NewProvider(NewClient("token", logger.New(logger.Info, nil)), "123")
}

func TestProvider_Send_EmptyChatID_ReturnsError(t *testing.T) {
	client := NewClient("token", logger.New(logger.Info, nil))
	provider := NewProvider(client, "")
	err := provider.Send(context.Background(), notification.Notification{})
	if err == nil {
		t.Fatal("expected error for empty chat id")
	}
}

func TestTemplates_FormatRecovery(t *testing.T) {
	n := notification.NewNotification("id", "Title", "Message", notification.SeverityInfo, notification.SourceRecovery)
	text := FormatRecovery(n)
	if len(text) == 0 {
		t.Fatal("expected non-empty recovery template")
	}
}

func TestTemplates_FormatTestNotification(t *testing.T) {
	text := FormatTestNotification("Desktop")
	if len(text) == 0 {
		t.Fatal("expected non-empty test notification")
	}
	if !strings.Contains(text, "Desktop") {
		t.Fatal("expected test notification to contain node name")
	}
}
