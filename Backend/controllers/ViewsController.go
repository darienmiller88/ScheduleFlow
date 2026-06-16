package controllers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct {
	Router    *chi.Mux
	templates *template.Template
}

func NewViewsController() *ViewsController {
	partialFiles, err := filepath.Glob("./templates/partials/*.html")

	if err != nil {
		panic(err)
	}

	partialFiles = append(partialFiles, "./templates/home.html")

	vc := &ViewsController{
		Router:    chi.NewRouter(),
		templates: template.Must(template.ParseFiles(partialFiles...)),
	}

	vc.registerViewRoutes()

	return vc
}

func (v *ViewsController) registerViewRoutes() {
	v.Router.Get("/home", v.homePage)
}

func (v *ViewsController) homePage(response http.ResponseWriter, request *http.Request) {
	err := v.templates.ExecuteTemplate(response, "home.html", nil)

	if err != nil {
		http.Error(response, "Error rendering template", http.StatusInternalServerError)
	}
}
