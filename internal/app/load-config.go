package app

import (
	"time"
	model "wildproject/internal/app/domain/models"
	"wildproject/pkg/env"
)

func loadConfigFromEnv() (*model.Config, error) {
	environment := env.String("ENV")
	projName := env.String("NAME")

	conn := env.String("DATABASE_CONN_STRING")

	addr := env.String("SERVER_ADDR")
	readTimeout := env.Int("SERVER_READ_TIMEOUT")
	writeTimeout := env.Int("SERVER_WRITE_TIMEOUT")
	idleTimeout := env.Int("SERVER_IDLE_TIMEOUT")

	authJwtSecret := env.String("AUTH_JWT_SECRET")
	accessTokenTTL := env.Int("AUTH_ACCESS_TOKEN_TTL")
	refreshTokenTTL := env.Int("AUTH_REFRESH_TOKEN_TTL")

	sentryDsn := env.String("SENTRY_DSN")
	sentryTSRate := env.Float64("SENTRY_TRACES_SAMPLE_RATE")
	sentryAttachST := env.Bool("SENTRY_ATTACH_STACKTRACE")
	sentryDebug := env.Bool("SENTRY_DEBUG")

	cfg := model.Config{
		Name: projName,
		Env:  environment,
		Database: model.DatabaseConfig{
			Conn: conn,
		},
		Server: model.ServerConfig{
			Addr:         addr,
			ReadTimeout:  time.Duration(readTimeout) * time.Second,
			WriteTimeout: time.Duration(writeTimeout) * time.Second,
			IdleTimeout:  time.Duration(idleTimeout) * time.Second,
		},
		Auth: model.AuthConfig{
			AuthJwtSecret:   []byte(authJwtSecret),
			AccessTokenTTL:  time.Duration(accessTokenTTL) * time.Minute,
			RefreshTokenTTL: time.Duration(refreshTokenTTL) * time.Minute,
		},
		Sentry: model.SentryConfig{
			Dsn:              sentryDsn,
			TracesSampleRate: sentryTSRate,
			AttachStackTrace: sentryAttachST,
			Debug:            sentryDebug,
		},
	}

	return &cfg, nil
}
