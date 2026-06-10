package v1

import (
	"context"
	"errors"
	"net/http"

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
	var req createShortURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	url, err := h.usecase.SaveOriginalURL(c.Request.Context(), req.URL)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, createShortURLResponse{
		ShortURL: url.ShortCode,
	})
}

func (h *HandlerFacade) GetOriginalURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := h.usecase.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, getOriginalURLResponse{
		OriginalURL: url.Original,
	})
}

func buildShortURL(c *gin.Context, shortCode string) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + c.Request.Host + "/" + shortCode
}

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, model.ErrEmptyOriginalURL),
		errors.Is(err, model.ErrInvalidOriginalURL),
		errors.Is(err, model.ErrInvalidShortCode),
		errors.Is(err, model.ErrInvalidShortCodeLength):
		c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})

	case errors.Is(err, model.ErrURLNotFound):
		c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})

	case errors.Is(err, model.ErrGenerationFailed):
		c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})

	default:
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal server error"})
	}
}
