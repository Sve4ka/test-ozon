package memory

import (
	"context"
	"testing"

	"backend/internal/cerr"
	"backend/internal/models"

	"github.com/stretchr/testify/require"
)

func TestRepo_CreateAndGet(t *testing.T) {
	repo := InitMemoryRepo(NewStorage())

	original := models.OriginalLink("https://example.com")
	code := models.ShortCode("aB12_cdEF3")

	err := repo.Create(context.Background(), original, code)
	require.NoError(t, err)

	gotCode, err := repo.GetByOriginalLink(context.Background(), original)
	require.NoError(t, err)
	require.Equal(t, code, gotCode)

	gotLink, err := repo.GetByShortCode(context.Background(), code)
	require.NoError(t, err)
	require.Equal(t, original, gotLink)
}

func TestRepo_Create_DuplicateOriginalLink(t *testing.T) {
	repo := InitMemoryRepo(NewStorage())

	original := models.OriginalLink("https://example.com")

	err := repo.Create(context.Background(), original, models.ShortCode("aB12_cdEF3"))
	require.NoError(t, err)

	err = repo.Create(context.Background(), original, models.ShortCode("ZZZ12_cdEF"))
	require.ErrorIs(t, err, cerr.ErrOriginalLinkAlreadyExists, "expected ErrOriginalLinkAlreadyExists, got %v", err)
}

func TestRepo_Create_DuplicateShortCode(t *testing.T) {
	repo := InitMemoryRepo(NewStorage())

	code := models.ShortCode("aB12_cdEF3")

	err := repo.Create(context.Background(), models.OriginalLink("https://example.com/1"), code)
	require.NoError(t, err)

	err = repo.Create(context.Background(), models.OriginalLink("https://example.com/2"), code)
	require.ErrorIs(t, err, cerr.ErrShortCodeAlreadyExists, "expected ErrShortCodeAlreadyExists, got %v", err)
}

func TestRepo_GetByOriginalLink_NotFound(t *testing.T) {
	repo := InitMemoryRepo(NewStorage())

	_, err := repo.GetByOriginalLink(context.Background(), models.OriginalLink("https://example.com"))
	require.ErrorIs(t, err, cerr.ErrNotFound, "expected ErrNotFound, got %v", err)
}

func TestRepo_GetByShortCode_NotFound(t *testing.T) {
	repo := InitMemoryRepo(NewStorage())

	_, err := repo.GetByShortCode(context.Background(), models.ShortCode("aB12_cdEF3"))
	require.ErrorIs(t, err, cerr.ErrNotFound, "expected ErrNotFound, got %v", err)
}
