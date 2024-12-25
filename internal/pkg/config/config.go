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
	Host               string
	Port               string
	DB                 string
	User               string
	Password           string
	SSLMode            string
	ConnMaxLifetime    time.Duration
	MaxIdleConnections int
	MaxOpenConnections int
	PingTimeout        time.Duration
}

type HTMLTemplateConfig struct {
	TemplateDir string
	LayoutFile  string
	PagesDir    string
	PartialsDir string
}

func Load() *Config {
	return &Config{
		Server: HTTPServerConfig{
			Addr:            env.Get("SERVER_HOST", "0.0.0.0"),
			Port:            env.Get("SERVER_PORT", "8888"),
			ShutdownTimeout: 10 * time.Second,
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     60 * time.Second,
		},
		DB: DBConfig{
			Driver:             "pgx",
			Host:               env.MustGet("DB_HOST"),
			Port:               env.MustGet("DB_PORT"),
			DB:                 env.MustGet("DB_NAME"),
			User:               env.MustGet("DB_USER"),
			Password:           env.MustGet("DB_PASSWORD"),
			SSLMode:            env.MustGet("DB_SSLMODE"),
			ConnMaxLifetime:    0,
			MaxIdleConnections: 50,
			MaxOpenConnections: 50,
			PingTimeout:        time.Duration(env.GetInt("DB_PING_TIMEOUT", 5)) * time.Second,
		},
		HTML: HTMLTemplateConfig{
			TemplateDir: "templates",
			LayoutFile:  "layout.html",
			PagesDir:    "pages",
			PartialsDir: "partials",
		},
	}
}
