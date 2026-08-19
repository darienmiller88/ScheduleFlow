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
	"ScheduleFlow/Backend/sessionManager"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	godotenv.Load()
	database.ConnectToPostgreSQL()

	//initialize session manager
	sm := sessionmanager.NewSessionManager(database.DB)
	
	//Initialize router
	router := chi.NewRouter()
	indexController := controllers.NewIndexController(database.DB, sm)

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
	http.ListenAndServe(fmt.Sprintf(":%s", port), sm.LoadAndSave(router))
}

