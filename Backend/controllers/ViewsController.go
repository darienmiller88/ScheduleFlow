package controllers

import (
	"ScheduleFlow/Backend/middlewares"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type ViewsController struct {
	db          	  *sqlx.DB
	Router            *chi.Mux
	pageTemplates     map[string]*template.Template
	sessionManger     *scs.SessionManager
	standardTemplates *template.Template
}

func NewViewsController(sessionManager *scs.SessionManager, db *sqlx.DB) *ViewsController {
	partials, _ := filepath.Glob("./templates/partials/*.html")
	pages, _ := filepath.Glob("./templates/pages/*.html")

	tmplMap := make(map[string]*template.Template)
	standardTemplates := template.Must(template.ParseGlob("./templates/*.html"))

	for _, page := range pages {
		// Get page name without extension (e.g., "home", "login")
		name := strings.TrimSuffix(filepath.Base(page), ".html")

		// Build file slice specifically for THIS page: base + partials + page
		files := []string{"templates/Base.html"}

		// append the partial files to the build
		files = append(files, partials...)

		//add the page to the files to be parsed
		files = append(files, fmt.Sprintf("templates/pages/%s.html", name))

		tmplMap[name] = template.Must(template.ParseFiles(files...))
	}

	vc := &ViewsController{
		Router:            chi.NewRouter(),
		pageTemplates:     tmplMap,
		sessionManger:     sessionManager,
		standardTemplates: standardTemplates,
	}

	vc.registerViewRoutes()

	return vc
}

func (v *ViewsController) registerViewRoutes() {
	v.Router.With(middlewares.RequireAuth(v.sessionManger)).Get("/home", v.homePage)
	v.Router.With(middlewares.SendBackToHome(v.sessionManger)).Get("/", v.loginPage)
	v.Router.With(middlewares.RequireAuth(v.sessionManger), middlewares.RequireVerification(v.sessionManger, v.db)).Get("/verification", v.verificationPage)
	v.Router.NotFound(v.notFound)
}

func (v *ViewsController) homePage(res http.ResponseWriter, req *http.Request) {
	if err := v.pageTemplates["home"].Execute(res, nil); err != nil {
		http.Error(res, "Error rendering template", http.StatusInternalServerError)
	}
}

func (v *ViewsController) loginPage(res http.ResponseWriter, req *http.Request) {
	if err := v.pageTemplates["login"].Execute(res, nil); err != nil {
		http.Error(res, "Error rendering template", http.StatusInternalServerError)
	}
}

func (v *ViewsController) verificationPage(res http.ResponseWriter, req *http.Request) {
	if err := v.standardTemplates.ExecuteTemplate(res, "verification.html", nil); err != nil {
		http.Error(res, "Error rendering template", http.StatusInternalServerError)
	}
}

func (v *ViewsController) notFound(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)

	if err := v.standardTemplates.ExecuteTemplate(res, "notfound.html", nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
