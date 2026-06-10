package v1

import (
	"github.com/gin-gonic/gin"
)

func (h *HandlerFacade) RegisterRoutesGin(handler *gin.Engine) {
	api := handler.Group("/api/v1")
	{
		api.POST("/shorten", h.CreateShortURL)
		api.GET("/:shortCode", h.GetOriginalURL)
	}

}
