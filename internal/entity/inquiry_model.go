package entity

import "time"

type Inquiry struct {
	ID                     int64     `gorm:"column:id;primaryKey;autoIncrement"`
	MerchantCode           string    `gorm:"column:merchant_code"`
	PartnerReferenceNo     string    `gorm:"column:partner_reference_no"`
	BeneficiaryAccount     string    `gorm:"column:beneficiary_account"`
	BeneficiaryBankCode    string    `gorm:"column:beneficiary_bank_code"`
	BeneficiaryAccountName string    `gorm:"column:beneficiary_account_name"`
	InquiryDate            time.Time `gorm:"column:inquiry_date"`
	Status                 string    `gorm:"column:status"`
}
