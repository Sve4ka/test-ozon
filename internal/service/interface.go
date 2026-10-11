package service

import (
	"context"
	"url_shortener/internal/models"
)

type Link interface {
	Create(ctx context.Context, originalLink models.OriginalLink) (*models.ShortCode, error)
	Get(ctx context.Context, code models.ShortCode) (*models.OriginalLink, error)
}
