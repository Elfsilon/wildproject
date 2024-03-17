package router

import (
	ctr "temp/internal/app/controllers"
	db "temp/internal/app/database"
	m "temp/internal/app/models"
	rep "temp/internal/app/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Repositories struct {
	Users *rep.Users
}

type Controllers struct {
	Users *ctr.Users
}

type Components struct {
	Database     db.Database
	Repositories Repositories
	Controllers  Controllers
}

type Router struct {
	app *fiber.App
	cfg *m.Config
}

func NewRouter(app *fiber.App, cfg *m.Config) Router {
	return Router{app, cfg}
}

func (r *Router) setupAppComponents(db *db.Database) (*Components, error) {
	log.Info("Setting up app components")
	rp := Repositories{
		Users: rep.NewUsers(db.DB),
	}

	c := Controllers{
		Users: ctr.NewUsers(rp.Users),
	}

	components := &Components{
		Repositories: rp,
		Controllers:  c,
	}

	return components, nil
}

func (r *Router) Setup(db *db.Database) error {
	c, err := r.setupAppComponents(db)
	if err != nil {
		return err
	}

	log.Info("Setting up app routes")

	api := r.app.Group("/api")
	api.Get("/health", ctr.HealthCheck)

	v1 := api.Group("/v1")

	users := v1.Group("/users")
	users.Post("/", c.Controllers.Users.Create)

	user := users.Group("/:user_id<guid>")
	user.Get("/", c.Controllers.Users.GetInfo)
	user.Put("/", c.Controllers.Users.Update)
	user.Delete("/", c.Controllers.Users.Delete)

	return nil
}
