package cache

import (
	"backend/internal/models"
	"context"
)

type Cache interface {
	Set(ctx context.Context, link models.OriginalLink, code models.ShortCode) error
	GetByOriginalLink(ctx context.Context, link models.OriginalLink) (models.ShortCode, error)
	GetByShortCode(ctx context.Context, code models.ShortCode) (models.OriginalLink, error)
}
