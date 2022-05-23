package main

import (
	"log"

	"github.com/aidlatyp/ya-pr-shortener/internal/server"
	"github.com/aidlatyp/ya-pr-shortener/internal/utils"
)

func main() {

	configs, err := utils.GetConfigs()
	if err != nil {
		log.Fatal(err)
	}

	resourcesCloser, err := server.Run(configs)
	defer func() {
		if resourcesCloser != nil {
			resourcesCloser()
		}
	}()

	if err != nil {
		resourcesCloser()
		log.Fatal(err)
	}
}
