package model

import "time"

type Config struct {
	Name     string
	Env      string
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
	Sentry   SentryConfig
}

type DatabaseConfig struct {
	Conn string
}

type ServerConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type AuthConfig struct {
	AuthJwtSecret   []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type SentryConfig struct {
	Dsn              string
	TracesSampleRate float64
	AttachStackTrace bool
	Debug            bool
}
