package app

import (
	m "temp/internal/app/models"
	"temp/pkg/env"
	"time"
)

func loadConfigFromEnv() (*m.Config, error) {
	version := env.String("VERSION")
	environment := env.String("ENV")

	conn := env.String("DATABASE_CONN_STRING")

	addr := env.String("SERVER_ADDR")
	readTimeout := env.Int("SERVER_READ_TIMEOUT")
	writeTimeout := env.Int("SERVER_WRITE_TIMEOUT")
	idleTimeout := env.Int("SERVER_IDLE_TIMEOUT")

	accessTokenTTL := env.Int("AUTH_ACCESS_TOKEN_TTL")
	refreshTokenTTL := env.Int("AUTH_REFRESH_TOKEN_TTL")

	cfg := m.Config{
		Version: version,
		Env:     environment,
		Database: m.DatabaseConfig{
			Conn: conn,
		},
		Server: m.ServerConfig{
			Addr:         addr,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		},
		Auth: m.AuthConfig{
			AccessTokenTTL:  time.Duration(accessTokenTTL) * time.Minute,
			RefreshTokenTTL: time.Duration(refreshTokenTTL) * time.Minute,
		},
	}

	return &cfg, nil
}
