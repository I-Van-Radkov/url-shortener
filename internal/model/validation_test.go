package model

import "testing"

func TestValidateOriginalURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid https url",
			input:   "https://example.com",
			wantErr: nil,
		},
		{
			name:    "valid http url",
			input:   "http://example.com/path?q=1",
			wantErr: nil,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: ErrEmptyOriginalURL,
		},
		{
			name:    "missing scheme",
			input:   "example.com",
			wantErr: ErrInvalidOriginalURL,
		},
		{
			name:    "unsupported scheme",
			input:   "ftp://example.com",
			wantErr: ErrInvalidOriginalURL,
		},
		{
			name:    "missing host",
			input:   "https://",
			wantErr: ErrInvalidOriginalURL,
		},
		{
			name:    "garbage string",
			input:   "not-a-url",
			wantErr: ErrInvalidOriginalURL,
		},
		{
			name:    "localhost is valid",
			input:   "http://localhost:8080/test",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateOriginalURL(tt.input)
			if err != tt.wantErr {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateShortCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid lowercase",
			input:   "abcdefghij",
			wantErr: nil,
		},
		{
			name:    "valid mixed charset",
			input:   "Ab3_kLm901",
			wantErr: nil,
		},
		{
			name:    "too short",
			input:   "abc",
			wantErr: ErrInvalidShortCodeLength,
		},
		{
			name:    "too long",
			input:   "abcdefghijk",
			wantErr: ErrInvalidShortCodeLength,
		},
		{
			name:    "contains dash",
			input:   "abcde-fghi",
			wantErr: ErrInvalidShortCode,
		},
		{
			name:    "contains exclamation mark",
			input:   "abcde!fghi",
			wantErr: ErrInvalidShortCode,
		},
		{
			name:    "contains cyrillic",
			input:   "абвгдеёжзи",
			wantErr: ErrInvalidShortCode,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: ErrInvalidShortCodeLength,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateShortCode(tt.input)
			if err != tt.wantErr {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
