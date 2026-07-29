package stream

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	"github.com/gorilla/websocket"
)

func TestHub_Register_AddsClient(t *testing.T) {
	log := logger.New(logger.Info, io.Discard)
	hub := New(log)
	hub.Start()
	defer hub.Stop()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		hub.Register(conn)
	}))
	defer server.Close()

	u := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)
}

func TestHub_Broadcast_NoClients(t *testing.T) {
	log := logger.New(logger.Info, io.Discard)
	hub := New(log)
	hub.Start()
	defer hub.Stop()

	hub.Broadcast(metrics.Info{Metadata: metrics.Metadata{Timestamp: time.Now().UTC()}})
}

func TestHub_Broadcast_SendsToRegisteredClient(t *testing.T) {
	log := logger.New(logger.Info, io.Discard)
	hub := New(log)
	hub.Start()
	defer hub.Stop()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		hub.Register(conn)
	}))
	defer server.Close()

	u := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	info := metrics.Info{
		Metadata: metrics.Metadata{
			Timestamp: time.Now().UTC(),
			Agent: metrics.AgentInfo{
				Name:    "test",
				Version: "1.0",
			},
		},
	}

	hub.Broadcast(info)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	var received metrics.Info
	if err := conn.ReadJSON(&received); err != nil {
		t.Fatalf("read json: %v", err)
	}

	if received.Metadata.Agent.Name != "test" {
		t.Fatalf("expected agent name test, got %s", received.Metadata.Agent.Name)
	}
}

func TestHub_Unregister_RemovesClient(t *testing.T) {
	log := logger.New(logger.Info, io.Discard)
	hub := New(log)
	hub.Start()
	defer hub.Stop()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		hub.Register(conn)
	}))
	defer server.Close()

	u := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	conn.Close()
}
