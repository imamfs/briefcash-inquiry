package entity

import "time"

type AccessToken struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	AccessToken string    `gorm:"column:access_token"`
	TokenType   string    `gorm:"column:token_type"`
	ExpiresIn   int16     `gorm:"column:expires_in"`
	ExpiresDate time.Time `gorm:"column:expires_date"`
}
