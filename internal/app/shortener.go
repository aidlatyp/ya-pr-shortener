package shortener

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/aidlatyp/ya-pr-shortener/internal/repositories"
	"github.com/aidlatyp/ya-pr-shortener/internal/utils"
)

const (
	letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterCount = 6
)

func GetNewShortener(repo repositories.Repository) *Shortener {
	return &Shortener{repo: repo}
}

type Shortener struct {
	repo repositories.Repository
}

func (s *Shortener) getRandomURL(longURL string) (string, error) {

	b := make([]byte, letterCount)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	res := string(b)
	err := s.repo.Set(longURL, res)

	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *Shortener) MakeShortURL(longURL, baseURL string) (string, error) {

	if !utils.IsValidURL(longURL) {
		return "", errors.New("uncorrect URL format")
	}

	// check if already exists
	shortURL, err := s.repo.GetByKey(longURL)

	if err != nil {
		return "", err
	}

	if len(shortURL) == 0 {
		shortURL, err = s.getRandomURL(longURL)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s/%s", baseURL, shortURL), nil
}

func (s *Shortener) GetRawURL(shortURL string) (string, error) {

	shortURLs, err := s.repo.GetAll()

	if err != nil {
		return "", err
	}

	for longValue, shortValue := range shortURLs {
		if shortValue == shortURL {
			return longValue, nil
		}
	}

	return "", fmt.Errorf("no URL for id = %s", shortURL)
}

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func (s *Shortener) GetAllURL(baseURL string) ([]*URL, error) {

	urlsMap, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	urls := make([]*URL, 0, 0)

	for longURL, shortURL := range urlsMap {
		urls = append(urls, &URL{
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, shortURL),
			OriginalURL: longURL,
		})
	}

	return urls, nil
}
