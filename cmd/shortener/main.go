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

	appConf := config.NewAppConfig()

	store := storage.NewStorage(appConf.FilePath)

	// db
	//var dbCheckUsecase *usecase.Liveliness
	//
	//if appConf.DBConnect != "" {
	//	pg, err := postgres.NewDB(appConf.DBConnect)
	//	if err != nil {
	//		log.Printf("can't start database %v", err.Error())
	//	} else {
	//		store = pg
	//		defer func() {
	//			if err := pg.Close(); err != nil {
	//				log.Print(err)
	//			}
	//		}()
	//	}
	//	dbCheckUsecase = usecase.NewLiveliness(pg)
	//}

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

	// end db

	// Domain
	gen := util.GetGenerator()
	shortener := domain.NewShortener(gen)

	// Usecases
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
