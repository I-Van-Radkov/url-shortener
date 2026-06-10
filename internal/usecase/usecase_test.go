package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
	"github.com/google/uuid"
)

type repositoryStub struct {
	findByShortCodeFn func(ctx context.Context, shortCode string) (*model.URL, error)
	createOrFindFn    func(ctx context.Context, url *model.URL) (*model.URL, error)
}

func (r *repositoryStub) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	return r.findByShortCodeFn(ctx, shortCode)
}

func (r *repositoryStub) CreateOrFind(ctx context.Context, url *model.URL) (*model.URL, error) {
	return r.createOrFindFn(ctx, url)
}

type generatorStub struct {
	generateFn func() (string, error)
}

func (g *generatorStub) Generate() (string, error) {
	return g.generateFn()
}

func TestSaveOriginalURLSuccess(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			return url, nil
		},
	}
	gen := &generatorStub{
		generateFn: func() (string, error) {
			return "Ab3_kLm901", nil
		},
	}

	uc := NewUsecase(repo, gen, 3)

	got, err := uc.SaveOriginalURL(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ShortCode != "Ab3_kLm901" {
		t.Fatalf("expected short code %q, got %q", "Ab3_kLm901", got.ShortCode)
	}

	if got.Original != "https://example.com" {
		t.Fatalf("expected original URL to be trimmed and saved, got %q", got.Original)
	}

	if got.ID == uuid.Nil {
		t.Fatal("expected generated UUID")
	}

	if got.CreatedAt.IsZero() {
		t.Fatal("expected created at to be set")
	}
}

func TestSaveOriginalURLTrimsSpaces(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			if url.Original != "https://example.com" {
				t.Fatalf("expected trimmed original URL, got %q", url.Original)
			}
			return url, nil
		},
	}
	gen := &generatorStub{
		generateFn: func() (string, error) {
			return "Ab3_kLm901", nil
		},
	}

	uc := NewUsecase(repo, gen, 3)

	_, err := uc.SaveOriginalURL(context.Background(), "   https://example.com   ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveOriginalURLInvalidURL(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			t.Fatal("repository must not be called for invalid URL")
			return nil, nil
		},
	}
	gen := &generatorStub{
		generateFn: func() (string, error) {
			t.Fatal("generator must not be called for invalid URL")
			return "", nil
		},
	}

	uc := NewUsecase(repo, gen, 3)

	_, err := uc.SaveOriginalURL(context.Background(), "not-a-url")
	if !errors.Is(err, model.ErrInvalidOriginalURL) {
		t.Fatalf("expected ErrInvalidOriginalURL, got %v", err)
	}
}

func TestSaveOriginalURLRetryOnShortCodeCollision(t *testing.T) {
	t.Parallel()

	generated := []string{"first_code", "Ab3_kLm901"}
	genCalls := 0
	gen := &generatorStub{
		generateFn: func() (string, error) {
			code := generated[genCalls]
			genCalls++
			return code, nil
		},
	}

	repoCalls := 0
	repo := &repositoryStub{
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			repoCalls++
			if repoCalls == 1 {
				return nil, model.ErrShortCodeExists
			}
			return url, nil
		},
	}

	uc := NewUsecase(repo, gen, 3)

	got, err := uc.SaveOriginalURL(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ShortCode != "Ab3_kLm901" {
		t.Fatalf("expected retry to return second short code, got %q", got.ShortCode)
	}

	if genCalls != 2 {
		t.Fatalf("expected generator to be called twice, got %d", genCalls)
	}
}

func TestSaveOriginalURLGenerationFailed(t *testing.T) {
	t.Parallel()

	gen := &generatorStub{
		generateFn: func() (string, error) {
			return "Ab3_kLm901", nil
		},
	}

	repo := &repositoryStub{
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			return nil, model.ErrShortCodeExists
		},
	}

	uc := NewUsecase(repo, gen, 2)

	_, err := uc.SaveOriginalURL(context.Background(), "https://example.com")
	if !errors.Is(err, model.ErrGenerationFailed) {
		t.Fatalf("expected ErrGenerationFailed, got %v", err)
	}
}

func TestSaveOriginalURLGeneratorError(t *testing.T) {
	t.Parallel()

	wantErr := fmt.Errorf("generator failed")
	gen := &generatorStub{
		generateFn: func() (string, error) {
			return "", wantErr
		},
	}

	repo := &repositoryStub{
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			t.Fatal("repository must not be called when generator fails")
			return nil, nil
		},
	}

	uc := NewUsecase(repo, gen, 3)

	_, err := uc.SaveOriginalURL(context.Background(), "https://example.com")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected generator error, got %v", err)
	}
}

func TestGetOriginalURLSuccess(t *testing.T) {
	t.Parallel()

	want := &model.URL{
		ID:        uuid.New(),
		ShortCode: "Ab3_kLm901",
		Original:  "https://example.com",
		CreatedAt: time.Now().UTC(),
	}

	repo := &repositoryStub{
		findByShortCodeFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			return want, nil
		},
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
	}

	gen := &generatorStub{generateFn: func() (string, error) { return "", nil }}
	uc := NewUsecase(repo, gen, 3)

	got, err := uc.GetOriginalURL(context.Background(), want.ShortCode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Original != want.Original {
		t.Fatalf("expected original %q, got %q", want.Original, got.Original)
	}
}

func TestGetOriginalURLInvalidShortCode(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{
		findByShortCodeFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			t.Fatal("repository must not be called for invalid short code")
			return nil, nil
		},
		createOrFindFn: func(ctx context.Context, url *model.URL) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
	}

	gen := &generatorStub{generateFn: func() (string, error) { return "", nil }}
	uc := NewUsecase(repo, gen, 3)

	_, err := uc.GetOriginalURL(context.Background(), "short")
	if !errors.Is(err, model.ErrInvalidShortCodeLength) {
		t.Fatalf("expected ErrInvalidShortCodeLength, got %v", err)
	}
}
