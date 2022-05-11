package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/handler"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/storage"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/aidlatyp/ya-pr-shortener/internal/config"
	"github.com/aidlatyp/ya-pr-shortener/internal/util"
	"github.com/caarlos0/env/v6"
)

func main() {

	// Configure server
	var serverConf config.Server
	err := env.Parse(&serverConf)
	if err != nil {
		log.Fatalf("can't load server config")
	}
	if serverConf.ServerAddr == "" {
		serverConf.ServerAddr = ":8080"
	}

	// Configure application
	var appConf config.App
	err = env.Parse(&appConf)
	if err != nil {
		log.Printf("can't load application config")
	}
	if appConf.BaseURL == "" {
		appConf.BaseURL = "http://localhost" + serverConf.ServerAddr
	}

	// Choose storage
	var store usecase.Repository

	// In memory data provider
	if appConf.FilePath == "" {
		store = storage.NewURLMemoryStorage()
		log.Println("no filename, fallback to in memory")

	} else {

		persistentStorage, err := storage.NewPersistentStorage(appConf.FilePath)
		if err != nil {
			log.Fatalf("cant start corrupted url file %v ", err.Error())
		}

		defer func() {
			if err := persistentStorage.Close(); err != nil {
				log.Printf("error while closing file,  not closed with %v", err)
			}
		}()

		store = persistentStorage
	}

	// Domain
	gen := util.GetGenerator()
	shortener := domain.NewShortener(gen)

	// Usecase
	uc := usecase.NewShorten(shortener, store)

	// Router
	bu := appConf.BaseURL + "/"
	appRouter := handler.NewAppRouter(bu, uc)

	// Start
	server := http.Server{
		Addr:              serverConf.ServerAddr,
		Handler:           appRouter,
		ReadHeaderTimeout: time.Duration(serverConf.ServerTimeout) * time.Second,
		ReadTimeout:       time.Duration(serverConf.ServerTimeout) * time.Second,
		WriteTimeout:      time.Duration(serverConf.ServerTimeout) * time.Second,
	}

	err = server.ListenAndServe()
	log.Printf("server finished with: %v", err)

}
