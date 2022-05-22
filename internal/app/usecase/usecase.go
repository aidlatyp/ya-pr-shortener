package usecase

import (
	"fmt"
	"log"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

type Repository interface {
	Store(*domain.URL) error
	FindByKey(string) (*domain.URL, error)
	FindAll(string) []*domain.URL
}

type InputPort interface {
	Shorten(string, string) string
	RestoreOrigin(string) (string, error)
	ShowAll(string) ([]*domain.URL, error)
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

func (s *Shorten) Shorten(url string, userID string) string {

	short := s.shortener.MakeShort(url)

	var user *domain.User = nil
	if userID != "" {
		user = &domain.User{
			ID: userID,
		}
		short.Owner = user.ID
	}

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

func (s *Shorten) ShowAll(user string) ([]*domain.URL, error) {
	list := s.repo.FindAll(user)
	if list == nil {
		return nil, fmt.Errorf("seems user %v do not have any links yet", user)
	}
	return list, nil
}
