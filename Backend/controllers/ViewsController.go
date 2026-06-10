package controllers

import (
	"html/template"
	"github.com/go-chi/chi/v5"
)

type ViewsController struct {
	Router *chi.Mux
}

func NewViewsController() *ViewsController {
	c := &ViewsController{
		Router: chi.NewRouter(),
	}

	c.registerRoutes()

	return c
}

func (v *ViewsController) registerRoutes() {
	panic("unimplemented")
}
