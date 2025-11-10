package repository

import (
	"briefcash-inquiry/internal/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PartnerRepository interface {
	FindAll(ctx context.Context) ([]entity.BankConfig, error)
}

type partnerRepository struct {
	db *gorm.DB
}

func NewPartnerRepository(db *gorm.DB) PartnerRepository {
	return &partnerRepository{db}
}

func (r *partnerRepository) FindAll(ctx context.Context) ([]entity.BankConfig, error) {
	var listConfig []entity.BankConfig

	err := r.db.WithContext(ctx).Table("partner").
		Select("partner.company_bank_code AS bank_code, domestic_bank.short_name AS bank_name, partner_settings.api_key AS client_key, partner_settings.api_secret AS client_secret, partner_settings.partner_id, partner_settings.channel_id, partner_url.internal_inquiry_url, partner_url.external_inquiry_url, partner_url.access_token_url, partner_url.base_url").
		Joins("INNER JOIN partner_url ON partner.company_id = partner_url.company_id").
		Joins("INNER JOIN domestic_bank ON partner.company_id = domestic_bank.company_id").
		Joins("INNER JOIN partner_settings ON partner.company_id = partner_settings.company_id").
		Scan(&listConfig).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return listConfig, nil
}
