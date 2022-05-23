package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	shortener "github.com/aidlatyp/ya-pr-shortener/internal/app"
	"github.com/aidlatyp/ya-pr-shortener/internal/middlewares"
	"github.com/aidlatyp/ya-pr-shortener/internal/repositories"
	"github.com/aidlatyp/ya-pr-shortener/internal/storage"
	"github.com/aidlatyp/ya-pr-shortener/internal/utils"
	"github.com/go-chi/chi/v5"
)

type PostURLRequest struct {
	URL string `json:"url"`
}

type PostURLResponse struct {
	Result string `json:"result"`
}

type Handler struct {
	configs        *utils.Config
	userShorteners map[string]*shortener.Shortener
	ReposClosers   []func() error
}

func NewHandler(configs *utils.Config) *Handler {
	return &Handler{
		configs:        configs,
		userShorteners: make(map[string]*shortener.Shortener),
		ReposClosers:   make([]func() error, 0),
	}
}

func (h *Handler) getUserShortener() (*shortener.Shortener, error) {

	userID := fmt.Sprintf("%x", middlewares.UserID)

	if sh, ok := h.userShorteners[userID]; ok {
		return sh, nil
	} else {

		var repo repositories.Repository
		var err error

		if len(h.configs.FileStoragePath) == 0 {
			repo = storage.NewMemoryStorage()
		} else {
			repo, err = storage.NewFileStorage(h.configs.FileStoragePath)
			if err != nil {
				return nil, err
			}
		}

		h.userShorteners[userID] = sh.GetNewShortener(repo)
		h.ReposClosers = append(h.ReposClosers, repo.CloseResources)

		return h.userShorteners[userID], nil
	}
}

func (h *Handler) PostShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	p := PostURLRequest{}

	if err := json.Unmarshal(b, &p); err != nil {
		http.Error(w, "Incorrect body JSON format", http.StatusBadRequest)
		return
	}

	if len(p.URL) == 0 {
		http.Error(w, "URL can not be empty", http.StatusBadRequest)
		return
	}

	sh, err := h.getUserShortener()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	shortURL, err := sh.MakeShortURL(p.URL, h.configs.BaseURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp, err := json.Marshal(PostURLResponse{Result: shortURL})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (h *Handler) GetURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")

	if len(id) == 0 {
		http.Error(w, "Need to set id", http.StatusBadRequest)
		return
	}

	sh, err := h.getUserShortener()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	longURL, err := sh.GetRawURL(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Location", longURL)
	w.Header().Set("content-type", "text/plain; charset=utf-8")

	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) SaveURLHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	rawURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	sh, err := h.getUserShortener()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	shortURL, err := sh.MakeShortURL(string(rawURL), h.configs.BaseURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (h *Handler) GetAllSavedURLs(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed by this route!", http.StatusMethodNotAllowed)
		return
	}

	sh, err := h.getUserShortener()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	urls, err := sh.GetAllURL(h.configs.BaseURL)
	if err != nil {
		http.Error(w, "Errors happens when get all saved URLS!", http.StatusInternalServerError)
		return
	}
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}
