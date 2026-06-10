package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewRepo(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresRepository) CreateOrFind(ctx context.Context, urlInput *model.URL) (*model.URL, error) {
	if urlInput == nil {
		return nil, fmt.Errorf("url input is nil")
	}

	query := `
        INSERT INTO urls (id, short_code, original, created_at)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (original) DO NOTHING
        RETURNING id, short_code, original, created_at
    `

	var url model.URL
	err := r.db.QueryRow(ctx, query, urlInput.ID, urlInput.ShortCode, urlInput.Original, urlInput.CreatedAt).Scan(&url.ID, &url.ShortCode, &url.Original, &url.CreatedAt)
	if err == nil {
		return &url, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return r.findByOriginalURL(ctx, urlInput.Original)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "urls_short_code_key":
			return nil, model.ErrShortCodeExists
		}
	}

	return nil, fmt.Errorf("failed to create or find: %w", err)
}

func (r *PostgresRepository) FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
	query, args, err := r.builder.
		Select("id", "short_code", "original", "created_at").
		From("urls").
		Where(squirrel.Eq{"short_code": shortCode}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var url model.URL
	err = r.db.QueryRow(ctx, query, args...).Scan(&url.ID, &url.ShortCode, &url.Original, &url.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrURLNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &url, nil
}

func (r *PostgresRepository) findByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	query, args, err := r.builder.
		Select("id", "short_code", "original", "created_at").
		From("urls").
		Where(squirrel.Eq{"original": originalURL}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var url model.URL
	err = r.db.QueryRow(ctx, query, args...).Scan(&url.ID, &url.ShortCode, &url.Original, &url.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrURLNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &url, nil
}
