package config

import (
	"time"

	"github.com/ferdiebergado/gopherkit/env"
)

type Config struct {
	Server HTTPServerConfig
	DB     DBConfig
}

type HTTPServerConfig struct {
	Addr            string
	Port            string
	ShutdownTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

type DBConfig struct {
	Driver             string
	DSN                string
	ConnMaxLifetime    time.Duration
	MaxIdleConnections int
	MaxOpenConnections int
	PingTimeout        time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Server: HTTPServerConfig{
			Port:            env.Get("PORT", "8888"),
			ShutdownTimeout: 10 * time.Second,
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     60 * time.Second,
		},
		DB: DBConfig{
			Driver:             "pgx",
			DSN:                env.MustGet("DATABASE_URL"),
			ConnMaxLifetime:    0,
			MaxIdleConnections: 50,
			MaxOpenConnections: 50,
			PingTimeout:        1 * time.Second,
		},
	}
}
