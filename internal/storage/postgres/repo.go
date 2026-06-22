package postgres

import (
	"backend/internal/cerr"
	"backend/internal/models"
	"backend/internal/storage"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repo struct {
	db *Pg
}

func InitPostgresRepo(db *Pg) storage.Storage {
	return &Repo{db: db}
}

func (r *Repo) Create(ctx context.Context, link models.OriginalLink, code models.ShortCode) error {
	var id int
	query := `INSERT INTO links (original_link, short_code) VALUES ($1, $2) returning id`
	err := r.db.Pool.QueryRow(ctx, query, link, code).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if !errors.As(err, &pgErr) {
			return err
		}
		if pgErr.Code == "23505" { // unique_violation
			switch pgErr.ConstraintName {
			case "links_original_link_key":
				return cerr.ErrOriginalLinkAlreadyExists
			case "links_short_code_key":
				return cerr.ErrShortCodeAlreadyExists
			default:
				return err
			}
		}
		return err
	}
	return nil
}

func (r *Repo) GetByOriginalLink(ctx context.Context, link models.OriginalLink) (models.ShortCode, error) {
	var code models.ShortCode
	query := `SELECT short_code FROM links WHERE original_link = $1`
	err := r.db.Pool.QueryRow(ctx, query, link).Scan(&code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", cerr.ErrNotFound
		}
		return "", err
	}
	return code, nil
}

func (r *Repo) GetByShortCode(ctx context.Context, code models.ShortCode) (models.OriginalLink, error) {
	var link models.OriginalLink
	query := `SELECT original_link FROM links WHERE short_code = $1`
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(&link)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", cerr.ErrNotFound
		}
		return "", err
	}
	return link, nil
}
