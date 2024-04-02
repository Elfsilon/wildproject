package app

import (
	"time"
	"wildproject/internal/app/data/database"
	model "wildproject/internal/app/domain/models"
	"wildproject/internal/app/router"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type AppFlags struct {
	ConfigPath string
}

type App struct {
	cfg *model.Config
}

func New() App {
	return App{}
}

func (a *App) Run(f *AppFlags) {
	a.LoadConfig(f)

	app := fiber.New(fiber.Config{
		ReadTimeout:  a.cfg.Server.ReadTimeout,
		WriteTimeout: a.cfg.Server.WriteTimeout,
		IdleTimeout:  a.cfg.Server.IdleTimeout,
	})

	dbInstance, dbDispose := a.InitDatabase(a.cfg.Database.Conn)
	defer dbDispose()

	sentryDispose := a.InitSentry()
	defer sentryDispose()

	r := router.NewRouter(app, a.cfg)
	r.Setup(dbInstance)

	log.Info("Running up server")
	if err := app.Listen(a.cfg.Server.Addr); err != nil {
		log.Error(err)
	}

	// TODO: Add graceful shutdown
	log.Info("Shut down server")
}

func (a *App) LoadConfig(f *AppFlags) {
	log.Info("loading app config")

	cfg, err := loadConfigFromEnv()
	if err != nil {
		log.Errorf("unable to load app config %s", err)
	}

	a.cfg = cfg
}

func (a *App) InitSentry() func() {
	log.Info("Setting up sentry")

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              a.cfg.Sentry.Dsn,
		Debug:            a.cfg.Sentry.Debug,
		ServerName:       a.cfg.Name,
		TracesSampleRate: a.cfg.Sentry.TracesSampleRate,
		AttachStacktrace: a.cfg.Sentry.AttachStackTrace,
	})
	if err != nil {
		log.Fatalf("sentry init error: %s", err)
	}

	return func() {
		sentry.Flush(2 * time.Second)
	}
}

func (a *App) InitDatabase(conn string) (database.Instance, func()) {
	log.Info("Establishing database connection")

	pg := database.NewPostgres()
	if err := pg.Open(conn); err != nil {
		log.Errorf("open database error: %s", err)
	}

	instance, err := pg.Instance()
	if err != nil {
		log.Errorf("get database instance error: %s", err)
	}

	return instance, func() {
		pg.Close()
	}
}
