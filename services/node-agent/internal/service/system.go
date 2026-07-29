package service

import "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/system"

// SystemService provides access to system-related information.
type SystemService struct{}

// NewSystemService creates a new SystemService.
func NewSystemService() *SystemService {
	return &SystemService{}
}

// GetInfo returns information about the current system.
func (s *SystemService) GetInfo() system.Info {
	return system.GetInfo()
}
