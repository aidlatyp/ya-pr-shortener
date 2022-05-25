package usecase

import (
	"errors"
	"fmt"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

type Repository interface {
	Store(*domain.URL) error
	FindByKey(string) (*domain.URL, error)
	FindAll(string) []*domain.URL
	BatchWrite([]domain.URL) error
}

type InputPort interface {
	Shorten(string, string) (string, error)
	RestoreOrigin(string) (string, error)
	ShowAll(string) ([]*domain.URL, error)
	ShortenBatch(input []Correlation, user string) ([]OutputBatchItem, error)
}

type Shorten struct {
	shortener *domain.Shortener
	repo      Repository
}

func NewShorten(
	shortener *domain.Shortener,
	repo Repository,
) *Shorten {
	return &Shorten{
		shortener: shortener,
		repo:      repo,
	}
}

type Correlation struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type OutputBatchItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
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

type ErrAlreadyExists struct {
	Err            error
	ExistShortenID string
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("url %v already exists", e.ExistShortenID)
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

		fmt.Println(err.Error())

		if errors.As(err, &ErrAlreadyExists{}) {
			// process error
			return "", err
		}
	}
	return short.Short, nil
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
	if list == nil || len(list) < 1 {
		return nil, fmt.Errorf("seems user %v do not have any links yet", user)
	}
	return list, nil
}
