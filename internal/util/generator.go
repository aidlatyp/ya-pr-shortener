package util

import (
	"math/rand"
	"time"

	"github.com/aidlatyp/ya-pr-shortener/internal/app/domain"
)

const symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// generator utility source of random bytes
func generator(l int) []byte {
	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, l)
	for i := range buf {
		randomIndex := rand.Intn(len(symbols))
		buf[i] = symbols[randomIndex]
	}
	return buf
}

// GenerateShorten generates random short bytes with exactly domain.ShortenedURLLen
func GenerateShorten() domain.Shorten {
	//return *(*[domain.ShortenedURLLen]byte)(generator(domain.ShortenedURLLen))
	id := generator(domain.ShortenedURLLen)
	var s domain.Shorten
	copy(s[:], id)
	return s
}

// genFunc function wrapper to satisfy the domain.Generator interface
type genFunc func() domain.Shorten

func (gf genFunc) Generate() domain.Shorten {
	return gf()
}

func GetShortenGenerator() domain.Generator {
	return genFunc(GenerateShorten)
}

// GenerateUserID generate bytes to represent user id
func GenerateUserID() []byte {
	return generator(6)
}
