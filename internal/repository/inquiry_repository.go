package repository

import (
	"context"
	"fmt"

	model "briefcash-inquiry/internal/entity"

	"gorm.io/gorm"
)

type InquiryRepository interface {
	SaveInquiry(ctx context.Context, inquiry *model.Inquiry) error
	WithTransaction(trx *gorm.DB) InquiryRepository
}

type inquiryRepository struct {
	db *gorm.DB
}

func NewInquiryRepository(db *gorm.DB) InquiryRepository {
	return &inquiryRepository{db}
}

func (ir *inquiryRepository) SaveInquiry(ctx context.Context, inquiry *model.Inquiry) error {
	err := ir.db.WithContext(ctx).Table("inquiry").Create(inquiry).Error
	if err != nil {
		return fmt.Errorf("failed to save inquiry data to database: %w", err)
	}
	return nil
}

func (ir *inquiryRepository) WithTransaction(trx *gorm.DB) InquiryRepository {
	return &inquiryRepository{db: trx}
}
