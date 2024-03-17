package app

import (
	db "temp/internal/app/database"
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
	log.Info("Initializing app")
	a.Init(f)

	app := fiber.New(fiber.Config{
		ReadTimeout:  a.cfg.Server.ReadTimeout,
		WriteTimeout: a.cfg.Server.WriteTimeout,
		IdleTimeout:  a.cfg.Server.IdleTimeout,
	})

	log.Info("Establishing database connection")
	log.Infof("conn: %v", a.cfg.Database.Conn)
	database := db.NewDatabase()
	if err := database.Open(a.cfg.Database.Conn); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	log.Info("Setting up router")
	r := router.NewRouter(app, a.cfg)
	r.Setup(database)

	log.Info("Running up server")
	if err := app.Listen(a.cfg.Server.Addr); err != nil {
		log.Error(err)
	}
	log.Info("Shut down server")
}
