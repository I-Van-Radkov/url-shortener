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

func NewRepo() *InMemoryRepository {
	return &InMemoryRepository{
		byShortCode:   make(map[string]model.URL),
		byOriginalURL: make(map[string]model.URL),
	}
}

func (r *InMemoryRepository) CreateOrFind(ctx context.Context, urlInput *model.URL) (*model.URL, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if url, exists := r.byOriginalURL[urlInput.Original]; exists {
		return &url, nil
	}

	if _, exists := r.byShortCode[urlInput.ShortCode]; exists {
		return nil, model.ErrShortCodeExists
	}

	url := *urlInput
	r.byShortCode[urlInput.ShortCode] = url
	r.byOriginalURL[urlInput.Original] = url

	return &url, nil
}

func (r *InMemoryRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if url, exists := r.byShortCode[shortCode]; exists {
		return &url, nil
	}

	return nil, model.ErrURLNotFound
}
