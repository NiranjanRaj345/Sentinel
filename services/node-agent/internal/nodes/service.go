package nodes

import (
	"context"
	"fmt"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Service struct {
	repo    Repository
	publish func(context.Context, events.Event) error
	log     *logger.Logger
}

func NewService(repo Repository, publish func(context.Context, events.Event) error, log *logger.Logger) *Service {
	if log == nil {
		log = logger.New(logger.Info, nil)
	}
	return &Service{repo: repo, publish: publish, log: log}
}

func (s *Service) Register(ctx context.Context, node Node) (Node, error) {
	if s == nil || s.repo == nil {
		return Node{}, fmt.Errorf("node service not initialized")
	}

	now := time.Now().UTC()
	node.CreatedAt = now
	node.LastSeen = now
	node.Status = StatusOnline

	if err := s.repo.Save(ctx, node); err != nil {
		return Node{}, err
	}

	s.publishEvent(ctx, events.EventTypeSystem, events.SeverityInfo, events.SourceSystem,
		"Node registered",
		fmt.Sprintf("Node %s registered as %s", node.Name, node.ID),
		map[string]interface{}{"nodeId": node.ID, "nodeName": node.Name},
	)

	return node, nil
}

func (s *Service) Heartbeat(ctx context.Context, req HeartbeatRequest) (Node, error) {
	if s == nil || s.repo == nil {
		return Node{}, fmt.Errorf("node service not initialized")
	}

	node, err := s.repo.Get(ctx, req.ID)
	if err != nil {
		return Node{}, err
	}

	prevStatus := node.Status
	now := time.Now().UTC()
	node.LastSeen = now
	node.Status = StatusOnline

	if err := s.repo.UpdateLastSeen(ctx, node.ID, now); err != nil {
		return Node{}, err
	}

	if prevStatus != StatusOnline {
		if err := s.repo.UpdateStatus(ctx, node.ID, StatusOnline); err != nil {
			return Node{}, err
		}
		node.Status = StatusOnline

		s.publishEvent(ctx, events.EventTypeSystem, events.SeverityInfo, events.SourceSystem,
			"Node online",
			fmt.Sprintf("Node %s is back online", node.Name),
			map[string]interface{}{"nodeId": node.ID, "nodeName": node.Name},
		)
	}

	return node, nil
}

func (s *Service) CheckOfflineNodes(ctx context.Context, timeout time.Duration) {
	if s == nil || s.repo == nil {
		return
	}

	nodes, err := s.repo.List(ctx)
	if err != nil {
		s.log.Error("failed to list nodes for offline check: %v", err)
		return
	}

	now := time.Now().UTC()
	for _, node := range nodes {
		if node.Status == StatusOnline && now.Sub(node.LastSeen) > timeout {
			if err := s.repo.UpdateStatus(ctx, node.ID, StatusOffline); err != nil {
				s.log.Error("failed to mark node offline: %v", err)
				continue
			}

			s.publishEvent(ctx, events.EventTypeSystem, events.SeverityWarning, events.SourceSystem,
				"Node offline",
				fmt.Sprintf("Node %s has not sent a heartbeat in %s", node.Name, now.Sub(node.LastSeen).Round(time.Second)),
				map[string]interface{}{"nodeId": node.ID, "nodeName": node.Name},
			)
		}
	}
}

func (s *Service) List(ctx context.Context) ([]Node, error) {
	if s == nil || s.repo == nil {
		return nil, fmt.Errorf("node service not initialized")
	}
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id string) (Node, error) {
	if s == nil || s.repo == nil {
		return Node{}, fmt.Errorf("node service not initialized")
	}
	return s.repo.Get(ctx, id)
}

func (s *Service) Remove(ctx context.Context, id string) error {
	if s == nil || s.repo == nil {
		return fmt.Errorf("node service not initialized")
	}
	return s.repo.Remove(ctx, id)
}

func (s *Service) Close() error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Close()
}

func (s *Service) publishEvent(ctx context.Context, eventType events.EventType, severity events.Severity, source events.Source, title, message string, metadata map[string]interface{}) {
	if s.publish == nil {
		return
	}
	_ = s.publish(ctx, events.Event{
		Type:     eventType,
		Severity: severity,
		Source:   source,
		Title:    title,
		Message:  message,
		Metadata: metadata,
		CreatedAt: time.Now().UTC(),
	})
}
