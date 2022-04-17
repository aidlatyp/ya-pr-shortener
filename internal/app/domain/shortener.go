package domain

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/config"
)

const shortenedUrlLen = config.ShortenedUrlLen

// Shorten restricts the generator return type
// by exact shortenedUrlLen number of elements
type Shorten [shortenedUrlLen]byte

// Generator is an interface used by Shortener to create a random strings
// and do not depend on concrete random generation algorithm
type Generator interface {
	Generate() Shorten
}

// Shortener is a structure which represents main "business logic"
type Shortener struct {
	Generator
}

func NewShortener(generator Generator) *Shortener {
	return &Shortener{
		Generator: generator,
	}
}

func (s *Shortener) MakeShort(inputUrl string) *Url {
	short := s.Generate()
	str := string(short[:])
	u := NewUrl(inputUrl, str)
	return u
}
