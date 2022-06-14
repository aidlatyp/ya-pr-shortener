package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

const bufLen = 3
const timeout = 5

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
	DeleteBatch(input []string, user string)
}

type Shorten struct {
	shortener *domain.Shortener
	repo      Repository
}

func NewShorten(shortener *domain.Shortener, repo Repository) *Shorten {
	s := Shorten{
		shortener: shortener,
		repo:      repo,
	}
	go s.deletionListener()
	return &s
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
			// deletionListener error at this layer if needed
			// i.e set user-friendly message
			return "", err
		}
	}
	return short.Short, nil
}

func (s *Shorten) RestoreOrigin(id string) (string, error) {
	url, err := s.repo.FindByKey(id)
	if err != nil {
		if errors.As(err, &ErrURLDeleted{}) {
			// deletionListener error at this layer if needed
			// i.e set user-friendly message
			return "", err
		}
		return "", err
	}
	return url.Orig, nil
}

type userURLsToDelete struct {
	UserId string
	delIDs []string
}

var inputChan = make(chan userURLsToDelete, 0)
var buf = make([]userURLsToDelete, 0, bufLen)
var timer = time.NewTimer(0 * time.Second)
var isTimeout = true

func (s *Shorten) flush(delBuf []userURLsToDelete) {
	for _, v := range delBuf {
		s.repo.BatchDelete(v.delIDs, v.UserId)
	}
}

func (s *Shorten) deletionListener() {
	for {
		select {
		case delRequest := <-inputChan:
			if isTimeout {
				timer.Reset(time.Second * timeout)
				isTimeout = false
			}
			buf = append(buf, delRequest)
			if len(buf) >= bufLen {
				cp := make([]userURLsToDelete, len(buf))
				copy(cp, buf)
				buf = make([]userURLsToDelete, 0)
				go s.flush(cp)
				timer.Stop()
				isTimeout = true
			}

		case <-timer.C:
			if len(buf) > 0 {
				cp := make([]userURLsToDelete, len(buf))
				copy(cp, buf)
				buf = make([]userURLsToDelete, 0)
				go s.flush(cp)
			}
			isTimeout = true
		}
	}
}

func (s *Shorten) DeleteBatch(delIDs []string, user string) {
	delRequest := userURLsToDelete{
		UserId: user,
		delIDs: delIDs,
	}
	go func() {
		inputChan <- delRequest
	}()
}

func (s *Shorten) ShowAll(user string) ([]*domain.URL, error) {
	list := s.repo.FindAll(user)
	if list == nil || len(list) < 1 {
		// deletionListener error at this layer if needed
		// i.e set user-friendly message
		return nil, fmt.Errorf("seems you %v do not have any shortened links yet", user)
	}
	return list, nil
}
