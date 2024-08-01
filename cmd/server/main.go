package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/mrkovshik/yametrics/api"
	"github.com/mrkovshik/yametrics/api/grpc"
	"github.com/mrkovshik/yametrics/internal/storage"
	"github.com/mrkovshik/yametrics/internal/util/retriable"
	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
)

var (
	buildVersion, buildDate, buildCommit string
)

func main() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)
	var metricService *service.MetricService
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()

	cfg, err := config.GetConfigs()
	if err != nil {
		sugar.Fatal("cfg.GetConfigs", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var db *sql.DB
	if cfg.DBEnable {
		db, err = sql.Open("postgres", cfg.DBAddress)
		if err != nil {
			sugar.Fatal("sql.Open", err)
		}
		ddl := `CREATE TABLE IF NOT EXISTS metrics  
		(
		    id    varchar not null
		constraint metrics_pk
		primary key,
			type  varchar not null,
			value double precision,
			delta BIGINT			
		);`

		if err := retriable.ExecRetryable(func() error {
			_, err := db.Exec(ddl)
			return err
		}); err != nil {
			sugar.Fatal("Exec", err)
		}

		defer db.Close() //nolint:all
		dbStorage := storage.NewPostgresStorage(db)
		metricService = service.NewMetricService(dbStorage, &cfg, sugar)
	} else {
		metricStorage := storage.NewInMemoryStorage()
		metricService = service.NewMetricService(metricStorage, &cfg, sugar)
	}
	apiService := grpc.NewServer(metricService, &cfg, sugar)

	if cfg.RestoreEnable {
		if err := metricService.RestoreMetrics(ctx); err != nil {
			sugar.Fatal("RestoreMetrics", err)
		}
	}

	if cfg.StoreEnable && !cfg.SyncStoreEnable {
		storeTicker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		go func() {
			for range storeTicker.C {
				if err := metricService.StoreMetrics(ctx); err != nil {
					sugar.Fatal("StoreMetrics", err)
				}
			}
		}()
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	run(stop, apiService)
	if err := metricService.StoreMetrics(context.Background()); err != nil {
		sugar.Fatal("StoreMetrics", err)
	}
}

func run(stop chan os.Signal, srv api.Server) {
	log.Fatal(srv.RunServer(stop))
}
