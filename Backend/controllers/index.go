package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type IndexController struct{
	Router *chi.Mux
}

func (c *IndexController) RegisterRoutes() {
	c.Router = chi.NewRouter()

	c.Router.Get("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(response, "Hello, World!")
	})
}