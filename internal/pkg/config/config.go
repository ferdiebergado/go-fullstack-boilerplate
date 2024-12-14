package config

import (
	"time"

	"github.com/ferdiebergado/gopherkit/env"
)

type Config struct {
	Server HTTPServerConfig
	DB     DBConfig
	HTML   HTMLTemplateConfig
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

type HTMLTemplateConfig struct {
	TemplateDir      string
	LayoutFile       string
	PartialTemplates string
}

func Load() *Config {
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
			PingTimeout:        5 * time.Second,
		},
		HTML: HTMLTemplateConfig{
			TemplateDir:      "templates",
			LayoutFile:       "layout.html",
			PartialTemplates: "partials/*.html",
		},
	}
}
