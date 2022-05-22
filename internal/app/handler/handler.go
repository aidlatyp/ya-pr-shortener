package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	appMiddle "github.com/aidlatyp/ya-pr-shortener/internal/app/handler/middlewares"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/aidlatyp/ya-pr-shortener/internal/config"
	"github.com/go-chi/chi"
	chiMiddle "github.com/go-chi/chi/middleware"
)

const minURLlen = 4

type AppRouter struct {
	usecase usecase.InputPort
	*chi.Mux
	baseURL string
}

func NewAppRouter(baseURL string, usecase usecase.InputPort) *AppRouter {
	// Root router
	rootRouter := chi.NewRouter()
	// Middlewares
	rootRouter.Use(chiMiddle.Recoverer)
	rootRouter.Use(appMiddle.AuthMiddleware)
	rootRouter.Use(appMiddle.CompressMiddleware)

	// configure application router
	appRouter := AppRouter{
		usecase: usecase,
		Mux:     rootRouter,
		baseURL: baseURL,
	}
	appRouter.apiRouter()
	return &appRouter
}

// apiRouter is a sub router which serve public api endpoints
func (a *AppRouter) apiRouter() *chi.Mux {

	apiRouter := chi.NewRouter()

	// Endpoints
	apiRouter.Get("/{id}", a.handleGet)
	apiRouter.Post("/", a.handlePost)

	// api
	apiRouter.Post("/api/shorten", a.handleShorten)
	apiRouter.Get("/api/user/urls", a.handleUserURLs)

	// Mount sub router to root router
	a.Mount("/", apiRouter)

	return apiRouter
}

// Handlers
func (a *AppRouter) handleUserURLs(writer http.ResponseWriter, request *http.Request) {

	ctxUserID, ok := request.Context().Value(config.UserIDCtxKey).(string)
	if !ok {
		writer.WriteHeader(404)
		return
	}

	resultList, err := a.usecase.ShowAll(ctxUserID)
	if err != nil {
		writer.WriteHeader(204)
		return
	}

	type Presentation struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	outputList := make([]Presentation, 0, len(resultList))

	for _, v := range resultList {
		p := Presentation{
			ShortURL:    a.baseURL + v.Short,
			OriginalURL: v.Orig,
		}
		outputList = append(outputList, p)
	}

	marshaled, _ := json.Marshal(outputList)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)

	_, err = writer.Write(marshaled)
	if err != nil {
		log.Printf("error while writing answer: %v", err)
	}

}

func (a *AppRouter) handleShorten(writer http.ResponseWriter, request *http.Request) {

	var ctxUserID string
	ctxUserID, _ = request.Context().Value(config.UserIDCtxKey).(string)

	inputBytes, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		return
	}

	input := make(map[string]string, 1)
	err = json.Unmarshal(inputBytes, &input)
	if err != nil {
		writer.WriteHeader(400)
		return
	}

	if origURL, ok := input["url"]; ok {

		id := a.usecase.Shorten(origURL, ctxUserID)

		output := map[string]string{
			"result": a.baseURL + id,
		}
		marshalled, err := json.Marshal(output)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(201)

		_, err = writer.Write(marshalled)
		if err != nil {
			log.Printf("error while writing answer: %v", err)
		}

	} else {
		writer.WriteHeader(400)
		return
	}
}

func (a *AppRouter) handleGet(writer http.ResponseWriter, request *http.Request) {

	id := chi.URLParam(request, "id")

	response, err := a.usecase.RestoreOrigin(id)
	if err != nil {
		writer.WriteHeader(404)
		return
	}

	writer.Header().Set("Location", response)
	writer.WriteHeader(307)
}

func (a *AppRouter) handlePost(writer http.ResponseWriter, request *http.Request) {

	ctxUserID, _ := request.Context().Value(config.UserIDCtxKey).(string)

	input, err := io.ReadAll(request.Body)
	if err != nil || len(input) < minURLlen {
		writer.WriteHeader(400)
		return
	}

	id := a.usecase.Shorten(string(input), ctxUserID)

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(201)

	_, err = writer.Write([]byte(a.baseURL + id))
	if err != nil {
		log.Printf("error while writing answer: %v", err)
	}
}
