package urls

import (
	"context"
	"time"

	repoUrls "github.com/cantylv/thumbnail-loader/internal/repository/urls"
)

type Usecase interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Init() error
}

type UsecaseLayer struct {
	repo repoUrls.Repo
}

func NewUsecaseLayer(repo repoUrls.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repo: repo,
	}
}

func (r *UsecaseLayer) Set(ctx context.Context, key, value string) error {
	return r.repo.Save(ctx, key, value)
}

func (r *UsecaseLayer) Get(ctx context.Context, key string) (string, error) {
	return r.repo.Get(ctx, key)
}

func (r *UsecaseLayer) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.repo.Init(ctx)
}
