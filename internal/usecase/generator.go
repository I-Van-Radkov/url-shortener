package usecase

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
)

type CodeGenerator struct {
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (cg *CodeGenerator) Generate() (string, error) {
	result := make([]byte, model.ShortCodeLength)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(model.Charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate code: %w", err)
		}
		result[i] = model.Charset[n.Int64()]
	}

	return string(result), nil
}
