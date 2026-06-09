package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
)

type Repository interface {
	FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
	CreateOrFind(ctx context.Context, url *model.URL) (*model.URL, error)
}

type Generator interface {
	Generate() (string, error)
}

type Usecase struct {
	repo Repository
	gen  Generator

	maxAttemptsToGen int
}

func New(repo Repository, gen Generator, maxAttemptsToGen int) *Usecase {
	if maxAttemptsToGen <= 0 {
		maxAttemptsToGen = 5
	}

	return &Usecase{
		repo:             repo,
		gen:              gen,
		maxAttemptsToGen: maxAttemptsToGen,
	}
}

func (u *Usecase) SaveOriginalURL(ctx context.Context, originalURL string) (*model.URL, error) {
	err := model.ValidateOriginalURL(originalURL)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	for i := 0; i < u.maxAttemptsToGen; i++ {
		code, err := u.gen.Generate()
		if err != nil {
			return nil, err
		}

		url := &model.URL{
			ShortCode: code,
			Original:  originalURL,
			CreatedAt: now,
		}

		url, err = u.repo.CreateOrFind(ctx, url)
		if errors.Is(err, model.ErrShortCodeExists) {
			continue
		}
		if err != nil {
			return nil, err
		}

		return url, nil
	}

	return nil, model.ErrGenerationFailed
}

func (u *Usecase) GetOriginalURL(ctx context.Context, shortCode string) (*model.URL, error) {
	err := model.ValidateShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	url, err := u.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	return url, nil
}
