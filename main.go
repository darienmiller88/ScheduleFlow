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

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	godotenv.Load()
	database.ConnectToPostgreSQL()

	router := chi.NewRouter()
	indexController := controllers.NewIndexController()

	//Set up middlewares
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.RequestSize(1 << 20))
	router.Use(middleware.Timeout(45 * time.Second))

	router.Mount("/", indexController.Router)
	
	//Serve static files along the "/static" route
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	port := os.Getenv("PORT")

	fmt.Println("Server is running on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}