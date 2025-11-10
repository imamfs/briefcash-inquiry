package repository

import (
	"briefcash-inquiry/internal/entity"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type TokenRepository interface {
	SaveToken(ctx context.Context, token *entity.AccessToken) error
	FindLatestValidToken(ctx context.Context) (string, error)
	FindToken(ctx context.Context) (*entity.AccessToken, error)
	WithTransaction(trx *gorm.DB) TokenRepository
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) TokenRepository {
	return &tokenRepository{db}
}

func (r *tokenRepository) SaveToken(ctx context.Context, token *entity.AccessToken) error {
	err := r.db.WithContext(ctx).Table("access_token").Create(token).Error

	if err != nil {
		return fmt.Errorf("failed to save token to database %w", err)
	}

	return nil
}

func (r *tokenRepository) FindLatestValidToken(ctx context.Context) (string, error) {
	var accessToken string

	err := r.db.WithContext(ctx).Table("access_token").
		Select("access_token").
		Where("expires_date > NOW()").Order("expires_date DESC").
		Limit(1).Scan(&accessToken).Error

	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (r *tokenRepository) FindToken(ctx context.Context) (*entity.AccessToken, error) {
	var accessToken entity.AccessToken

	err := r.db.WithContext(ctx).Table("access_token").
		Where("expires_date > NOW()").Order("expires_date DESC").
		Limit(1).First(&accessToken).Error

	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (r *tokenRepository) WithTransaction(trx *gorm.DB) TokenRepository {
	return &tokenRepository{db: trx}
}
