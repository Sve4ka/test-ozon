package redis

import (
	"backend/internal/cache"
	"backend/internal/cerr"
	"backend/internal/models"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	db       *redis.Client
	cacheTTL time.Duration
}

func InitRedisCache(db *redis.Client, cacheTTL time.Duration) cache.Cache {
	return &Cache{db: db, cacheTTL: cacheTTL}
}
func originalKey(link models.OriginalLink) string {
	return "shortener:original:" + string(link)
}

func shortKey(code models.ShortCode) string {
	return "shortener:short:" + string(code)
}

func (c Cache) Set(ctx context.Context, link models.OriginalLink, code models.ShortCode) error {
	pipe := c.db.TxPipeline()
	pipe.Set(ctx, originalKey(link), string(code), c.cacheTTL)
	pipe.Set(ctx, shortKey(code), string(link), c.cacheTTL)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c Cache) GetByOriginalLink(ctx context.Context, link models.OriginalLink) (models.ShortCode, error) {
	code, err := c.db.GetEx(ctx, originalKey(link), c.cacheTTL).Result()
	if errors.Is(err, redis.Nil) {
		return "", cerr.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return models.ShortCode(code), nil
}

func (c Cache) GetByShortCode(ctx context.Context, code models.ShortCode) (models.OriginalLink, error) {
	link, err := c.db.GetEx(ctx, shortKey(code), c.cacheTTL).Result()

	if errors.Is(err, redis.Nil) {
		return "", cerr.ErrNotFound
	}
	if err != nil {
		return "", err
	}

	return models.OriginalLink(link), nil
}
