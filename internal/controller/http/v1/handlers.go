package v1

import (
	"context"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
	"github.com/gin-gonic/gin"
)

type Usecase interface {
	SaveOriginalURL(ctx context.Context, originalURL string) (*model.URL, error)
	GetOriginalURL(ctx context.Context, shortCode string) (*model.URL, error)
}

type HandlerFacade struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) *HandlerFacade {
	return &HandlerFacade{
		usecase: usecase,
	}
}

func (h *HandlerFacade) CreateShortURL(c *gin.Context) {
}

func (h *HandlerFacade) GetOriginalURL(c *gin.Context) {
}
