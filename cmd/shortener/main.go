package main

import (
	"log"
	"net/http"
	"time"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/handler"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/storage"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/storage/postgres"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/aidlatyp/ya-pr-shortener/internal/config"
	"github.com/aidlatyp/ya-pr-shortener/internal/util"
)

func main() {
	// configure from flags, env or by default
	appConf := config.NewAppConfig()

	// choose storage depending on if specified filepath or not
	store := storage.NewStorage(appConf.FilePath)

	// connect database if connect string configured
	pg, err := postgres.NewDB(appConf.DBConnect)
	if err != nil {
		log.Printf("can't start database due to: %v", err.Error())
	} else {
		store = pg
		defer func() {
			if err := pg.Close(); err != nil {
				log.Print(err)
			}
		}()
	}

	// Domain
	gen := util.GetShortenGenerator()
	shortener := domain.NewShortener(gen)

	// Use_cases
	shortenUsecase := usecase.NewShorten(shortener, store)
	dbCheckUsecase := usecase.NewLiveliness(pg)

	// Application Router
	appRouter := handler.NewAppRouter(
		appConf.BaseURL,
		shortenUsecase,
		dbCheckUsecase,
	)

	// Start
	server := http.Server{
		Addr:              appConf.ServerAddr,
		Handler:           appRouter,
		ReadHeaderTimeout: time.Duration(appConf.ServerTimeout) * time.Second,
		ReadTimeout:       time.Duration(appConf.ServerTimeout) * time.Second,
		WriteTimeout:      time.Duration(appConf.ServerTimeout) * time.Second,
	}
	log.Printf("server is starting at %v", appConf.ServerAddr)

	err = server.ListenAndServe()
	log.Printf("server finished with: %v", err)
}
