package main

import (
	"briefcash-inquiry/config"
	"briefcash-inquiry/internal/controller"
	"briefcash-inquiry/internal/helper/dbhelper"
	"briefcash-inquiry/internal/helper/loghelper"
	"briefcash-inquiry/internal/helper/redishelper"
	"briefcash-inquiry/internal/repository"
	"briefcash-inquiry/internal/service"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	loghelper.InitLogger("./resource/app.log", logrus.InfoLevel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// graceful shutdown
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
		sign := <-signalChannel
		loghelper.Logger.WithField("signal", sign.String()).Info("Received shutdown signal")
		cancel()
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to load configuration")
	}

	dbCfg := dbhelper.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		DBName:   cfg.DBName,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
		SSLMode:  "disable",
	}

	dbHelper, err := dbhelper.NewDBHelper(dbCfg)
	if err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to create connection to databases")
	}
	defer dbHelper.Close()

	redisClient, err := redishelper.NewRedisHelper(cfg)
	if err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to established connection with redis")
	}
	defer redisClient.Close()

	inquiryRepo := repository.NewInquiryRepository(dbHelper.DB)
	partnerRepo := repository.NewPartnerRepository(dbHelper.DB)
	tokenRepo := repository.NewTokenRepository(dbHelper.DB)
	tokenRedis := repository.NewTokenRedisRepository(redisClient.Client)
	partnerService := service.NewPartnerService(partnerRepo)

	if err := partnerService.LoadAllBankPartner(ctx); err != nil {
		loghelper.Logger.WithError(err).Fatal("Failed to load bank route config to memory")
	}

	tokenService := service.NewTokenService(dbHelper.DB, tokenRepo, tokenRedis)
	inquiryService := service.NewInquiryService(inquiryRepo, tokenService, partnerService, dbHelper.DB)
	inquiryController := controller.NewInquiryController(inquiryService)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(RequestLoggerMiddleware())

	api := router.Group("/api/v1/")
	api.POST("/inquiry", inquiryController.InquiryAccountNumber)

	server := &http.Server{
		Addr:    cfg.AppPort,
		Handler: router,
	}

	go func() {
		loghelper.Logger.WithField("port", cfg.AppPort).Info("Inquiry Account Service is running...")
		if err := router.Run(cfg.AppPort); err != nil {
			loghelper.Logger.WithError(err).Fatal("Failed to start Inquiry Account Service")
		}
	}()

	<-ctx.Done()

	loghelper.Logger.Info("Shutting down server properly...")

	// beri jeda waktu untuk shutdown
	shutDownCtx, cancelShutDown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutDown()

	if err := server.Shutdown(shutDownCtx); err != nil {
		loghelper.Logger.WithError(err).Error("Forced shutdown due to timeout")
	} else {
		loghelper.Logger.Info("Inquiry Account Service shutdown completed")
	}

}

func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		duration := time.Since(start)
		status := c.Writer.Status()

		loghelper.Logger.WithFields(logrus.Fields{
			"method":   c.Request.Method,
			"path":     c.FullPath(),
			"status":   status,
			"duration": duration.String(),
			"clientIp": c.ClientIP(),
		}).Info("Handled request")
	}
}
