package handler

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/repository"
	"io"
	"log"
	"net/http"
	"strings"
)

type AppHandler struct {
	shortener *domain.Shortener
	repo      *repository.URLRepo
}

func NewAppHandler(shortener *domain.Shortener, repo *repository.URLRepo) *AppHandler {
	return &AppHandler{
		shortener: shortener,
		repo:      repo,
	}
}

func (a *AppHandler) HandleMain() http.HandlerFunc {

	return func(writer http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			a.handleMainGet(writer, request)
		case http.MethodPost:
			a.handleMainPost(writer, request)
		default:
			writer.WriteHeader(400)
		}
	}
}

func (a *AppHandler) handleMainGet(writer http.ResponseWriter, request *http.Request) {

	urlTokens := strings.Split(request.URL.Path, "/")
	if len(urlTokens) == 2 && urlTokens[1] != "" {

		url, err := a.repo.FindByKey(urlTokens[1])
		if err != nil {
			writer.WriteHeader(400)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Header().Set("Location", url.Orig)
		writer.WriteHeader(307)
		return
	}

	writer.WriteHeader(400)
}

func (a *AppHandler) handleMainPost(writer http.ResponseWriter, request *http.Request) {

	if request.URL.Path == "/" {
		url, err := io.ReadAll(request.Body)
		if err != nil || len(url) < len("xx.xx") {
			writer.WriteHeader(400)
			return
		}
		u := a.shortener.MakeShort(string(url))
		err = a.repo.Store(u)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(400)
			return
		}
		writer.WriteHeader(201)
		writer.Write([]byte("http://localhost:8080/" + u.Short))
		return
	}
	writer.WriteHeader(400)
}
