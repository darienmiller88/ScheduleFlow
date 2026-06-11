package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct {
	Router    *chi.Mux
	templates *template.Template
}

func NewViewsController() *ViewsController {
	c := &ViewsController{
		Router: chi.NewRouter(),
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	c.registerViewRoutes()

	return c
}

func (v *ViewsController) registerViewRoutes() {
	v.Router.Get("/home", v.HomePage)
}

func (v *ViewsController) HomePage(response http.ResponseWriter, request *http.Request) {
	err := v.templates.ExecuteTemplate(response, "home.html", nil)
	
	if err != nil {
		fmt.Println("Error rendering template:", err)
		http.Error(response, "Error rendering template", http.StatusInternalServerError)
	}
}