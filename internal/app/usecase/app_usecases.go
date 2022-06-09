package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

type Repository interface {
	Store(*domain.URL) error
	FindByKey(string) (*domain.URL, error)
	FindAll(string) []*domain.URL
	BatchWrite([]domain.URL) error
	BatchDelete([]string, string) error
}

type InputPort interface {
	Shorten(string, string) (string, error)
	RestoreOrigin(string) (string, error)
	ShowAll(string) ([]*domain.URL, error)
	ShortenBatch(input []Correlation, user string) ([]OutputBatchItem, error)
	DeleteBatch(input []string, user string) error
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

func (s *Shorten) ShortenBatch(input []Correlation, user string) ([]OutputBatchItem, error) {

	output := make([]OutputBatchItem, 0)
	urls := make([]domain.URL, 0)
	for _, inputPair := range input {

		url := s.shortener.MakeShort(inputPair.OriginalURL)
		url.Owner = user

		out := OutputBatchItem{
			CorrelationID: inputPair.CorrelationID,
			ShortURL:      url.Short,
		}
		urls = append(urls, *url)
		output = append(output, out)
	}

	err := s.repo.BatchWrite(urls)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *Shorten) Shorten(url string, userID string) (string, error) {
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
		if errors.As(err, &ErrAlreadyExists{}) {
			return "", err
		}
	}
	return short.Short, nil
}

func (s *Shorten) RestoreOrigin(id string) (string, error) {
	url, err := s.repo.FindByKey(id)
	if err != nil {
		// process errors
		return "", err
	}
	return url.Orig, nil
}

func (s *Shorten) DeleteBatch(delIDs []string, user string) error {
	fmt.Println("usecase del")
	err := s.repo.BatchDelete(delIDs, user)
	if err != nil {
		// process
		log.Print("error deleting", err)
	}
	return nil
}

func (s *Shorten) ShowAll(user string) ([]*domain.URL, error) {
	list := s.repo.FindAll(user)
	if list == nil || len(list) < 1 {
		return nil, fmt.Errorf("seems you %v do not have any shortened links yet", user)
	}
	return list, nil
}
