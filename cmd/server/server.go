package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/mrkovshik/yametrics/api/rest"
	"github.com/mrkovshik/yametrics/internal/storage"
	"github.com/mrkovshik/yametrics/internal/util/retriable"
	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
)

func main() {
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
	metricStorage := storage.NewMapStorage()
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
		metricStorage = storage.NewDBStorage(db)
	}
	ctx := context.Background()
	metricService := service.NewMetricService(metricStorage, &cfg, sugar)
	apiService := rest.NewRestAPIServer(metricService, &cfg, sugar)
	if cfg.RestoreEnable {
		if err := metricStorage.RestoreMetrics(ctx, cfg.StoreFilePath); err != nil {
			sugar.Fatal("RestoreMetrics", err)
		}
	}

	if cfg.StoreEnable && !cfg.SyncStoreEnable {
		storeTicker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
		go func() {
			for range storeTicker.C {
				if err := metricStorage.StoreMetrics(ctx, cfg.StoreFilePath); err != nil {
					sugar.Fatal("StoreMetrics", err)
				}
			}
		}()
	}

	apiService.RunServer(ctx)
	if err := metricStorage.StoreMetrics(ctx, cfg.StoreFilePath); err != nil {
		sugar.Fatal("StoreMetrics", err)
	}
}
