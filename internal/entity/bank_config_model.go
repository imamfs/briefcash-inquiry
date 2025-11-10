package entity

type BankConfig struct {
	BankCode           string `gorm:"column:bank_code"`
	BankName           string `gorm:"column:bank_name"`
	ExternalInquiryURL string `gorm:"column:internal_inquiry_url"`
	InternalInquiryURL string `gorm:"column:external_inquiry_url"`
	AccessTokenURL     string `gorm:"column:access_token_url"`
	BaseURL            string `gorm:"column:base_url"`
	ClientKey          string `gorm:"column:client_key"`
	ClientSecret       string `gorm:"column:client_secret"`
	PartnerId          string `gorm:"column:partner_id"`
	ChannelId          string `gorm:"column:channel_id"`
}
