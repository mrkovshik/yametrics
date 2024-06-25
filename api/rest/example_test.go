package rest

import (
	"context"
	"fmt"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/model"
	service "github.com/mrkovshik/yametrics/internal/service/server"
	"github.com/mrkovshik/yametrics/internal/storage"
	"go.uber.org/zap"
)

func Example() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	cfg, err := config.GetConfigs()
	if err != nil {
		sugar.Fatal("cfg.GetConfigs", err)
	}
	ctx := context.Background()
	metricStorage := storage.NewMapStorage()
	metricService := service.NewMetricService(metricStorage, &cfg, sugar)
	gauge := 2.5
	counter := int64(2)
	metric1 := model.Metrics{
		ID:    "test_gauge",
		MType: model.MetricTypeGauge,
		Value: &gauge,
	}
	metric2 := model.Metrics{
		ID:    "test_counter",
		MType: model.MetricTypeCounter,
		Delta: &counter,
	}
	if err := metricService.UpdateMetrics(ctx, []model.Metrics{
		metric1,
		metric2,
	}); err != nil {
		sugar.Fatal("metricService.UpdateMetrics", err)
	}

	m1, err := metricService.GetMetric(ctx, model.Metrics{
		ID:    "test_counter",
		MType: model.MetricTypeCounter,
	})

	m2, err := metricService.GetMetric(ctx, model.Metrics{
		ID:    "test_gauge",
		MType: model.MetricTypeGauge,
	})

	if err != nil {
		sugar.Fatal("metricService.UpdateMetrics", err)
	}
	fmt.Println(*m1.Delta, m1.ID, m1.MType)
	fmt.Println(*m2.Value, m2.ID, m2.MType)

	// Output:
	// 2 test_counter counter
	// 2.5 test_gauge gauge
}
