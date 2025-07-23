package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type URLRepository interface {
	Save(ctx context.Context, originalURL, alias string) (string, error)
	FindByCode(ctx context.Context, alias string) (string, error)
	IncrementClicks(ctx context.Context, alias string) error
	GetStats(ctx context.Context, alias string) (int, time.Time, string, error)
}

type URLService struct {
	repo URLRepository
}

func NewURLService(repo URLRepository) *URLService {
	return &URLService{repo: repo}
}

func (u *URLService) Shorten(ctx context.Context, originalURL string) (string, error) {
	alias, err := generateAlias()
	if err != nil {
		return "", err
	}
	aliasFromDB, err := u.repo.Save(ctx, originalURL, alias)
	if err != nil {
		return "", fmt.Errorf("failed to save url: %w", err)
	}
	if aliasFromDB != alias {
		return aliasFromDB, nil
	}
	return alias, nil
}

func (u *URLService) Resolve(ctx context.Context, alias string) (string, error) {
	originalURL, err := u.repo.FindByCode(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("failed to find url: %w", err)
	}
	err = u.repo.IncrementClicks(ctx, alias)
	if err != nil {
		return "", fmt.Errorf("failed to increment clicks: %w", err)
	}
	return originalURL, nil
}

func (u *URLService) Stats(ctx context.Context, alias string) (int, time.Time, string, error) {
	clicks, createdAt, originalURL, err := u.repo.GetStats(ctx, alias)
	if err != nil {
		return 0, time.Time{}, "", fmt.Errorf("failed to get stats: %w", err)
	}
	return clicks, createdAt, originalURL, nil
}

func generateAlias() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random alias: %s", err)
	}
	return hex.EncodeToString(b), nil
}
