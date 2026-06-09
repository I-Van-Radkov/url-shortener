package memory

import (
	"context"
	"sync"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
)

type InMemoryRepository struct {
	mu            sync.RWMutex
	byShortCode   map[string]model.URL
	byOriginalURL map[string]model.URL
}

func New() *InMemoryRepository {
	return &InMemoryRepository{
		byShortCode:   make(map[string]model.URL),
		byOriginalURL: make(map[string]model.URL),
	}
}

func (r *InMemoryRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if url, exists := r.byOriginalURL[originalURL]; exists {
		return url.Clone(), nil
	}

	return nil, model.ErrURLNotFound
}

func (r *InMemoryRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if url, exists := r.byShortCode[shortCode]; exists {
		return url.Clone(), nil
	}

	return nil, model.ErrURLNotFound
}

func (r *InMemoryRepository) Save(ctx context.Context, url *model.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if url == nil {
		return model.ErrInvalidOriginalURL
	}

	if _, exists := r.byOriginalURL[url.Original]; exists {
		return model.ErrOriginalURLExists
	}

	if _, exists := r.byShortCode[url.ShortCode]; exists {
		return model.ErrShortCodeExists
	}

	stored := *url
	r.byOriginalURL[url.Original] = stored
	r.byShortCode[url.ShortCode] = stored

	return nil
}
