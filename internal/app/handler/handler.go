package handler

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io"
	"log"
	"net/http"
)

type AppRouter struct {
	usecase usecase.InputPort
	*chi.Mux
}

func NewAppRouter(usecase usecase.InputPort) *AppRouter {

	// Root router
	rootRouter := chi.NewRouter()

	// Middlewares
	rootRouter.Use(middleware.Recoverer)
	rootRouter.Use(CustomMiddleware(struct{}{}))

	// configure application router
	appRouter := AppRouter{
		usecase: usecase,
		Mux:     rootRouter,
	}

	appRouter.apiRouter()
	// other sub routers..

	return &appRouter
}

// apiRouter is a sub router which serve public api endpoints
func (a *AppRouter) apiRouter() *chi.Mux {

	apiRouter := chi.NewRouter()

	// Endpoints
	apiRouter.Get("/{id}", a.handleGet)
	apiRouter.Post("/", a.handlePost)

	// Mount sub router to root router
	a.Mount("/", apiRouter)

	return apiRouter
}

// Handlers
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
	if err != nil || len(input) < len("xx.xx") {
		writer.WriteHeader(400)
		return
	}

	id := a.usecase.Shorten(string(input))

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(201)
	_, err = writer.Write([]byte("http://localhost:8080/" + id))
	if err != nil {
		log.Printf("error write answer: %v", err)
	}
}
