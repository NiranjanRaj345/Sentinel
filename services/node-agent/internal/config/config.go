package config

import (
	"fmt"
	"os"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/alert"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Agent   AgentConfig   `yaml:"agent"`
	Logging LoggingConfig `yaml:"logging"`
	Metrics MetricsConfig `yaml:"metrics"`
	Storage StorageConfig `yaml:"storage"`
	Alerts  AlertsConfig  `yaml:"alerts"`
	Auth    AuthConfig    `yaml:"auth"`
}

type AuthConfig struct {
	Tokens []AuthTokenConfig `yaml:"tokens"`
}

type AuthTokenConfig struct {
	Name        string     `yaml:"name"`
	Value       string     `yaml:"value"`
	Permissions []string   `yaml:"permissions"`
	ExpiresAt   *time.Time `yaml:"expiresAt"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type AgentConfig struct {
	Name string `yaml:"name"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type MetricsConfig struct {
	Interval string `yaml:"interval"`
}

type StorageConfig struct {
	Path string `yaml:"path"`
}

type AlertsConfig struct {
	Rules []alert.Rule `yaml:"rules"`
}

func Default() Config {
	return Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Agent: AgentConfig{
			Name: "sentinel-node",
		},
		Logging: LoggingConfig{
			Level: "info",
		},
		Metrics: MetricsConfig{
			Interval: "5s",
		},
		Storage: StorageConfig{
			Path: "sentinel.db",
		},
		Alerts: AlertsConfig{
			Rules: []alert.Rule{
				{
					ID:        "cpu-warning",
					Name:      "CPU Usage",
					Metric:    "cpu.usage",
					Operator:  alert.GreaterThan,
					Threshold: 80.0,
					Severity:  alert.SeverityWarning,
					Enabled:   true,
				},
				{
					ID:        "cpu-critical",
					Name:      "CPU Usage",
					Metric:    "cpu.usage",
					Operator:  alert.GreaterThanOrEqual,
					Threshold: 90.0,
					Severity:  alert.SeverityCritical,
					Enabled:   true,
				},
				{
					ID:        "memory-warning",
					Name:      "Memory Usage",
					Metric:    "memory.used_percent",
					Operator:  alert.GreaterThan,
					Threshold: 85.0,
					Severity:  alert.SeverityWarning,
					Enabled:   true,
				},
				{
					ID:        "disk-critical",
					Name:      "Disk Usage",
					Metric:    "disk.used_percent",
					Operator:  alert.GreaterThan,
					Threshold: 95.0,
					Severity:  alert.SeverityCritical,
					Enabled:   true,
				},
			},
		},
	}
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (c Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}

	if err := c.Agent.Validate(); err != nil {
		return err
	}

	if err := c.Logging.Validate(); err != nil {
		return err
	}

	if err := c.Metrics.Validate(); err != nil {
		return err
	}

	if err := c.Storage.Validate(); err != nil {
		return err
	}

	if err := c.Alerts.Validate(); err != nil {
		return err
	}

	if err := c.Auth.Validate(); err != nil {
		return err
	}

	return nil
}

func (c ServerConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("server.host cannot be empty")
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

	return nil
}

func (c AgentConfig) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("agent.name cannot be empty")
	}

	return nil
}

func (c LoggingConfig) Validate() error {
	if _, err := logger.ParseLevel(c.Level); err != nil {
		return fmt.Errorf("logging.level: %w", err)
	}

	return nil
}

func (c MetricsConfig) Validate() error {
	if c.Interval == "" {
		return fmt.Errorf("metrics.interval cannot be empty")
	}

	if _, err := time.ParseDuration(c.Interval); err != nil {
		return fmt.Errorf("metrics.interval must be a valid duration: %w", err)
	}

	return nil
}

func (c StorageConfig) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("storage.path cannot be empty")
	}

	return nil
}

func (c AlertsConfig) Validate() error {
	for _, rule := range c.Rules {
		if err := rule.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c AuthConfig) Validate() error {
	for _, token := range c.Tokens {
		if token.Name == "" {
			return fmt.Errorf("auth token name cannot be empty")
		}
		if token.Value == "" {
			return fmt.Errorf("auth token value cannot be empty for token %q", token.Name)
		}
		for _, perm := range token.Permissions {
			switch perm {
			case "read", "operate", "configure":
			default:
				return fmt.Errorf("invalid permission %q for token %q", perm, token.Name)
			}
		}
	}
	return nil
}

func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
