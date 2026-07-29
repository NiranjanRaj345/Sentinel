package config

import "fmt"

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Agent   AgentConfig   `yaml:"agent"`
	Logging LoggingConfig `yaml:"logging"`
	Metrics MetricsConfig `yaml:"metrics"`
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
	}
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
