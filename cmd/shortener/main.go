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

// Configuration and main is relatively simple right now
// serve state with global vars
// later if size will up grow - move to separate app struct
var serverConf config.ServerConf
var appConf config.AppConf

func init() {
	flags := config.NewParsedFlags()

	serverConf = config.NewServerConf(&flags)
	appConf = config.NewAppConf(&flags, serverConf.ServerAddr)
}

func main() {
	// Choose storage mode dep on conf
	var store usecase.Repository

	// anyway need im-memory as main storage or as a cache
	store = storage.NewURLMemoryStorage()

	if appConf.IsFilePathSet() {
		// Connect persistence
		persistentStorage, err := storage.NewPersistentStorage(appConf.FilePath, store)
		if err != nil {
			log.Fatalf("cant start corrupted url file %v ", err.Error())
		}

		defer func() {
			if err := persistentStorage.Close(); err != nil {
				log.Printf("error while closing file, potentional data loss, not closed with %v", err)
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
		Addr:              serverConf.ServerAddr,
		Handler:           appRouter,
		ReadHeaderTimeout: time.Duration(serverConf.ServerTimeout) * time.Second,
		ReadTimeout:       time.Duration(serverConf.ServerTimeout) * time.Second,
		WriteTimeout:      time.Duration(serverConf.ServerTimeout) * time.Second,
	}

	log.Printf("server starting at %v", serverConf.ServerAddr)

	err := server.ListenAndServe()
	log.Printf("server finished with: %v", err)
}
