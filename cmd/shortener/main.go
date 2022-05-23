package main

import (
	"log"

	"github.com/aidlatyp/ya-pr-shortener/internal/server"
	"github.com/aidlatyp/ya-pr-shortener/internal/utils"
)

func main() {

	handler, err := server.Run(utils.GetConfigs())

	defer func() {
		for _, closers := range handler.ReposClosers {
			err = closers()
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	if err != nil {
		log.Fatal(err)
	}
}
