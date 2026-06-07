package controllers

import (
	"fmt"
	"net/http"

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

func (c *IndexController) registerRoutes() {
	c.Router.Get("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(response, "Hello, World!")
	})
}