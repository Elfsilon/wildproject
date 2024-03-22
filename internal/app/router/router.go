package router

import (
	ctr "temp/internal/app/controllers"
	db "temp/internal/app/database"
	model "temp/internal/app/domain/models"
	rep "temp/internal/app/repositories"
	tokenmanager "temp/internal/app/token-manager"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type Repositories struct {
	Users    *rep.Users
	Sessions *rep.Sessions
}

type Controllers struct {
	Users    *ctr.Users
	Sessions *ctr.Sessions
}

type Services struct {
	TokenManager *tokenmanager.TokenManager
}

type Components struct {
	Repositories Repositories
	Controllers  Controllers
	Services     Services
}

type Router struct {
	app *fiber.App
	cfg *model.Config
}

func NewRouter(app *fiber.App, cfg *model.Config) Router {
	return Router{app, cfg}
}

func (r *Router) setupAppComponents(db *db.Database) (*Components, error) {
	log.Info("Setting up app components")
	s := Services{
		TokenManager: tokenmanager.NewTokenManager(r.cfg.Auth.AuthJwtSecret, r.cfg.Auth.AccessTokenTTL),
	}

	rp := Repositories{
		Users:    rep.NewUsers(db.DB),
		Sessions: rep.NewSessions(db.DB),
	}

	c := Controllers{
		Users:    ctr.NewUsers(rp.Users),
		Sessions: ctr.NewSessions(&r.cfg.Auth, rp.Sessions, rp.Users, s.TokenManager),
	}

	components := &Components{
		Repositories: rp,
		Controllers:  c,
		Services:     s,
	}

	return components, nil
}

func (r *Router) Setup(db *db.Database) error {
	c, err := r.setupAppComponents(db)
	if err != nil {
		return err
	}
	log.Info("Setting up app routes")
	authGuard := NewAuthGuard(c.Services.TokenManager, c.Repositories.Sessions)

	api := r.app.Group("/api")
	api.Get("/health", ctr.HealthCheck)

	v1 := api.Group("/v1")

	// Unprotected routes
	unUsers := v1.Group("/users")
	unUsers.Post("/", c.Controllers.Users.Create)

	unUser := unUsers.Group("/:user_id<guid>")

	unSessions := unUser.Group("/sessions")
	unSessions.Post("/", c.Controllers.Sessions.New)
	unSessions.Put("/", authGuard.RefreshGuard, c.Controllers.Sessions.Refresh)

	// Protected routes
	protected := v1.Group("/protected", authGuard.AccessGuard)

	users := protected.Group("/users")

	user := users.Group("/:user_id<guid>")
	user.Get("/", c.Controllers.Users.GetInfo)
	user.Put("/", c.Controllers.Users.Update)
	user.Delete("/", c.Controllers.Users.Delete)

	sessions := user.Group("/sessions")
	sessions.Get("/", c.Controllers.Sessions.GetAllByUserID)
	sessions.Delete("/", c.Controllers.Sessions.DropAll)
	sessions.Get("/:session_id<int>/", c.Controllers.Sessions.Get)
	sessions.Delete("/:session_id<int>/", c.Controllers.Sessions.Drop)

	return nil
}
