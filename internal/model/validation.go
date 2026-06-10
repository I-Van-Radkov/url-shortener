package model

import (
	"net/url"
	"unicode/utf8"
)

func ValidateOriginalURL(raw string) error {
	if raw == "" {
		return ErrEmptyOriginalURL
	}

	parsed, err := url.ParseRequestURI(raw)
	if err != nil {
		return ErrInvalidOriginalURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidOriginalURL
	}

	if parsed.Hostname() == "" {
		return ErrInvalidOriginalURL
	}

	return nil
}

func ValidateShortCode(code string) error {
	if utf8.RuneCountInString(code) != ShortCodeLength {
		return ErrInvalidShortCodeLength
	}

	for _, ch := range code {
		if ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			continue
		}

		return ErrInvalidShortCode
	}

	return nil
}
