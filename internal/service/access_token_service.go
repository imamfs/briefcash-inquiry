package service

import (
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/repository"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TokenService interface {
	SaveAccessTokenDB(ctx context.Context, token *entity.AccessToken) error
	SaveAccessTokenRedis(ctx context.Context, bank string, token *entity.AccessToken) error
	GetActiveAccessToken(ctx context.Context, bank string) (string, error)
}

type tokenService struct {
	db         *gorm.DB
	tokenRepo  repository.TokenRepository
	tokenRedis repository.TokenRedisRepository
}

func NewTokenService(db *gorm.DB, tokenRepo repository.TokenRepository, tokenRedis repository.TokenRedisRepository) TokenService {
	return &tokenService{db, tokenRepo, tokenRedis}
}

func (s *tokenService) SaveAccessTokenDB(ctx context.Context, token *entity.AccessToken) error {
	// Save new access token to database
	if err := s.saveToken(ctx, token); err != nil {
		return fmt.Errorf("failed to save new access token to database: %w", err)
	}
	return nil
}

func (s *tokenService) SaveAccessTokenRedis(ctx context.Context, bank string, token *entity.AccessToken) error {
	// Save new access token to redis
	expiresIn := token.ExpiresIn
	if expiresIn <= 30 {
		expiresIn = 30
	}
	ttl := time.Duration(expiresIn-30) * time.Second

	key := fmt.Sprintf("%s:access_token", bank)
	if err := s.tokenRedis.SetToken(ctx, key, token.AccessToken, ttl); err != nil {
		return fmt.Errorf("failed to set new access token to redis")
	}
	return nil
}

func (s *tokenService) GetActiveAccessToken(ctx context.Context, bank string) (string, error) {
	// Check latest active access token in redis
	key := fmt.Sprintf("%s:access_token", bank)
	value, err := s.tokenRedis.GetToken(ctx, key)
	if err == nil {
		return value, nil
	}

	// fallback to database
	tokenEntity, err := s.tokenRepo.FindToken(ctx)
	if err != nil {
		return "", fmt.Errorf("token not found in redis and database: %w", err)
	}

	expiresIn := tokenEntity.ExpiresIn
	if expiresIn <= 30 {
		expiresIn = 30
	}

	ttl := time.Duration(expiresIn-30) * time.Second
	_ = s.tokenRedis.SetToken(ctx, key, tokenEntity.AccessToken, ttl)

	return tokenEntity.AccessToken, nil
}

func (s *tokenService) saveToken(ctx context.Context, token *entity.AccessToken) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		repo := s.tokenRepo.WithTransaction(tx)
		return repo.SaveToken(ctx, token)
	})
}
