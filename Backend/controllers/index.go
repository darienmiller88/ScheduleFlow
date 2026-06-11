package controllers

import (
	"github.com/go-chi/chi/v5"
)

type IndexController struct{
	Router *chi.Mux
}

func NewIndexController() *IndexController {
	c := &IndexController{
		Router: chi.NewRouter(),
	}

	c.registerRoutes()

	return c
}

// registerRoutes sets up the routes for the IndexController by mounting all sub-routes in controllers.
func (c *IndexController) registerRoutes() {
	vc := NewViewsController()

	c.Router.Mount("/", vc.Router)
}