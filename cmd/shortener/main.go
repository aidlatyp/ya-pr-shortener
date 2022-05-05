package main

import (
	"log"
	"net/http"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/handler"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/storage"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/aidlatyp/ya-pr-shortener/internal/config"
	"github.com/aidlatyp/ya-pr-shortener/internal/util"
	"github.com/caarlos0/env/v6"
)

func main() {

	// Config
	var serverConf config.Server
	var appConf config.App

	err := env.Parse(&serverConf)
	if err != nil {
		log.Fatalf("can't load server config")
	}

	log.Println("addr is - >", serverConf.ServerAddr)
	log.Println("host is - >", serverConf.ServerHost)
	log.Println("port is - >", serverConf.ServerPort)

	if serverConf.ServerAddr == "" {
		serverConf.ServerAddr = "8080"
	}

	err = env.Parse(&appConf)
	if err != nil {
		log.Printf("can't load application config")
	}

	serverConf.ServerTimeout = 30

	// Domain
	gen := util.GetGenerator()
	shortener := domain.NewShortener(gen)

	// Data Provider
	store := storage.NewURLStorage()

	// Usecase
	uc := usecase.NewShorten(shortener, store)

	// Router
	bu := appConf.BaseURL + "/"
	appRouter := handler.NewAppRouter(bu, uc)

	server := http.Server{
		Addr:    "127.0.0.1" + serverConf.ServerAddr,
		Handler: appRouter,
		//ReadHeaderTimeout: time.Duration(serverConf.ServerTimeout) * time.Second,
		//ReadTimeout:       time.Duration(serverConf.ServerTimeout) * time.Second,
		//WriteTimeout:      time.Duration(serverConf.ServerTimeout) * time.Second,
	}

	err = server.ListenAndServe()
	log.Printf("server finished with: %v", err)
}
