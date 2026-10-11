package noop

import (
	"context"
	"url_shortener/internal/cache"
	"url_shortener/internal/cerr"
	"url_shortener/internal/models"
)

type Cache struct {
}

func InitNoopCache() cache.Cache {
	return &Cache{}
}

func (c Cache) Set(ctx context.Context, link models.OriginalLink, code models.ShortCode) error {
	return nil
}

func (c Cache) GetByOriginalLink(ctx context.Context, link models.OriginalLink) (models.ShortCode, error) {
	return "", cerr.ErrNotFound
}

func (c Cache) GetByShortCode(ctx context.Context, code models.ShortCode) (models.OriginalLink, error) {
	return "", cerr.ErrNotFound
}
