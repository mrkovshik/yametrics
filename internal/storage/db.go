// Package storage provides implementations of the service.Storage interface for metrics storage.
package storage

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/util/retriable"
)

// DBStorage implements the service.Storage interface using a SQL database.
type DBStorage struct {
	db *sql.DB
}

// NewDBStorage creates a new instance of dBStorage with the provided SQL database connection.
func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{
		db: db,
	}
}

// UpdateMetricValue updates a single metric value in the database transactionally.
func (s *DBStorage) UpdateMetricValue(ctx context.Context, newMetrics model.Metrics) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:all
	err1 := s.updateMetricValue(ctx, newMetrics, tx)
	if err1 != nil {
		return err1
	}
	return tx.Commit()
}

// UpdateMetrics updates multiple metrics in the database transactionally.
func (s *DBStorage) UpdateMetrics(ctx context.Context, newMetrics []model.Metrics) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:all
	for _, metric := range newMetrics {
		err := s.updateMetricValue(ctx, metric, tx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// GetMetricByModel retrieves a metric from the database based on the provided model.
func (s *DBStorage) GetMetricByModel(ctx context.Context, newMetrics model.Metrics) (model.Metrics, error) {
	query := `SELECT id, type, value, delta FROM metrics WHERE id = $1`
	row, err := retriable.QueryRowRetryable(func() *sql.Row {
		return s.db.QueryRowContext(ctx, query, newMetrics.ID)
	})
	if err != nil {
		return model.Metrics{}, err
	}

	var foundMetric model.Metrics
	if err := row.Scan(&foundMetric.ID, &foundMetric.MType, &foundMetric.Value, &foundMetric.Delta); err != nil {
		return model.Metrics{}, err
	}
	return foundMetric, nil
}

// GetAllMetrics retrieves all metrics from the storage and returns them as a map
func (s *DBStorage) GetAllMetrics(ctx context.Context) (map[string]model.Metrics, error) {
	metricMap := make(map[string]model.Metrics)
	query := `SELECT id, type, value, delta FROM metrics`
	rows, err := retriable.QueryRetryable(func() (*sql.Rows, error) {
		return s.db.QueryContext(ctx, query)
	})
	if err != nil {
		return map[string]model.Metrics{}, err
	}
	defer rows.Close() //nolint:all
	for rows.Next() {
		currentMetric := model.Metrics{}
		if err := rows.Scan(&currentMetric.ID, &currentMetric.MType, &currentMetric.Value, &currentMetric.Delta); err != nil {
			return map[string]model.Metrics{}, err
		}
		metricMap[currentMetric.ID] = currentMetric
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return metricMap, nil
}

// StoreMetrics stores all metrics into a JSON file at the specified path.
func (s *DBStorage) StoreMetrics(ctx context.Context, path string) error {
	file, err := retriable.OpenRetryable(func() (*os.File, error) {
		return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	})
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
	metricMap, err := s.scanAllMetricsToMap(ctx)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(metricMap)
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	return err
}

// RestoreMetrics restores metrics from a JSON file at the specified path into the database.
func (s *DBStorage) RestoreMetrics(ctx context.Context, path string) error {
	file, err := retriable.OpenRetryable(func() (*os.File, error) {
		return os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	})
	if err != nil {
		return err
	}
	defer file.Close() //nolint:all
	reader := bufio.NewReader(file)
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	metricMap := make(map[string]model.Metrics)
	if err := json.Unmarshal(data, &metricMap); err != nil {
		return err
	}
	for _, value := range metricMap {
		if err := s.UpdateMetricValue(ctx, value); err != nil {
			return err
		}
	}
	return nil
}

// Ping pings the database to check the connectivity.
func (s *DBStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// scanAllMetricsToMap scans all metrics from the database and returns them as a map.
func (s *DBStorage) scanAllMetricsToMap(ctx context.Context) (map[string]model.Metrics, error) {
	metricMap := make(map[string]model.Metrics)
	query := `SELECT id, type, value, delta FROM metrics`
	rows, err := retriable.QueryRetryable(func() (*sql.Rows, error) {
		return s.db.QueryContext(ctx, query)
	})
	if err != nil {
		return map[string]model.Metrics{}, err
	}
	defer rows.Close() //nolint:all
	for rows.Next() {
		currentMetric := model.Metrics{}
		if err := rows.Scan(&currentMetric.ID, &currentMetric.MType, &currentMetric.Value, &currentMetric.Delta); err != nil {
			return map[string]model.Metrics{}, err
		}
		metricMap[currentMetric.ID] = currentMetric
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return metricMap, nil
}

// updateMetricValue updates the metric value in the database transactionally.
func (s *DBStorage) updateMetricValue(ctx context.Context, newMetrics model.Metrics, tx *sql.Tx) error {
	query := `SELECT id, type, value, delta FROM metrics WHERE id=$1 AND type= $2`
	row, err := retriable.QueryRowRetryable(func() *sql.Row {
		return tx.QueryRowContext(ctx, query, newMetrics.ID, newMetrics.MType)
	})
	if err != nil {
		return err
	}
	var (
		id, mType string
		value     sql.NullFloat64
		delta     sql.NullInt64
	)
	if errScan := row.Scan(&id, &mType, &value, &delta); errScan != nil {
		if errors.Is(errScan, sql.ErrNoRows) {
			query = `INSERT INTO metrics (id, type, value, delta)
		VALUES ($1, $2, $3, $4)`

			if errExecRetryable := retriable.ExecRetryable(func() error {
				_, errExecContext := tx.ExecContext(ctx, query, newMetrics.ID, newMetrics.MType, newMetrics.Value, newMetrics.Delta)
				return errExecContext
			}); errExecRetryable != nil {
				return errExecRetryable
			}
			return nil
		}
		return err
	}

	if mType == model.MetricTypeCounter {
		query = `UPDATE metrics value SET delta = $1 WHERE id = $2 AND type = $3`
		if !delta.Valid {
			return errors.New("unexpected null in delta field")
		}

		if errExecRetryable := retriable.ExecRetryable(func() error {
			_, errExecContext := tx.ExecContext(ctx, query, *newMetrics.Delta+delta.Int64, id, mType)
			return errExecContext
		}); errExecRetryable != nil {
			return errExecRetryable
		}
		return nil
	}
	query = `UPDATE metrics value SET value = $1 WHERE id = $2 AND type = $3`
	if !value.Valid {
		return errors.New("unexpected null in value field")
	}
	if errExecRetryable := retriable.ExecRetryable(func() error {
		_, errExecContext := tx.ExecContext(ctx, query, newMetrics.Value, id, mType)
		return errExecContext
	}); errExecRetryable != nil {
		return errExecRetryable
	}
	return nil
}
