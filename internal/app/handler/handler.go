package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	rootRouter.Use(middleware.Recoverer)
	rootRouter.Use(CompressMiddleware)

	//func(next http.Handler) http.Handler

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
	// Api
	apiRouter.Post("/api/shorten", a.handleShorten)

	// Mount sub router to root router
	a.Mount("/", apiRouter)

	return apiRouter
}

// Handlers
func (a *AppRouter) handleShorten(writer http.ResponseWriter, request *http.Request) {

	fmt.Println("HANDLER REQUEST ")

	inputBytes, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Println("READ ALL ERROR ")
		log.Println(err)
		writer.WriteHeader(400)
		return
	}

	input := make(map[string]string, 1)
	err = json.Unmarshal(inputBytes, &input)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(400)
		return
	}

	if origURL, ok := input["url"]; ok {

		id := a.usecase.Shorten(origURL)

		output := map[string]string{
			"result": a.baseURL + id,
		}
		marshalled, err := json.Marshal(output)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		fmt.Println("SET HEADER")

		writer.WriteHeader(201)
		fmt.Println("WRITE HEADERS")

		_, err = writer.Write(marshalled)
		if err != nil {
			log.Printf("error while writing answer: %v", err)
		}

		fmt.Println("WRITE DONE ")

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

	input, err := io.ReadAll(request.Body)
	if err != nil || len(input) < minURLlen {
		writer.WriteHeader(400)
		return
	}

	id := a.usecase.Shorten(string(input))

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(201)
	_, err = writer.Write([]byte(a.baseURL + id))
	if err != nil {
		log.Printf("error while writing answer: %v", err)
	}
}
