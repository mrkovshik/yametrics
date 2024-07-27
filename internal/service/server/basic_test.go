package server

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	config "github.com/mrkovshik/yametrics/internal/config/server"
	"github.com/mrkovshik/yametrics/internal/model"
	mock_storage "github.com/mrkovshik/yametrics/internal/storage/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	testGaugeID1   = "gauge_1"
	testCounterID1 = "counter_1"
)

var (
	loggerConfig = zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Encoding:    "json", // You can use "console" for a more readable format
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "", // Disable caller key to remove caller information
			MessageKey:     "message",
			StacktraceKey:  "", // Disable stacktrace key to remove stack traces
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	testGaugeVal1     = 20.5
	testCounterDelta1 = int64(2)

	testGauge1 = model.Metrics{
		ID:    testGaugeID1,
		MType: model.MetricTypeGauge,
		Value: &testGaugeVal1,
	}
	testCounter1 = model.Metrics{
		ID:    testCounterID1,
		MType: model.MetricTypeCounter,
		Delta: &testCounterDelta1,
	}
)

func TestMetricService(t *testing.T) {
	cfg, errGetTestConfig := config.GetTestConfig()
	assert.NoError(t, errGetTestConfig)
	logger, errBuild := loggerConfig.Build()
	assert.NoError(t, errBuild)
	defer logger.Sync() //nolint:all

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := defineStorage(ctx, ctrl)
	basicSvs := NewMetricService(mockStorage, &cfg, logger.Sugar())

	t.Run("update", func(t *testing.T) {
		err := basicSvs.UpdateMetrics(ctx, []model.Metrics{testCounter1, testGauge1})
		assert.NoError(t, err)
	})

	t.Run("get_metric", func(t *testing.T) {
		m, err := basicSvs.GetMetric(ctx, model.Metrics{ID: testCounterID1})
		assert.NoError(t, err)
		assert.Equal(t, testCounter1, m)
	})

	t.Run("get_all", func(t *testing.T) {
		s, err := basicSvs.GetAllMetrics(ctx)
		assert.NoError(t, err)
		assert.NotEqual(t, "", s)
	})

	t.Run("store", func(t *testing.T) {
		err := basicSvs.StoreMetrics(ctx)
		assert.NoError(t, err)
	})

	t.Run("restore", func(t *testing.T) {
		err := basicSvs.RestoreMetrics(ctx)
		assert.NoError(t, err)
	})
}

func defineStorage(ctx context.Context, ctrl *gomock.Controller) *mock_storage.MockStorage {
	strg := mock_storage.NewMockStorage(ctrl)
	strg.EXPECT().UpdateMetrics(ctx, []model.Metrics{testCounter1, testGauge1}).Return(nil).AnyTimes()
	strg.EXPECT().GetMetricByModel(ctx, model.Metrics{ID: testCounterID1}).Return(testCounter1, nil).AnyTimes()
	strg.EXPECT().GetAllMetrics(ctx).Return(map[string]model.Metrics{testGauge1.ID: testGauge1, testCounter1.ID: testCounter1}, nil).AnyTimes()
	strg.EXPECT().StoreMetrics(ctx, "./tmp/metrics-test.json").Return(nil).AnyTimes()
	strg.EXPECT().RestoreMetrics(ctx, "./tmp/metrics-test.json").Return(nil).AnyTimes()
	return strg
}
