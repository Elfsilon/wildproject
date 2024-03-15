package app

import (
	m "temp/internal/app/models"
	"temp/internal/app/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
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

	cfg, err := loadConfigFromEnv()
	if err != nil {
		log.Fatalf("unable to load app config %s", err)
	}

	a.cfg = cfg
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
