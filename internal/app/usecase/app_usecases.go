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
	shortener    *domain.Shortener
	repo         Repository
	deletionChan chan userURLsToDelete
	buf          []userURLsToDelete
	timer        *time.Timer
	isTimeout    bool
}

func NewShorten(shortener *domain.Shortener, repo Repository) *Shorten {
	s := Shorten{
		shortener: shortener,
		repo:      repo,
		// deletion request goes into this channel
		deletionChan: make(chan userURLsToDelete),
		// buffer to accumulate delete requests
		buf: make([]userURLsToDelete, 0, bufLen),
		// was timeout exceeded from first deletion request
		isTimeout: true,
		// timer to count batch time
		timer: time.NewTimer(0),
	}
	// listener for users deletions
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

// Represent separate user request to delete bunch of urls
type userURLsToDelete struct {
	UserID string
	delIDs []string
}

func (s *Shorten) flush(delBuf []userURLsToDelete) {
	for _, v := range delBuf {
		s.repo.BatchDelete(v.delIDs, v.UserID)
	}
}

//  listener listens for user deletions
func (s *Shorten) deletionListener() {
	for {
		select {

		// buffer is full flush
		case delRequest := <-s.deletionChan:
			if s.isTimeout {
				s.timer.Reset(time.Second * timeout)
				s.isTimeout = false
			}
			s.buf = append(s.buf, delRequest)
			if len(s.buf) >= bufLen {
				cp := make([]userURLsToDelete, len(s.buf))
				copy(cp, s.buf)
				s.buf = make([]userURLsToDelete, 0)
				go s.flush(cp)
				s.timer.Stop()
				s.isTimeout = true
			}

		// timer fired flush
		case <-s.timer.C:
			if len(s.buf) > 0 {
				cp := make([]userURLsToDelete, len(s.buf))
				copy(cp, s.buf)
				s.buf = make([]userURLsToDelete, 0)
				go s.flush(cp)
			}
			s.isTimeout = true
		}
	}
}

func (s *Shorten) DeleteBatch(delIDs []string, user string) {
	delRequest := userURLsToDelete{
		UserID: user,
		delIDs: delIDs,
	}
	// send income del requests to [deletionListener]
	// via channel
	go func() {
		s.deletionChan <- delRequest
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
