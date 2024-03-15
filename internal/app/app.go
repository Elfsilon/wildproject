package app

import (
	"os"
	m "temp/internal/app/models"
	"temp/internal/app/router"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gopkg.in/yaml.v3"
)

type AppFlags struct {
	ConfigPath string
}

type App struct {
	cfg *m.Config
}

func New() App {
	return App{}
}

func (a *App) Init(f *AppFlags) {
	log.Info("loading app config")
	if err := a.LoadConfigFromFile(f.ConfigPath); err != nil {
		log.Fatalf("unable to load app config %s", err)
	}
}

func (a *App) Run(f *AppFlags) {
	a.Init(f)

	app := fiber.New(fiber.Config{
		ReadTimeout:  a.cfg.Server.ReadTimeout,
		WriteTimeout: a.cfg.Server.WriteTimeout,
		IdleTimeout:  a.cfg.Server.IdleTimeout,
	})

	r := router.NewRouter(app, a.cfg)
	r.Setup()

	app.Listen(a.cfg.Server.Addr)
}

func (a *App) LoadConfigFromFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg m.RawConfig

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}

	a.cfg = &m.Config{
		Version: cfg.Version,
		Env:     cfg.Env,
		Database: m.DatabaseConfig{
			Conn: cfg.Database.Conn,
		},
		Server: m.ServerConfig{
			Addr:         cfg.Server.Addr,
			ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		},
		Auth: m.AuthConfig{
			AccessTokenTTL:  time.Duration(cfg.Auth.AccessTokenTTL) * time.Minute,
			RefreshTokenTTL: time.Duration(cfg.Auth.RefreshTokenTTL) * time.Minute,
		},
	}

	return err
}
