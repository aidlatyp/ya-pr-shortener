package util

import (
	"math/rand"
	"time"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Generate() domain.Shorten {
	rand.Seed(time.Now().UnixNano())
	var buf domain.Shorten
	for i := range buf {
		randomIndex := rand.Intn(len(symbols))
		buf[i] = symbols[randomIndex]
	}
	return buf
}

// genFunc function wrapper to satisfy the Generator interface
type genFunc func() domain.Shorten

func (gf genFunc) Generate() domain.Shorten {
	return gf()
}

func GetGenerator() domain.Generator {
	return genFunc(Generate)
}

// GenerateUserID todo: remove
func GenerateUserID() []byte {
	id := generator(6)
	return id
}

func generator(l int) []byte {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, l)
	for i := range buf {
		randomIndex := rand.Intn(len(symbols))
		buf[i] = symbols[randomIndex]
	}
	return buf
}
