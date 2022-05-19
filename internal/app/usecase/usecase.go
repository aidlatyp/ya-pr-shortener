package usecase

import (
	"log"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

type Repository interface {
	Store(*domain.URL) error
	FindByKey(string) (*domain.URL, error)
}

type InputPort interface {
	Shorten(string) string
	RestoreOrigin(string) (string, error)
}

type Shorten struct {
	shortener *domain.Shortener
	repo      Repository
}

func NewShorten(shortener *domain.Shortener, repo Repository) *Shorten {
	return &Shorten{
		shortener: shortener,
		repo:      repo,
	}
}

func (s *Shorten) Shorten(url string) string {
	short := s.shortener.MakeShort(url)
	err := s.repo.Store(short)
	if err != nil {
		// process an error in the future
		log.Println(err.Error())
	}
	return short.Short
}

func (s *Shorten) RestoreOrigin(id string) (string, error) {
	url, err := s.repo.FindByKey(id)
	if err != nil {
		return "", err
	}
	return url.Orig, nil
}
