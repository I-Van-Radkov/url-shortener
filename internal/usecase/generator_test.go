package usecase

import (
	"strings"
	"testing"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
)

func TestCodeGeneratorGenerate_LengthAndCharset(t *testing.T) {
	t.Parallel()

	generator := NewCodeGenerator()

	for i := 0; i < 100; i++ {
		code, err := generator.Generate()
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		if len(code) != model.ShortCodeLength {
			t.Fatalf("expected code length %d, got %d", model.ShortCodeLength, len(code))
		}

		for _, ch := range code {
			if !strings.ContainsRune(model.Charset, ch) {
				t.Fatalf("generated code contains invalid character %q", ch)
			}
		}
	}
}
