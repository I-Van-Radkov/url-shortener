package v1

type createShortURLRequest struct {
	URL string `json:"url" binding:"required"`
}

type createShortURLResponse struct {
	ShortURL string `json:"short_url"`
}

type getOriginalURLResponse struct {
	OriginalURL string `json:"original_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}
