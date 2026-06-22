package service

import (
	"backend/internal/models"
	"context"
)

type Link interface {
	Create(ctx context.Context, originalLink models.OriginalLink) (*models.ShortCode, error)
	Get(ctx context.Context, code models.ShortCode) (*models.OriginalLink, error)
}
