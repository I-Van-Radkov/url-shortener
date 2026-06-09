package memory

import (
	"context"
	"sync"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
)

type InMemoryRepository struct {
	mu            sync.RWMutex
	byShortCode   map[string]*model.URL
	byOriginalUrl map[string]*model.URL
}

func New() *InMemoryRepository {
	return &InMemoryRepository{
		byShortCode:   make(map[string]*model.URL),
		byOriginalUrl: make(map[string]*model.URL),
	}
}

func (r *InMemoryRepository) FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if url, exists := r.byOriginalUrl[originalURL]; exists {
		return url, nil
	}

	return nil, model.ErrURLNotFound
}

func (r *InMemoryRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if url, exists := r.byShortCode[shortCode]; exists {
		return url, nil
	}

	return nil, model.ErrURLNotFound
}

func (r *InMemoryRepository) Save(ctx context.Context, url *model.URL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byOriginalUrl[url.Original]; exists {
		return model.ErrOriginalURLExists
	}

	if _, exists := r.byShortCode[url.ShortCode]; exists {
		return model.ErrShortCodeExists
	}

	r.byOriginalUrl[url.Original] = url
	r.byShortCode[url.ShortCode] = url

	return nil
}
