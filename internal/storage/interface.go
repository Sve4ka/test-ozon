package storage

import (
	"context"
	"url_shortener/internal/models"
)

type Storage interface {
	Create(ctx context.Context, link models.OriginalLink, code models.ShortCode) error
	GetByOriginalLink(ctx context.Context, link models.OriginalLink) (models.ShortCode, error)
	GetByShortCode(ctx context.Context, code models.ShortCode) (models.OriginalLink, error)
}
