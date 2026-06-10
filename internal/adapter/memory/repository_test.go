package memory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
	"github.com/google/uuid"
)

func TestCreateOrFindNewURL(t *testing.T) {
	t.Parallel()

	repo := NewRepo()
	input := &model.URL{
		ID:        uuid.New(),
		ShortCode: "Ab3_kLm901",
		Original:  "https://example.com",
		CreatedAt: time.Now().UTC(),
	}

	got, err := repo.CreateOrFind(context.Background(), input)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if got == nil {
		t.Fatal("expected non-nil url")
	}

	if got.Original != input.Original || got.ShortCode != input.ShortCode || got.ID != input.ID {
		t.Fatalf("unexpected saved url: %#v", got)
	}
}

func TestCreateOrFindExistingOriginalReturnsExisting(t *testing.T) {
	t.Parallel()

	repo := NewRepo()
	existing := &model.URL{
		ID:        uuid.New(),
		ShortCode: "Ab3_kLm901",
		Original:  "https://example.com",
		CreatedAt: time.Now().UTC(),
	}

	_, err := repo.CreateOrFind(context.Background(), existing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := &model.URL{
		ID:        uuid.New(),
		ShortCode: "ZZZZZZZZZZ",
		Original:  existing.Original,
		CreatedAt: time.Now().UTC(),
	}

	got, err := repo.CreateOrFind(context.Background(), second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.ShortCode != existing.ShortCode {
		t.Fatalf("expected existing short code %q, got %q", existing.ShortCode, got.ShortCode)
	}
}

func TestCreateOrFindShortCodeCollision(t *testing.T) {
	t.Parallel()

	repo := NewRepo()
	first := &model.URL{
		ID:        uuid.New(),
		ShortCode: "Ab3_kLm901",
		Original:  "https://example.com/one",
		CreatedAt: time.Now().UTC(),
	}

	_, err := repo.CreateOrFind(context.Background(), first)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second := &model.URL{
		ID:        uuid.New(),
		ShortCode: first.ShortCode,
		Original:  "https://example.com/two",
		CreatedAt: time.Now().UTC(),
	}

	_, err = repo.CreateOrFind(context.Background(), second)
	if !errors.Is(err, model.ErrShortCodeExists) {
		t.Fatalf("expected ErrShortCodeExists, got %v", err)
	}
}

func TestCreateOrFindNilInput(t *testing.T) {
	t.Parallel()

	repo := NewRepo()

	_, err := repo.CreateOrFind(context.Background(), nil)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestFindByShortCode(t *testing.T) {
	t.Parallel()

	repo := NewRepo()
	input := &model.URL{
		ID:        uuid.New(),
		ShortCode: "Ab3_kLm901",
		Original:  "https://example.com",
		CreatedAt: time.Now().UTC(),
	}

	_, err := repo.CreateOrFind(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := repo.FindByShortCode(context.Background(), input.ShortCode)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Original != input.Original {
		t.Fatalf("expected original %q, got %q", input.Original, got.Original)
	}
}

func TestFindByShortCodeNotFound(t *testing.T) {
	t.Parallel()

	repo := NewRepo()

	_, err := repo.FindByShortCode(context.Background(), "Ab3_kLm901")
	if !errors.Is(err, model.ErrURLNotFound) {
		t.Fatalf("expected ErrURLNotFound, got %v", err)
	}
}
