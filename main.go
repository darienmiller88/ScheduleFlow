package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"ScheduleFlow/Backend/database"
)

func main() {
	godotenv.Load()
	database.ConnectToMongoDB()

	router := chi.NewRouter()

	router.Get("/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(response, "Hello, World!")
	})

	port := os.Getenv("PORT")

	fmt.Println("Server is running on port:", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}