package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"ScheduleFlow/Backend/controllers"
	"ScheduleFlow/Backend/database"
	"ScheduleFlow/Backend/services"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager

func main() {
	godotenv.Load()
	database.ConnectToPostgreSQL()

	// Initialize a new session manager and configure the session lifetime for one week.
	sessionManager = scs.New()
	sessionManager.Lifetime = 168 * time.Hour

	//Initialize router
	router := chi.NewRouter()
	indexController := controllers.NewIndexController()

	emailService := services.NewEmailSendService()

	if err := emailService.SendEmail(); err != nil {
		fmt.Println("Error sending email:", err)
	}

	//Set up middlewares
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.ClientIPFromRemoteAddr)
	router.Use(middleware.RequestSize(1 << 20))
	router.Use(middleware.Timeout(45 * time.Second))

	//Mount index router, which has all other controllers mounted on it
	router.Mount("/", indexController.Router)
	
	//Serve static files along the "/static" route
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	port := os.Getenv("PORT")

	fmt.Println("Server is running on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), sessionManager.LoadAndSave(router))
}

