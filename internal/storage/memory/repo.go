package memory

import (
	"context"
	"url_shortener/internal/cerr"
	"url_shortener/internal/models"
	"url_shortener/internal/storage"
)

type Repo struct {
	storage *Storage
}

func InitMemoryRepo(storage *Storage) storage.Storage {
	return &Repo{storage: storage}
}

func (r *Repo) Create(ctx context.Context, link models.OriginalLink, code models.ShortCode) error {
	r.storage.mu.Lock()
	defer r.storage.mu.Unlock()

	if _, exists := r.storage.byOriginal[link]; exists {
		return cerr.ErrOriginalLinkAlreadyExists
	}

	if _, exists := r.storage.byCode[code]; exists {
		return cerr.ErrShortCodeAlreadyExists
	}

	r.storage.byOriginal[link] = code
	r.storage.byCode[code] = link

	return nil
}

func (r *Repo) GetByOriginalLink(ctx context.Context, link models.OriginalLink) (models.ShortCode, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	code, ok := r.storage.byOriginal[link]
	if !ok {
		return "", cerr.ErrNotFound
	}
	return code, nil
}

func (r *Repo) GetByShortCode(ctx context.Context, code models.ShortCode) (models.OriginalLink, error) {
	r.storage.mu.RLock()
	defer r.storage.mu.RUnlock()

	link, ok := r.storage.byCode[code]
	if !ok {
		return "", cerr.ErrNotFound
	}
	return link, nil
}
