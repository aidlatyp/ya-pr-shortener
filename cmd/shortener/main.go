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
)

func main() {

	flags := config.ParseFlags()
	appConf := config.NewAppConfig(&flags)

	// anyway need im-memory as main storage or as a cache
	var store usecase.Repository = storage.NewURLMemoryStorage()

	// Configuration and main is relatively simple right now
	// later if size will grow up - move to separate app struct
	if appConf.IsFilePathSet() {
		persistentStorage, err := storage.NewPersistentStorage(appConf.FilePath, store)
		if err != nil {
			log.Fatalf("can't start in persistent mode %v ", err.Error())
		}
		defer func() {
			if err := persistentStorage.Close(); err != nil {
				log.Print(err)
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
	appRouter := handler.NewAppRouter(appConf.BaseURL, uc)

	// Start
	server := http.Server{
		Addr:              appConf.ServerAddr,
		Handler:           appRouter,
		ReadHeaderTimeout: time.Duration(appConf.ServerTimeout) * time.Second,
		ReadTimeout:       time.Duration(appConf.ServerTimeout) * time.Second,
		WriteTimeout:      time.Duration(appConf.ServerTimeout) * time.Second,
	}

	log.Printf("server starting at %v", appConf.ServerAddr)

	err := server.ListenAndServe()
	log.Printf("server finished with: %v", err)
}
