package handler

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"io"
	"net/http"
	"strings"
)

type AppHandler struct {
	usecase usecase.InputPort
}

func NewAppHandler(usecase usecase.InputPort) *AppHandler {
	return &AppHandler{
		usecase: usecase,
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

	token := strings.Split(request.URL.Path, "/")
	if len(token) == 2 && token[1] != "" {

		response, err := a.usecase.RestoreOrigin(token[1])
		if err != nil {
			writer.WriteHeader(404)
			return
		}

		writer.Header().Set("Location", response)
		writer.WriteHeader(307)
		return
	}

	writer.WriteHeader(400)
}

func (a *AppHandler) handleMainPost(writer http.ResponseWriter, request *http.Request) {

	if request.URL.Path == "/" {

		input, err := io.ReadAll(request.Body)
		if err != nil || len(input) < len("xx.xx") {
			writer.WriteHeader(400)
			return
		}

		id := a.usecase.Shorten(string(input))

		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(201)
		writer.Write([]byte("http://localhost:8080/" + id))

		return
	}
	writer.WriteHeader(400)
}
