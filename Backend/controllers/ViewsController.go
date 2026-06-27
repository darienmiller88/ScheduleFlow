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
	partials, _ := filepath.Glob("./templates/partials/*.html")
	pages, _    := filepath.Glob("./templates/pages/*.html")
	htmlFiles   := append(partials, pages...)

	vc := &ViewsController{
		Router:    chi.NewRouter(),
		templates: template.Must(template.ParseFiles(htmlFiles...)),
	}

	vc.registerViewRoutes()

	return vc
}

func (v *ViewsController) registerViewRoutes() {
	v.Router.Get("/home", v.homePage)
	v.Router.Get("/", v.loginPage)
}


func (v *ViewsController) homePage(response http.ResponseWriter, request *http.Request) {
	err := v.templates.ExecuteTemplate(response, "home.html", nil)

	if err != nil {
		http.Error(response, "Error rendering template", http.StatusInternalServerError)
	}
}

func (v *ViewsController) loginPage(response http.ResponseWriter, request *http.Request) {
	err := v.templates.ExecuteTemplate(response, "login.html", nil)

	if err != nil {
		http.Error(response, "Error rendering template", http.StatusInternalServerError)
	}
}
