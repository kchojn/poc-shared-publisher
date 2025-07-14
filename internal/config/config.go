package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Network NetworkConfig  `mapstructure:"network"`
	Rollups []RollupConfig `mapstructure:"rollups"`
	Log     LogConfig      `mapstructure:"log"`
}

type NetworkConfig struct {
	ListenAddr     string        `mapstructure:"listen_addr"`
	MaxMessageSize int           `mapstructure:"max_message_size"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout"`
}

type RollupConfig struct {
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
	viper.SetDefault("network.listen_addr", ":8080")
	viper.SetDefault("network.max_message_size", 10*1024*1024) // 10MB
	viper.SetDefault("network.read_timeout", 30*time.Second)
	viper.SetDefault("network.write_timeout", 30*time.Second)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
