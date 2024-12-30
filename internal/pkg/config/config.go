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
			ShutdownTimeout: time.Duration(env.GetInt("SERVER_SHUTDOWN_TIMEOUT", 10)) * time.Second,
			ReadTimeout:     time.Duration(env.GetInt("SERVER_READ_TIMEOUT", 10)) * time.Second,
			WriteTimeout:    time.Duration(env.GetInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second,
			IdleTimeout:     time.Duration(env.GetInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second,
		},
		DB: DBConfig{
			Driver:             "pgx",
			Host:               env.MustGet("DB_HOST"),
			Port:               env.MustGet("DB_PORT"),
			DB:                 env.MustGet("DB_NAME"),
			User:               env.MustGet("DB_USER"),
			Password:           env.MustGet("DB_PASSWORD"),
			SSLMode:            env.MustGet("DB_SSLMODE"),
			ConnMaxLifetime:    time.Duration(env.GetInt("DB_CONN_MAX_LIFETIME", 0)) * time.Second,
			MaxIdleConnections: env.GetInt("DB_MAX_IDLE_CONNS", 50),
			MaxOpenConnections: env.GetInt("DB_MAX_OPEN_CONNS", 50),
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
