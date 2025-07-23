package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type URLRepo struct {
	db *pgxpool.Pool
}

func NewURLRepository(db *pgxpool.Pool) *URLRepo {
	return &URLRepo{db: db}
}

func (u *URLRepo) Save(ctx context.Context, originalURL, alias string) (string, error) {
	const query = `INSERT INTO urls(original_url, alias) VALUES ($1, $2) ON CONFLICT (original_url) DO UPDATE SET alias = urls.alias RETURNING alias`
	var aliasFromDB string
	err := u.db.QueryRow(ctx, query, originalURL, alias).Scan(&aliasFromDB)
	if err != nil {
		return "", err
	}
	return aliasFromDB, nil
}

func (u *URLRepo) FindByCode(ctx context.Context, alias string) (string, error) {
	const query = `SELECT original_url FROM urls WHERE alias = $1`
	var originalURL string
	err := u.db.QueryRow(ctx, query, alias).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf(`url with this short_code "%s" not found`, alias)
		}
		return "", err
	}
	return originalURL, nil
}

func (u *URLRepo) IncrementClicks(ctx context.Context, alias string) error {
	const query = `UPDATE urls SET click_count = click_count + 1 WHERE alias = $1`
	_, err := u.db.Exec(ctx, query, alias)
	return err
}

func (u *URLRepo) GetStats(ctx context.Context, alias string) (int, time.Time, string, error) {
	const query = `SELECT click_count, created_at, original_url FROM urls WHERE alias = $1`
	var clicks int
	var createdAt time.Time
	var originalURL string
	err := u.db.QueryRow(ctx, query, alias).Scan(&clicks, &createdAt, &originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, time.Time{}, "", fmt.Errorf(`url with this short_code "%s" not found`, alias)
		}
		return 0, time.Time{}, "", err
	}
	return clicks, createdAt, originalURL, nil
}
