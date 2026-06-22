package link

import (
	"backend/internal/cache"
	"backend/internal/cerr"
	"backend/internal/generate"
	"backend/internal/log"
	"backend/internal/models"
	"backend/internal/service"
	"backend/internal/storage"
	"context"
	"errors"
)

type Service struct {
	repo            storage.Storage
	cache           cache.Cache
	generate        generate.Generator
	generateAttempt int
}

func InitLinkService(repo storage.Storage, cache cache.Cache, generate generate.Generator, generateAttempt int) service.Link {
	return Service{cache: cache, repo: repo, generate: generate, generateAttempt: generateAttempt}
}

func (s Service) Create(ctx context.Context, originalLink models.OriginalLink) (*models.ShortCode, error) {
	if code, err := s.cache.GetByOriginalLink(ctx, originalLink); err == nil {

		return &code, nil
	}

	for i := 0; i < s.generateAttempt; i++ {
		code, err := s.generate.Generate()
		if err != nil {
			log.Log.Error(err)
			return nil, err
		}

		shortCode := models.ShortCode(code)

		err = s.repo.Create(ctx, originalLink, shortCode)
		if err == nil {

			if err = s.cache.Set(ctx, originalLink, shortCode); err != nil {
				log.Log.Error(err)
			}

			return &shortCode, nil
		}

		if errors.Is(err, cerr.ErrOriginalLinkAlreadyExists) {
			shortCode, err = s.repo.GetByOriginalLink(ctx, originalLink)
			if err != nil {
				log.Log.Error(err)

				return nil, err
			}

			if err = s.cache.Set(ctx, originalLink, shortCode); err != nil {
				log.Log.Error(err)
			}

			return &shortCode, nil
		}

		if errors.Is(err, cerr.ErrShortCodeAlreadyExists) {
			continue
		}

		log.Log.Error(err)

		return nil, err
	}

	log.Log.Error(cerr.ErrNotGenerated)

	return nil, cerr.ErrNotGenerated
}

func (s Service) Get(ctx context.Context, code models.ShortCode) (*models.OriginalLink, error) {
	if originalLink, err := s.cache.GetByShortCode(ctx, code); err == nil {
		return &originalLink, nil
	}

	originalLink, err := s.repo.GetByShortCode(ctx, code)

	if err != nil {
		log.Log.Error(err)
		return nil, err
	}

	if err = s.cache.Set(ctx, originalLink, code); err != nil {
		log.Log.Error(err)
	}

	return &originalLink, nil
}
