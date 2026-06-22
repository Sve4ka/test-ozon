package link

import (
	"context"
	"errors"
	"testing"

	"backend/internal/cerr"
	"backend/internal/mocks"
	"backend/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Create_CacheHit(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")
	shortCode := models.ShortCode("aB12_cdEF3")

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(shortCode, nil).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Create(ctx, originalLink)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, shortCode, *got)
}

func TestService_Create_NewLink(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")
	generatedCode := "aB12_cdEF3"
	shortCode := models.ShortCode(generatedCode)

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(models.ShortCode(""), cerr.ErrNotFound).
		Once()

	generator.
		On("Generate").
		Return(generatedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, shortCode).
		Return(nil).
		Once()

	cache.
		On("Set", mock.Anything, originalLink, shortCode).
		Return(nil).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Create(ctx, originalLink)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, shortCode, *got)
}

func TestService_Create_CacheSetError(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")
	generatedCode := "aB12_cdEF3"
	shortCode := models.ShortCode(generatedCode)
	cacheErr := errors.New("cache set failed")

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(models.ShortCode(""), cerr.ErrNotFound).
		Once()

	generator.
		On("Generate").
		Return(generatedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, shortCode).
		Return(nil).
		Once()

	cache.
		On("Set", mock.Anything, originalLink, shortCode).
		Return(cacheErr).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Create(ctx, originalLink)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, shortCode, *got)
}

func TestService_Create_OriginalAlreadyExists(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")
	generatedCode := "aB12_cdEF3"
	generatedShortCode := models.ShortCode(generatedCode)
	existingShortCode := models.ShortCode("aB12_cdEF4")

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(models.ShortCode(""), cerr.ErrNotFound).
		Once()

	generator.
		On("Generate").
		Return(generatedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, generatedShortCode).
		Return(cerr.ErrOriginalLinkAlreadyExists).
		Once()

	repo.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(existingShortCode, nil).
		Once()

	cache.
		On("Set", mock.Anything, originalLink, existingShortCode).
		Return(nil).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Create(ctx, originalLink)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, existingShortCode, *got)
}

func TestService_Create_ShortCodeAlreadyExists(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")

	busyGeneratedCode := "aB12_cdEF3"
	busyShortCode := models.ShortCode(busyGeneratedCode)

	freeGeneratedCode := "aB12_cdEF4"
	freeShortCode := models.ShortCode(freeGeneratedCode)

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(models.ShortCode(""), cerr.ErrNotFound).
		Once()

	generator.
		On("Generate").
		Return(busyGeneratedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, busyShortCode).
		Return(cerr.ErrShortCodeAlreadyExists).
		Once()

	generator.
		On("Generate").
		Return(freeGeneratedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, freeShortCode).
		Return(nil).
		Once()

	cache.
		On("Set", mock.Anything, originalLink, freeShortCode).
		Return(nil).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Create(ctx, originalLink)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, freeShortCode, *got)
}

func TestService_Create_GeneratorError(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")
	generatorErr := errors.New("generator failed")

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(models.ShortCode(""), cerr.ErrNotFound).
		Once()

	generator.
		On("Generate").
		Return("", generatorErr).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Create(ctx, originalLink)

	require.Nil(t, got)
	require.ErrorIs(t, err, generatorErr)
}

func TestService_Create_StorageUnexpectedError(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")
	generatedCode := "aB12_cdEF3"
	shortCode := models.ShortCode(generatedCode)
	storageErr := errors.New("storage is unavailable")

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(models.ShortCode(""), cerr.ErrNotFound).
		Once()

	generator.
		On("Generate").
		Return(generatedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, shortCode).
		Return(storageErr).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Create(ctx, originalLink)

	require.Nil(t, got)
	require.ErrorIs(t, err, storageErr)
}

func TestService_Create_AttemptsExceeded(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	originalLink := models.OriginalLink("https://example.com")

	firstGeneratedCode := "aB12_cdEF3"
	firstShortCode := models.ShortCode(firstGeneratedCode)

	secondGeneratedCode := "aB12_cdEF4"
	secondShortCode := models.ShortCode(secondGeneratedCode)

	cache.
		On("GetByOriginalLink", mock.Anything, originalLink).
		Return(models.ShortCode(""), cerr.ErrNotFound).
		Once()

	generator.
		On("Generate").
		Return(firstGeneratedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, firstShortCode).
		Return(cerr.ErrShortCodeAlreadyExists).
		Once()

	generator.
		On("Generate").
		Return(secondGeneratedCode, nil).
		Once()

	repo.
		On("Create", mock.Anything, originalLink, secondShortCode).
		Return(cerr.ErrShortCodeAlreadyExists).
		Once()

	service := InitLinkService(repo, cache, generator, 2)

	got, err := service.Create(ctx, originalLink)

	require.Nil(t, got)
	require.ErrorIs(t, err, cerr.ErrNotGenerated)
}

func TestService_Get_CacheHit(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	shortCode := models.ShortCode("aB12_cdEF3")
	originalLink := models.OriginalLink("https://example.com")

	cache.
		On("GetByShortCode", mock.Anything, shortCode).
		Return(originalLink, nil).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Get(ctx, shortCode)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, originalLink, *got)
}

func TestService_Get_StorageHit(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	shortCode := models.ShortCode("aB12_cdEF3")
	originalLink := models.OriginalLink("https://example.com")

	cache.
		On("GetByShortCode", mock.Anything, shortCode).
		Return(models.OriginalLink(""), cerr.ErrNotFound).
		Once()

	repo.
		On("GetByShortCode", mock.Anything, shortCode).
		Return(originalLink, nil).
		Once()

	cache.
		On("Set", mock.Anything, originalLink, shortCode).
		Return(nil).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Get(ctx, shortCode)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, originalLink, *got)
}

func TestService_Get_CacheSetError(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	shortCode := models.ShortCode("aB12_cdEF3")
	originalLink := models.OriginalLink("https://example.com")
	cacheErr := errors.New("cache set failed")

	cache.
		On("GetByShortCode", mock.Anything, shortCode).
		Return(models.OriginalLink(""), cerr.ErrNotFound).
		Once()

	repo.
		On("GetByShortCode", mock.Anything, shortCode).
		Return(originalLink, nil).
		Once()

	cache.
		On("Set", mock.Anything, originalLink, shortCode).
		Return(cacheErr).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Get(ctx, shortCode)

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, originalLink, *got)
}

func TestService_Get_StorageNotFound(t *testing.T) {
	ctx := context.Background()

	repo := mocks.NewStorage(t)
	cache := mocks.NewCache(t)
	generator := mocks.NewGenerator(t)

	shortCode := models.ShortCode("aB12_cdEF3")

	cache.
		On("GetByShortCode", mock.Anything, shortCode).
		Return(models.OriginalLink(""), cerr.ErrNotFound).
		Once()

	repo.
		On("GetByShortCode", mock.Anything, shortCode).
		Return(models.OriginalLink(""), cerr.ErrNotFound).
		Once()

	service := InitLinkService(repo, cache, generator, 10)

	got, err := service.Get(ctx, shortCode)

	require.Nil(t, got)
	require.ErrorIs(t, err, cerr.ErrNotFound)
}
