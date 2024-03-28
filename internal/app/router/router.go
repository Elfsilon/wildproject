package router

import (
	"time"
	"wildproject/internal/app/data/database"
	repo "wildproject/internal/app/data/repositories"
	manager "wildproject/internal/app/domain/managers"
	model "wildproject/internal/app/domain/models"
	service "wildproject/internal/app/domain/services"
	controller "wildproject/internal/app/router/controllers"
	"wildproject/internal/app/router/middleware"

	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Router struct {
	app *fiber.App
	cfg *model.Config
}

func NewRouter(app *fiber.App, cfg *model.Config) Router {
	return Router{app, cfg}
}

func (r *Router) Setup(db database.Instance) error {
	log.Info("Setting up router")

	tm := manager.NewJwtManager(r.cfg.Auth.AuthJwtSecret, r.cfg.Auth.AccessTokenTTL)

	ur := repo.NewUsers(db)
	sr := repo.NewSessions(db)

	us := service.NewUsers(ur)
	ss := service.NewSessions(&r.cfg.Auth, sr, tm)

	uc := controller.NewUsers(us)
	sc := controller.NewSessions(ss, us)

	// Setup middlewares
	sentryMiddleware := fibersentry.New(fibersentry.Config{
		WaitForDelivery: true,
		Timeout:         10 * time.Second,
	})

	// sm := func(c *fiber.Ctx) error {
	// 	if hub := fibersentry.GetHubFromContext(c); hub != nil {
	// 		hub.Scope().setta()
	// 		// Add current route, user_id, headers, session_id, ...
	// 	}

	// 	return c.Next()
	// }

	authGuard := middleware.NewAuthGuard(ss, tm)

	// Setup routes
	r.app.Use(sentryMiddleware)

	api := r.app.Group("/api")
	api.Get("/health", controller.HealthCheck)

	v1 := api.Group("/v1")

	// Unprotected routes
	unUsers := v1.Group("/users")
	unUsers.Post("/", uc.Create)

	unSessions := v1.Group("/sessions")
	unSessions.Post("/", sc.Create)
	unSessions.Put("/", authGuard.RefreshGuard, sc.Refresh)

	// Protected routes
	protected := v1.Group("/protected", authGuard.AccessGuard)

	users := protected.Group("/users")

	user := users.Group("/:user_id<guid>")
	user.Get("/", uc.GetInfo)
	user.Put("/", uc.Update)
	user.Delete("/", uc.Delete)

	sessions := user.Group("/sessions")
	sessions.Get("/", sc.GetAllByUserID)
	sessions.Delete("/", sc.DropAll)

	session := sessions.Group("/:session_id<int>")
	session.Get("/", sc.GetByID)
	session.Delete("/", sc.Drop)

	return nil
}
