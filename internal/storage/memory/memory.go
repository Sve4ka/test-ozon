package memory

import (
	"backend/internal/models"
	"sync"
)

type Storage struct {
	mu sync.RWMutex

	byOriginal map[models.OriginalLink]models.ShortCode
	byCode     map[models.ShortCode]models.OriginalLink
}

func NewStorage() *Storage {
	return &Storage{
		byOriginal: make(map[models.OriginalLink]models.ShortCode),
		byCode:     make(map[models.ShortCode]models.OriginalLink),
		mu:         sync.RWMutex{},
	}
}
