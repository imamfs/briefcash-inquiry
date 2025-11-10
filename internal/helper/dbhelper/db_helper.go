package dbhelper

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBHelper struct {
	DB *gorm.DB
}

type DBConfig struct {
	Host     string
	Port     string
	DBName   string
	Username string
	Password string
	SSLMode  string
}

func NewDBHelper(cfg DBConfig) (*DBHelper, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDb, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("failed to get generic database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	if err := sqlDb.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping to database: %w", err)
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetConnMaxLifetime(time.Hour)

	return &DBHelper{DB: db}, nil
}

func (h *DBHelper) Close() error {
	sqlDb, err := h.DB.DB()

	if err == nil {
		return sqlDb.Close()
	}

	return fmt.Errorf("failed to close database: %w", err)
}
