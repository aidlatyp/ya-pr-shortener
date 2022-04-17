package util

import (
	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
	"math/rand"
	"time"
)

const Symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Generate() domain.Shorten {
	rand.Seed(time.Now().UnixNano())
	var buf domain.Shorten
	for i := range buf {
		randomIndex := rand.Intn(len(Symbols))
		buf[i] = Symbols[randomIndex]
	}
	return buf
}

// GenFunc function wrapper to satisfy the Generator interface
type GenFunc func() domain.Shorten

func (gf GenFunc) Generate() domain.Shorten {
	return gf()
}
