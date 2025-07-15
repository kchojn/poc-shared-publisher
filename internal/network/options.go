package network

import "time"

// ServerOption configures a server.
type ServerOption func(*ServerConfig)

// WithMaxConnections sets the maximum number of connections.
func WithMaxConnections(max int) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.MaxConnections = max
	}
}

// WithTimeouts sets read/write timeouts.
func WithTimeouts(read, write time.Duration) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.ReadTimeout = read
		cfg.WriteTimeout = write
	}
}

// ClientOption configures a client.
type ClientOption func(*ClientConfig)

// WithReconnectDelay sets the reconnection delay.
func WithReconnectDelay(delay time.Duration) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.ReconnectDelay = delay
	}
}

// WithConnectTimeout sets the connection timeout.
func WithConnectTimeout(timeout time.Duration) ClientOption {
	return func(cfg *ClientConfig) {
		cfg.ConnectTimeout = timeout
	}
}
