package randstr

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"math/rand/v2"
)

type Generator struct {
	charset        string
	minLen, maxLen int
}

func New(minLen, maxLen int, charset string) *Generator {
	return &Generator{
		charset: charset,
		minLen:  minLen,
		maxLen:  maxLen,
	}
}

func (g *Generator) RandomString() (string, error) {
	//nolint: gosec
	size := rand.IntN(g.maxLen-g.minLen) + g.minLen
	random := make([]byte, size)

	n, err := cryptorand.Read(random)
	if n != len(random) || err != nil {
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(random)[:min(g.maxLen, len(random))], nil
}
