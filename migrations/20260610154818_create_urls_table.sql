-- +goose Up
CREATE TABLE urls (
    id UUID PRIMARY KEY,
    short_code VARCHAR(10) NOT NULL,
    original TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT urls_short_code_key UNIQUE (short_code),
    CONSTRAINT urls_original_key UNIQUE (original)
);

-- +goose Down
DROP TABLE urls;
