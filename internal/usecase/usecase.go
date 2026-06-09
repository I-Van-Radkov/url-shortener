package usecase

import (
	"context"

	"github.com/I-Van-Radkov/url-shortener/internal/model"
)

type Repository interface {
	Save(ctx context.Context, url *model.URL) error
	FindByOriginalURL(ctx context.Context, originalURL string) (*model.URL, error)
	FindByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
}

type Usecase struct {
	repo Repository
}

func New(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) SaveOriginalURL(ctx context.Context, originalUrl string) (*model.URL, error) {
	// TODO: валидация

	// TODO: проверка на существование ссылки, если существует - возвращаем короткий код

	// TODO: генерация короткого кода

	// TODO: проверка на существование короткого кода, если существует - перегенерируем

	// TODO: сохраняем ссылку

	// TODO: возвращаем результат
}

func (u *Usecase) GetOriginalURL(ctx context.Context, shortCode string) (*model.URL, error) {
	// TODO: валидация

	// TODO: получение ссылки

	// TODO: если ссылка не найдена, то возвращаем ошибку

	// TODO: возвращаем результат
}
