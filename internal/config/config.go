package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server     ServerConfig      `mapstructure:"server"`
	Sequencers []SequencerConfig `mapstructure:"sequencers"`
	Log        LogConfig         `mapstructure:"log"`
}

type ServerConfig struct {
	ListenAddr     string        `mapstructure:"listen_addr"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
	MaxMessageSize int           `mapstructure:"max_message_size"`
}

type SequencerConfig struct {
	ChainID  string `mapstructure:"chain_id"`
	Endpoint string `mapstructure:"endpoint"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// Set defaults
	viper.SetDefault("server.listen_addr", ":8080")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.max_message_size", 10*1024*1024) // 10MB
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if len(c.Sequencers) < 2 {
		return fmt.Errorf("at least 2 sequencers required")
	}
	return nil
}
