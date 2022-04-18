package main

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/handler"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/repository"
	"github.com/aidlatyp/ya-pr-shortener/internal/config"
	"github.com/aidlatyp/ya-pr-shortener/internal/util"
	"net/http"
	"time"
)

//func Middleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
//		// do something
//		next.ServeHTTP(writer, request)
//	})
//}

func main() {

	repo := repository.NewURLRepo()

	generator := util.GenFunc(util.Generate)
	service := domain.NewShortener(&generator)

	appHandler := handler.NewAppHandler(service, repo)

	mux := http.NewServeMux()
	mux.Handle("/", appHandler.HandleMain())

	server := http.Server{
		Addr:              config.ServerAddr,
		Handler:           mux, //Middleware(mux),
		ReadHeaderTimeout: config.ServerTimeout * time.Second,
		ReadTimeout:       config.ServerTimeout * time.Second,
		WriteTimeout:      config.ServerTimeout * time.Second,
	}

	server.ListenAndServe()

}
