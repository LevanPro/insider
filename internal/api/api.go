package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LevanPro/insider/internal/config"
	"github.com/LevanPro/insider/internal/infra/database"
	"github.com/LevanPro/insider/internal/infra/logger"
	"github.com/LevanPro/insider/internal/infra/scheduler"
	"github.com/LevanPro/insider/internal/infra/sender"
	"github.com/LevanPro/insider/internal/repository"
	"github.com/LevanPro/insider/internal/service"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type App struct {
	log       *zap.SugaredLogger
	db        *sqlx.DB
	service   *service.MessageService
	scheduler *scheduler.Scheduler
}

func Run() error {
	log, err := logger.InitLogger("MESSAGE-SERVICE")
	if err != nil {
		return fmt.Errorf("unable to create logger %w", err)
	}

	defer log.Sync()

	log.Infow("startup", "time", time.Now().UTC())
	defer log.Infow("shutdown complete")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error loading config %w", err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// =========================================================================
	log.Infow("startup", "status", "initializing database support", cfg.DB.Host)

	cfgDB := database.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	}

	db, err := database.Open(cfgDB)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.DB.Host)
		db.Close()
	}()
	// =========================================================================

	// ========== Run migrations =========================================

	err = database.RunMigrations(db)
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	// ===================================================================
	postgresMessageRepo := repository.NewPostgresMessageRepository(db)
	senderClient := sender.NewClient(cfg.Application.WebhookURL, cfg.Application.WebhookAuthKey)

	messageService := service.NewMessageService(postgresMessageRepo, senderClient, cfg.Application.BatchSize, cfg.Application.NumberOfWorkers, log)
	scheduler := scheduler.NewScheduler(messageService.ProcessNextUnsent, cfg.Application.SchedulerInterval, cfg.Application.SchedulerStartImmediate)
	scheduler.Start()

	app := &App{
		db:        db,
		log:       log,
		scheduler: scheduler,
		service:   messageService,
	}

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      app.setupRoutes(),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
