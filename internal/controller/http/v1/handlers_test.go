package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type usecaseStub struct {
	saveOriginalURLFn func(ctx context.Context, originalURL string) (*model.URL, error)
	getOriginalURLFn  func(ctx context.Context, shortCode string) (*model.URL, error)
}

func (u *usecaseStub) SaveOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	return u.saveOriginalURLFn(ctx, originalURL)
}

func (u *usecaseStub) GetOriginalURL(ctx context.Context, shortCode string) (*model.URL, error) {
	return u.getOriginalURLFn(ctx, shortCode)
}

func TestCreateShortURLSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&usecaseStub{
		saveOriginalURLFn: func(ctx context.Context, originalURL string) (*model.URL, error) {
			return &model.URL{
				ID:        uuid.New(),
				ShortCode: "Ab3_kLm901",
				Original:  originalURL,
			}, nil
		},
		getOriginalURLFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
	})

	router := gin.New()
	handler.RegisterRoutesGin(router)

	body := bytes.NewBufferString(`{"url":"https://example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", body)
	req.Header.Set("Content-Type", "application/json")
	req.Host = "localhost:8080"

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp createShortURLResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.ShortURL != "http://localhost:8080/Ab3_kLm901" {
		t.Fatalf("unexpected short url %q", resp.ShortURL)
	}
}

func TestCreateShortURLInvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&usecaseStub{
		saveOriginalURLFn: func(ctx context.Context, originalURL string) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
		getOriginalURLFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
	})

	router := gin.New()
	handler.RegisterRoutesGin(router)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", bytes.NewBufferString(`{"url":`))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestCreateShortURLUsecaseValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&usecaseStub{
		saveOriginalURLFn: func(ctx context.Context, originalURL string) (*model.URL, error) {
			return nil, model.ErrInvalidOriginalURL
		},
		getOriginalURLFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
	})

	router := gin.New()
	handler.RegisterRoutesGin(router)

	body := bytes.NewBufferString(`{"url":"not-a-url"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/shorten", body)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestGetOriginalURLSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&usecaseStub{
		saveOriginalURLFn: func(ctx context.Context, originalURL string) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
		getOriginalURLFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			return &model.URL{
				ID:        uuid.New(),
				ShortCode: shortCode,
				Original:  "https://example.com",
			}, nil
		},
	})

	router := gin.New()
	handler.RegisterRoutesGin(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/Ab3_kLm901", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp getOriginalURLResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.OriginalURL != "https://example.com" {
		t.Fatalf("unexpected original url %q", resp.OriginalURL)
	}
}

func TestGetOriginalURLNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&usecaseStub{
		saveOriginalURLFn: func(ctx context.Context, originalURL string) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
		getOriginalURLFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			return nil, model.ErrURLNotFound
		},
	})

	router := gin.New()
	handler.RegisterRoutesGin(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/Ab3_kLm901", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestGetOriginalURLInvalidShortCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&usecaseStub{
		saveOriginalURLFn: func(ctx context.Context, originalURL string) (*model.URL, error) {
			return nil, errors.New("unexpected call")
		},
		getOriginalURLFn: func(ctx context.Context, shortCode string) (*model.URL, error) {
			return nil, model.ErrInvalidShortCodeLength
		},
	})

	router := gin.New()
	handler.RegisterRoutesGin(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/short", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
