package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig      `mapstructure:"server"`
	Metrics    MetricsConfig     `mapstructure:"metrics"`
	Sequencers []SequencerConfig `mapstructure:"sequencers"`
	Log        LogConfig         `mapstructure:"log"`
}

type ServerConfig struct {
	ListenAddr     string        `mapstructure:"listen_addr"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxMessageSize int           `mapstructure:"max_message_size"`
	MaxConnections int           `mapstructure:"max_connections"`
}

type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

type SequencerConfig struct {
	ChainID  string `mapstructure:"chain_id"`
	Endpoint string `mapstructure:"endpoint"`
	Name     string `mapstructure:"name"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Output string `mapstructure:"output"` // stdout, stderr, file
	File   string `mapstructure:"file"`   // log file path if output=file
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.SetDefault("server.listen_addr", ":8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.max_message_size", 10*1024*1024) // 10MB
	viper.SetDefault("server.max_connections", 100)

	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.port", 8081)
	viper.SetDefault("metrics.path", "/metrics")

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.pretty", false)
	viper.SetDefault("log.output", "stdout")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.MaxMessageSize <= 0 {
		return fmt.Errorf("server.max_message_size must be positive")
	}
	if c.Server.MaxConnections <= 0 {
		return fmt.Errorf("server.max_connections must be positive")
	}

	if c.Metrics.Enabled && c.Metrics.Port <= 0 {
		return fmt.Errorf("metrics.port must be positive when metrics enabled")
	}

	if len(c.Sequencers) < 2 {
		return fmt.Errorf("at least 2 sequencers required for production")
	}

	chainIDs := make(map[string]bool)
	for _, seq := range c.Sequencers {
		if chainIDs[seq.ChainID] {
			return fmt.Errorf("duplicate chain_id: %s", seq.ChainID)
		}
		chainIDs[seq.ChainID] = true
	}

	return nil
}
