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

var serverConf config.ServerConf
var appConf config.App

func init() {
	flags := config.NewParsedFlags()
	serverConf = config.NewServerConf(flags)
	appConf = config.NewAppConf(flags, serverConf)
}

func main() {

	var store usecase.Repository
	// In memory data provider
	if appConf.FilePath == "" {
		store = storage.NewURLMemoryStorage()
		log.Println("no filename, fallback to in memory")
	} else {
		// Persistent data provider
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

	log.Printf("server starting at %v", serverConf.ServerAddr)

	err := server.ListenAndServe()
	log.Printf("server finished with: %v", err)
}
