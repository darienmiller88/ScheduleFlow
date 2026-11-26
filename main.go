package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"

	"ScheduleFlow/Backend/controllers"
	"ScheduleFlow/Backend/database"
	"ScheduleFlow/Backend/middlewares"
	"ScheduleFlow/Backend/sessionManager"

	"github.com/go-chi/chi/v5/middleware"
	// "github.com/go-co-op/gocron/v2"
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
	router.Use(middleware.RequestSize(5 << 20))
	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(httprate.LimitBy(300, time.Minute, middlewares.ClientIPKey))


	//Mount index router, which has all other controllers mounted on it
	router.Mount("/", indexController.Router)

	//Serve static files along the "/static" route
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	port := os.Getenv("PORT")

	fmt.Println("Server is running on port:", port)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: sm.LoadAndSave(router),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	//Listen for OS signals to gracefully shut down the server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Wait for the OS signal to shut down the server gracefully
	<-ctx.Done()

	// Create a new context with a timeout of 10 seconds for the shutdown process
    ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()

	// Gracefully shut down the server with a timeout of 10 seconds. If the server does not shut down
	// within this time, it will forcefully close.
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("HTTP server shutdown error: %v", err)
    }
}