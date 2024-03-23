package app

import (
	"temp/internal/app/data/database"
	model "temp/internal/app/domain/models"
	"temp/internal/app/router"

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

	db := database.NewPostgres()
	if err := db.Open(a.cfg.Database.Conn); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbInstance, err := db.Instance()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Setting up router")

	r := router.NewRouter(app, a.cfg)
	r.Setup(dbInstance)

	log.Info("Running up server")
	if err := app.Listen(a.cfg.Server.Addr); err != nil {
		log.Error(err)
	}

	// TODO: Add graceful shutdown
	log.Info("Shut down server")
}
