package router

import (
	ctr "temp/internal/app/controllers"
	m "temp/internal/app/models"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app *fiber.App
	cfg *m.Config
}

func NewRouter(app *fiber.App, cfg *m.Config) Router {
	return Router{app, cfg}
}

func (r *Router) Setup() {
	api := r.app.Group("/api")
	api.Get("/health", ctr.HealthCheck)

	// v1 := api.Group("/v1")
}
