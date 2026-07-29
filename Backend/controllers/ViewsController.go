package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

type ViewsController struct {
	Router    *chi.Mux
	templates map[string]*template.Template
}

func NewViewsController() *ViewsController {
	partials, _ := filepath.Glob("./templates/partials/*.html")
	pages, _ := filepath.Glob("./templates/pages/*.html")

	tmplMap := make(map[string]*template.Template)

	for _, page := range pages {
		// Get page name without extension (e.g., "home", "login")
		name := strings.TrimSuffix(filepath.Base(page), ".html")

		fmt.Println("name:", name)

		// // Build file slice specifically for THIS page: base + partials + page
		// files := append([]string{base, page}, partials...)

		// // Parse the isolated template set for this page
		// tmplMap[name] = template.Must(template.ParseFiles(files...))
		files := []string{"templates/Base.html"}
		files = append(files, partials...)
		files = append(files, fmt.Sprintf("templates/pages/%s.html", name))

		tmplMap[name] = template.Must(template.ParseFiles(files...))
	}

	vc := &ViewsController{
		Router:    chi.NewRouter(),
		templates: tmplMap,
	}

	vc.registerViewRoutes()

	return vc
}

func (v *ViewsController) registerViewRoutes() {
	v.Router.Get("/home", v.homePage)
	v.Router.Get("/", v.loginPage)
}

func (v *ViewsController) homePage(res http.ResponseWriter, req *http.Request) {
	if err := v.templates["home"].Execute(res, nil); err != nil {
		http.Error(res, "Error rendering template", http.StatusInternalServerError)
	}
}

func (v *ViewsController) loginPage(res http.ResponseWriter, req *http.Request) {
	if err := v.templates["login"].Execute(res, nil); err != nil {
		http.Error(res, "Error rendering template", http.StatusInternalServerError)
	}
}
