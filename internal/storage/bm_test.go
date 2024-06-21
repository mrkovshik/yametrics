package storage

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/mrkovshik/yametrics/internal/model"
	"go.uber.org/zap"
)

func BenchmarkUpdateMetricValue(b *testing.B) {
	var (
		testGaugeMetricValue1 = 10.1
		testGaugeMetricID1    = "test_gauge_metric_1"
		testGaugeMetric1      = model.Metrics{
			ID:    testGaugeMetricID1,
			MType: model.MetricTypeGauge,
			Delta: nil,
			Value: &testGaugeMetricValue1,
		}
		testCounterMetricID1    = "test_counter_metric"
		testCounterMetricValue1 = int64(10)
		testCounterMetric1      = model.Metrics{
			ID:    testCounterMetricID1,
			MType: model.MetricTypeCounter,
			Delta: &testCounterMetricValue1,
			Value: nil,
		}
	)
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()

	db, err := sql.Open("postgres", "host=localhost port=5432 user=yandex password=yandex dbname=yandex sslmode=disable")
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
	_, err = db.Exec(ddl)

	if err != nil {
		sugar.Fatal("Exec", err)
	}
	ctx := context.Background()
	defer db.Close() //nolint:all
	postgresStorage := NewDBStorage(db)
	runtimeStorage := NewMapStorage()
	if err := postgresStorage.RestoreMetrics(ctx, "/metrics-db_bm.json"); err != nil {
		sugar.Fatal("RestoreMetrics", err)
	}

	b.Run("postgresStorage update metric", func(b *testing.B) {
		metrics := []model.Metrics{testGaugeMetric1, testCounterMetric1}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := postgresStorage.UpdateMetrics(ctx, metrics); err != nil {
				sugar.Fatal("UpdateMetrics", err)
			}
		}
	})

	b.Run("runtimeStorage update metric", func(b *testing.B) {
		metrics := []model.Metrics{testGaugeMetric1, testCounterMetric1}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := runtimeStorage.UpdateMetrics(ctx, metrics); err != nil {
				sugar.Fatal("UpdateMetrics", err)
			}
		}
	})

	b.Run("postgresStorage get metrics", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := postgresStorage.GetAllMetrics(ctx); err != nil {
				sugar.Fatal("UpdateMetrics", err)
			}
		}
	})

	b.Run("runtimeStorage get metrics", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := runtimeStorage.GetAllMetrics(ctx); err != nil {
				sugar.Fatal("UpdateMetrics", err)
			}
		}
	})

	b.Run("postgresStorage get metric by model", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := postgresStorage.GetMetricByModel(ctx, model.Metrics{
				ID:    testGaugeMetricID1,
				MType: model.MetricTypeGauge,
			}); err != nil {
				sugar.Fatal("UpdateMetrics", err)
			}
		}
	})

	b.Run("runtimeStorage get metric by model", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := runtimeStorage.GetMetricByModel(ctx, model.Metrics{
				ID:    testGaugeMetricID1,
				MType: model.MetricTypeGauge,
			}); err != nil {
				sugar.Fatal("UpdateMetrics", err)
			}
		}
	})

}
