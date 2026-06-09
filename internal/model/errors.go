package model

import "errors"

var (
	ErrURLNotFound            = errors.New("url not found")
	ErrShortCodeExists        = errors.New("short code already exists")
	ErrOriginalURLExists      = errors.New("original url already has a short link")
	ErrEmptyOriginalURL       = errors.New("original url cannot be empty")
	ErrInvalidOriginalURL     = errors.New("original url is invalid")
	ErrInvalidShortCodeLength = errors.New("short code must be exactly 10 characters")
	ErrGenerationFailed       = errors.New("failed to generate unique short code")
	ErrInvalidStorageType     = errors.New("invalid storage type")
)
