package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"ScheduleFlow/Backend/database"
	"ScheduleFlow/Backend/controllers"
)

func main() {
	godotenv.Load()
	database.ConnectToMongoDB()

	router := chi.NewRouter()
	indexController := controllers.NewIndexController()
	
	router.Mount("/", indexController.Router)

	port := os.Getenv("PORT")

	fmt.Println("Server is running on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}