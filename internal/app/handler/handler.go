package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	appMiddle "github.com/aidlatyp/ya-pr-shortener/internal/app/handler/middlewares"
	"github.com/aidlatyp/ya-pr-shortener/internal/app/usecase"
	"github.com/go-chi/chi"
	chiMiddle "github.com/go-chi/chi/middleware"
)

const minURLlen = 4

type AppRouter struct {
	usecase    usecase.InputPort
	liveliness *usecase.Liveliness
	*chi.Mux
	baseURL string
}

func NewAppRouter(
	baseURL string,
	appUsecase usecase.InputPort,
	liveliness *usecase.Liveliness,
) *AppRouter {

	// Root router
	rootRouter := chi.NewRouter()

	// Root Middlewares
	rootRouter.Use(chiMiddle.Recoverer)
	rootRouter.Use(appMiddle.AuthMiddleware)

	// configure application router
	appRouter := AppRouter{
		usecase:    appUsecase,
		Mux:        rootRouter,
		baseURL:    baseURL,
		liveliness: liveliness,
	}

	appRouter.apiRouter()
	appRouter.infraRouter()

	return &appRouter
}

// apiRouter is a sub router which serve public api endpoints
func (a *AppRouter) apiRouter() {

	apiRouter := chi.NewRouter()

	// compress api endpoints only
	apiRouter.Use(appMiddle.CompressMiddleware)

	// Endpoints
	apiRouter.Get("/{id}", a.handleGet)
	apiRouter.Post("/", a.handlePost)

	// api
	apiRouter.Post("/api/shorten", a.handleShorten)
	apiRouter.Get("/api/user/urls", a.handleUserURLs)
	apiRouter.Post("/api/shorten/batch", a.handleBatch)

	// Mount sub router to root router
	a.Mount("/", apiRouter)
}

// infraRouter is a sub router which serve infrastructure endpoints
func (a *AppRouter) infraRouter() {
	infraRouter := chi.NewRouter()
	infraRouter.Get("/", a.handlePing)
	a.Mount("/ping", infraRouter)
}

func (a *AppRouter) handleBatch(writer http.ResponseWriter, request *http.Request) {

	ctxUserID, ok := request.Context().Value(appMiddle.UserIDCtxKey).(string)
	if !ok {
		writer.WriteHeader(404)
		return
	}

	inputBytes, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		return
	}

	inputCollection := make([]usecase.Correlation, 0)
	err = json.Unmarshal(inputBytes, &inputCollection)
	if err != nil {
		log.Printf("cant unmarshal due to %v", err)
	}

	outputList, err := a.usecase.ShortenBatch(inputCollection, ctxUserID)
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	for index := range outputList {
		outputList[index].ShortURL = a.baseURL + outputList[index].ShortURL
	}

	marshaled, _ := json.Marshal(outputList)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(201)

	_, err = writer.Write(marshaled)
	if err != nil {
		log.Printf("error while writing answer: %v", err)
	}

	//for k, v := range inputCollection {
	//	fmt.Println(k, v)
	//}

}

// Handlers
func (a *AppRouter) handlePing(writer http.ResponseWriter, _ *http.Request) {
	err := a.liveliness.Do()
	if err != nil {
		writer.WriteHeader(500)
	}
	writer.WriteHeader(200)
}

func (a *AppRouter) handleUserURLs(writer http.ResponseWriter, request *http.Request) {

	log.Println("fetch all foe user")

	ctxUserID, ok := request.Context().Value(appMiddle.UserIDCtxKey).(string)
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
	ctxUserID, _ = request.Context().Value(appMiddle.UserIDCtxKey).(string)

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
		id, err := a.usecase.Shorten(origURL, ctxUserID)
		if err != nil {
			if errors.As(err, &usecase.ErrAlreadyExists{}) {
				e := err.(usecase.ErrAlreadyExists)
				id = e.ExistShortenID
				writer.Header().Set("Content-Type", "text/plain")
				writer.WriteHeader(409)
				_, err = writer.Write([]byte(a.baseURL + id))
				if err != nil {
					log.Printf("error while writing answer: %v", err)
				}
				return
			}
		}

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

	ctxUserID, _ := request.Context().Value(appMiddle.UserIDCtxKey).(string)

	input, err := io.ReadAll(request.Body)
	if err != nil || len(input) < minURLlen {
		writer.WriteHeader(400)
		return
	}

	id, err := a.usecase.Shorten(string(input), ctxUserID)
	if err != nil {
		if errors.As(err, &usecase.ErrAlreadyExists{}) {
			e := err.(usecase.ErrAlreadyExists)
			id = e.ExistShortenID
			writer.Header().Set("Content-Type", "text/plain")
			writer.WriteHeader(409)
			_, err = writer.Write([]byte(a.baseURL + id))
			if err != nil {
				log.Printf("error while writing answer: %v", err)
			}
			return
		}
	}

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(201)

	_, err = writer.Write([]byte(a.baseURL + id))
	if err != nil {
		log.Printf("error while writing answer: %v", err)
	}
}
