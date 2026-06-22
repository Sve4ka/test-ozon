package memory

import (
	"backend/internal/cerr"
	"backend/internal/models"
	"backend/internal/storage"
	"context"
)

type Repo struct {
	db *Storage
}

func InitMemoryRepo(db *Storage) storage.Storage {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, link models.OriginalLink, code models.ShortCode) error {
	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, exists := r.db.byOriginal[link]; exists {
		return cerr.ErrOriginalLinkAlreadyExists
	}

	if _, exists := r.db.byCode[code]; exists {
		return cerr.ErrShortCodeAlreadyExists
	}

	r.db.byOriginal[link] = code
	r.db.byCode[code] = link

	return nil
}

func (r *Repo) GetByOriginalLink(ctx context.Context, link models.OriginalLink) (models.ShortCode, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	code, ok := r.db.byOriginal[link]
	if !ok {
		return "", cerr.ErrNotFound
	}
	return code, nil
}

func (r *Repo) GetByShortCode(ctx context.Context, code models.ShortCode) (models.OriginalLink, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	link, ok := r.db.byCode[code]
	if !ok {
		return "", cerr.ErrNotFound
	}
	return link, nil
}
