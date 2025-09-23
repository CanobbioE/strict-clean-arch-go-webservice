// Package main is the application entry point.
// Wires up all dependency and handle injection, starts the HTTP server, and sets environment config.
package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/config"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/infrastructure/db"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/infrastructure/webservice"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/interfaces/controller"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/interfaces/presenter"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/usecase/interactor"
)

func main() {
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		panic("CONFIG_PATH environment variable not set")
	}

	cfg, err := config.Load(configFile)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	repo := db.NewInMemoryBookRepo(logger)
	interact := interactor.NewBookInteractor(logger, repo)
	bookPresenter := presenter.NewBookPresenter(logger)
	errPresenter := presenter.NewErrorPresenter(logger)
	ctl := controller.NewBookController(logger, interact, bookPresenter, errPresenter)

	router := webservice.NewHandler(ctl)

	s := &http.Server{
		Addr:              cfg.ServerAddress,
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}
	logger.Info("Starting bookshop service on " + cfg.ServerAddress)

	log.Fatal(s.ListenAndServe())
}
