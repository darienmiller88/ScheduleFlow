package controllers

import (
	"ScheduleFlow/Backend/repositories"
	"ScheduleFlow/Backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type IndexController struct{
	Router *chi.Mux
}

func NewIndexController(db *sqlx.DB) *IndexController {
	c := &IndexController{
		Router: chi.NewRouter(),
	}

	c.registerRoutes(db)

	return c
}

// registerRoutes sets up the routes for the IndexController by mounting all sub-routes in controllers.
func (c *IndexController) registerRoutes(db *sqlx.DB) {
	vc := NewViewsController()
	sc := NewSpecialistController(services.NewSpecialistService(repositories.NewSpecialistRepository(db)))

	c.Router.Mount("/", vc.Router)
	c.Router.Mount("/specialists", sc.Router)
}