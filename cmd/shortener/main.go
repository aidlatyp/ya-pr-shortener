package main

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/handler"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/storage"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/aidlatyp/ya-pr-shortener/internal/config"
	"github.com/aidlatyp/ya-pr-shortener/internal/util"
	"log"
	"net/http"
	"time"
)

func main() {

	// Domain
	gen := util.GetGenerator()
	shortener := domain.NewShortener(gen)

	// Storage (Data Provider)
	store := storage.NewURLStorage()

	// Usecase
	uc := usecase.NewShorten(shortener, store)

	// Main application handler
	appHandler := handler.NewAppRouter(uc)

	// Server
	server := http.Server{
		Addr:              config.ServerAddr,
		Handler:           appHandler,
		ReadHeaderTimeout: config.ServerTimeout * time.Second,
		ReadTimeout:       config.ServerTimeout * time.Second,
		WriteTimeout:      config.ServerTimeout * time.Second,
	}
	err := server.ListenAndServe()
	log.Printf("server finished with: %v", err)
}
