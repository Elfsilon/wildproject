package m

import "time"

type RawConfig struct {
	Version  string `yaml:"version"`
	Env      string `yaml:"env"`
	Database struct {
		Conn string `yaml:"conn_string"`
	}
	Server struct {
		Addr         string `yaml:"addr"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
		IdleTimeout  int    `yaml:"idle_timeout"`
	}
	Auth struct {
		AccessTokenTTL  int `yaml:"access_token_ttl"`
		RefreshTokenTTL int `yaml:"refresh_token_ttl"`
	}
}

type Config struct {
	Version  string
	Env      string
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
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
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}
