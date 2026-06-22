package generate

import (
	"crypto/rand"
	"math/big"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	lenAlphabet = len(alphabet)
	CodeLength  = 10
)

type Generator interface {
	Generate() (string, error)
}

type RandomGenerator struct {
	alphabet    string
	lenAlphabet *big.Int
	codeLength  int
}

func NewRandomGenerator() *RandomGenerator {
	return &RandomGenerator{
		alphabet:    alphabet,
		lenAlphabet: big.NewInt(int64(lenAlphabet)),
		codeLength:  CodeLength,
	}
}

func (g *RandomGenerator) Generate() (string, error) {
	result := make([]byte, g.codeLength)

	for i := 0; i < g.codeLength; i++ {
		n, err := rand.Int(rand.Reader, g.lenAlphabet)
		if err != nil {
			return "", err
		}

		result[i] = g.alphabet[n.Int64()]
	}

	return string(result), nil
}

func IsValidCode(code string) bool {
	if len(code) != CodeLength {
		return false
	}

	for _, ch := range code {
		if (ch < 'a' || ch > 'z') &&
			(ch < 'A' || ch > 'Z') &&
			(ch < '0' || ch > '9') &&
			ch != '_' {
			return false
		}
	}

	return true
}
