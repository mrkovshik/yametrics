//go:build !coverage

package storage

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/stretchr/testify/assert"
)

func Test_DBStorage(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=yandex password=yandex dbname=yandex sslmode=disable")
	assert.NoError(t, err)
	ddl := `CREATE TABLE IF NOT EXISTS metrics  
		(
		    id    varchar not null
		constraint metrics_pk
		primary key,
			type  varchar not null,
			value double precision,
			delta BIGINT			
		);
TRUNCATE TABLE metrics;`
	_, err = db.Exec(ddl)
	assert.NoError(t, err)

	ctx := context.Background()

	//Удаляем все записи из таблицы
	defer db.Exec(`TRUNCATE TABLE metrics;`) //nolint:all
	defer db.Close()                         //nolint:all
	testDBStorage := NewDBStorage(db)
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
		errUpdateMetricValue := testDBStorage.UpdateMetricValue(ctx, testGaugeMetric1)
		assert.NoError(t, errUpdateMetricValue)
		errUpdateMetricValue2 := testDBStorage.UpdateMetricValue(ctx, testCounterMetric1)
		assert.NoError(t, errUpdateMetricValue2)
	})

	t.Run("update batch", func(t *testing.T) {
		errUpdateMetrics := testDBStorage.UpdateMetrics(ctx, []model.Metrics{testCounterMetric2, testGaugeMetric2})
		assert.NoError(t, errUpdateMetrics)
	})
	f, errCreate := os.CreateTemp("", testFilePath)
	assert.NoError(t, errCreate)
	defer os.Remove(f.Name()) //nolint:all
	t.Run("store", func(t *testing.T) {
		errUpdateMetrics := testDBStorage.StoreMetrics(ctx, f.Name())
		assert.NoError(t, errUpdateMetrics)
	})

	//Удаляем все записи из таблицы, чтобы проверить, как она будет восстанавливать данные из файла
	_, err = db.Exec(`TRUNCATE TABLE metrics;`) //nolint:all
	assert.NoError(t, err)

	t.Run("restore", func(t *testing.T) {
		errUpdateMetrics := testDBStorage.RestoreMetrics(ctx, f.Name())
		assert.NoError(t, errUpdateMetrics)
	})
	t.Run("get all metrics", func(t *testing.T) {
		_, errGetMetricByModel1 := testDBStorage.GetAllMetrics(ctx)
		assert.NoError(t, errGetMetricByModel1)

	})

	t.Run("get metric", func(t *testing.T) {
		metric1, errGetMetricByModel1 := testDBStorage.GetMetricByModel(ctx, model.Metrics{ID: testCounterMetricID1, MType: model.MetricTypeCounter})
		assert.NoError(t, errGetMetricByModel1)
		assert.Equal(t, *testCounterMetric1.Delta+*testCounterMetric2.Delta, *metric1.Delta)
		metric2, errGetMetricByModel2 := testDBStorage.GetMetricByModel(ctx, model.Metrics{ID: testGaugeMetricID2, MType: model.MetricTypeGauge})
		assert.NoError(t, errGetMetricByModel2)
		assert.Equal(t, testGaugeMetric2, metric2)
	})

}
