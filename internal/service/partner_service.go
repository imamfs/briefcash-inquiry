package service

import (
	"briefcash-inquiry/internal/entity"
	"briefcash-inquiry/internal/helper/loghelper"
	"briefcash-inquiry/internal/repository"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type BankPartner interface {
	LoadAllBankPartner(ctx context.Context) error
	GetBankConfig(bankCode string) entity.BankConfig
}

type bankPartner struct {
	mu        sync.RWMutex
	dbRepo    repository.PartnerRepository
	bankCache map[string]entity.BankConfig
}

func NewPartnerService(dbRepo repository.PartnerRepository) BankPartner {
	return &bankPartner{
		dbRepo:    dbRepo,
		bankCache: make(map[string]entity.BankConfig),
	}
}

func (s *bankPartner) LoadAllBankPartner(ctx context.Context) error {
	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":   "partner_service",
		"operation": "load_bank_partner_config",
	})

	log.WithField("step", "get_data_db").Info("Get existing bank route from db")
	banks, err := s.dbRepo.FindAll(ctx)
	if err != nil {
		log.WithField("step", "get_data_db").WithError(err).Error("Failed to fetch bank config from database")
		return err
	}

	log.WithField("step", "caching_config").Infof("Cache data bank config to memory, with total data %d", len(banks))
	s.mu.Lock()
	for _, bank := range banks {
		s.bankCache[bank.BankCode] = bank
	}
	s.mu.Unlock()
	return nil
}

func (s *bankPartner) GetBankConfig(bankCode string) entity.BankConfig {
	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":   "partner_service",
		"operation": "load_bank_partner_config",
	})

	s.mu.RLock()
	defer s.mu.RUnlock()

	log.WithField("step", "get_bank_config").Info("Get specific bank config from memory")
	bank, ok := s.bankCache[bankCode]

	if !ok {
		log.WithField("step", "get_bank_config").
			Warnf("Bank config not found for %s, fallback to default Bank BCA", bankCode)
		return s.bankCache["014"]
	}

	log.WithField("step", "get_bank_config").Infof("Bank %s is selected", bank.BankName)
	return bank
}
