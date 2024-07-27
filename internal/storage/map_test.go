//go:build !coverage

package storage

import (
	"context"
	"os"
	"testing"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_mapStorage(t *testing.T) {
	testMapStorage := NewInMemoryStorage()
	ctx := context.Background()
	const testFilePath = "test.json"

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
		testGaugeMetricValue2 = 20.2
		testGaugeMetricID2    = "test_gauge_metric_2"
		testGaugeMetric2      = model.Metrics{
			ID:    testGaugeMetricID2,
			MType: model.MetricTypeGauge,
			Delta: nil,
			Value: &testGaugeMetricValue2,
		}
		testCounterMetricValue2 = int64(20)
		testCounterMetric2      = model.Metrics{
			ID:    testCounterMetricID1,
			MType: model.MetricTypeCounter,
			Delta: &testCounterMetricValue2,
			Value: nil,
		}
	)

	t.Run("update val", func(t *testing.T) {
		errUpdateMetricValue := testMapStorage.UpdateMetricValue(ctx, testGaugeMetric1)
		assert.NoError(t, errUpdateMetricValue)
		errUpdateMetricValue2 := testMapStorage.UpdateMetricValue(ctx, testCounterMetric1)
		assert.NoError(t, errUpdateMetricValue2)
	})

	t.Run("update batch", func(t *testing.T) {
		errUpdateMetrics := testMapStorage.UpdateMetrics(ctx, []model.Metrics{testCounterMetric2, testGaugeMetric2})
		assert.NoError(t, errUpdateMetrics)
	})
	f, err := os.CreateTemp("", testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name()) //nolint:all
	t.Run("store", func(t *testing.T) {
		errUpdateMetrics := testMapStorage.StoreMetrics(ctx, f.Name())
		assert.NoError(t, errUpdateMetrics)
	})
	testMapStorage2 := NewInMemoryStorage()
	t.Run("restore", func(t *testing.T) {
		errUpdateMetrics := testMapStorage2.RestoreMetrics(ctx, f.Name())
		assert.NoError(t, errUpdateMetrics)
	})
	t.Run("get metric", func(t *testing.T) {
		metric1, errGetMetricByModel1 := testMapStorage2.GetMetricByModel(ctx, model.Metrics{ID: testCounterMetricID1, MType: model.MetricTypeCounter})
		assert.NoError(t, errGetMetricByModel1)
		assert.Equal(t, *testCounterMetric1.Delta+*testCounterMetric2.Delta, *metric1.Delta)
		metric2, errGetMetricByModel2 := testMapStorage2.GetMetricByModel(ctx, model.Metrics{ID: testGaugeMetricID2, MType: model.MetricTypeGauge})
		assert.NoError(t, errGetMetricByModel2)
		assert.Equal(t, testGaugeMetric2, metric2)
	})

	t.Run("get all metrics", func(t *testing.T) {
		_, errGetMetricByModel1 := testMapStorage2.GetAllMetrics(ctx)
		assert.NoError(t, errGetMetricByModel1)

	})

}
