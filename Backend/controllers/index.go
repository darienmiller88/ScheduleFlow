package controllers

import (
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/alexedwards/scs/v2"
)

type IndexController struct{
	Router *chi.Mux
}

func NewIndexController(db *sqlx.DB, sessionManager *scs.SessionManager) *IndexController {
	c := &IndexController{
		Router: chi.NewRouter(),
	}

	c.registerRoutes(db, sessionManager)

	return c
}

// registerRoutes sets up the routes for the IndexController by mounting all sub-routes in controllers.
func (c *IndexController) registerRoutes(db *sqlx.DB, sessionManager *scs.SessionManager) {
	vc := NewViewsController(sessionManager)
	sc := NewSpecialistController(
		services.NewSpecialistService(repositories.NewSpecialistRepository(db)),
		services.NewEmailVerificationService(repositories.NewEmailVerificationRepository(db)),
		services.NewEmailSendService(),
		sessionManager,
	)

	c.Router.Mount("/", vc.Router)
	c.Router.Mount("/specialists", sc.Router)
}