package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"x-platform-emailer/internal/config"
	"x-platform-emailer/internal/lib/logger/sl"
	"x-platform-emailer/internal/server/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Envs
	configPath := os.Getenv("CONFIG_PATH")

	// Config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	// Logger
	log := sl.SetupLogger(cfg.Env)
	log.Info("initializing server", slog.String("address", cfg.HTTPServer.Address))

	// Router
	router := chi.NewRouter()

	// Handlers
	router.Post("/send", handlers.NewSendHandler(log, &cfg.Mailbox))

	// Middleware
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	// Server
	server := http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Startup
	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))
	if err = server.ListenAndServe(); err != nil {
		os.Exit(-1)
	}
}
