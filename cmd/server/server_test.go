package main

import (
	"bytes"
	"encoding/json"
	"github.com/mrkovshik/yametrics/internal/model"
	"go.uber.org/zap"

	"net/http"
	"testing"
	"time"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
	"github.com/mrkovshik/yametrics/internal/storage/server"
	"github.com/stretchr/testify/require"
)

func Test_server(t *testing.T) {
	var (
		testCounter1      = int64(3)
		testCounterResult = int64(3) * 2
		testGauge1        = 2.5
		testGauge2        = 3.5
	)
	type (
		want struct {
			code        int
			response    model.Metrics
			contentType string
		}
		request struct {
			method      string
			url         string
			contentType string
			req         model.Metrics
		}
	)
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "positive update #1",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test1",
					MType: model.MetricTypeGauge,
					Value: &testGauge1,
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive update #2",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test2",
					MType: model.MetricTypeCounter,
					Delta: &testCounter1,
				},
			},

			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},

		{
			name: "positive get #1",
			request: request{
				method:      http.MethodGet,
				url:         "http://localhost:8080/value/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test2",
					MType: model.MetricTypeCounter,
				},
			},

			want: want{
				code: http.StatusOK,
				response: model.Metrics{
					ID:    "test2",
					MType: model.MetricTypeCounter,
					Delta: &testCounter1,
				},
				contentType: "application/json",
			},
		},

		{
			name: "positive get #2",
			request: request{
				method:      http.MethodGet,
				url:         "http://localhost:8080/value/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test1",
					MType: model.MetricTypeGauge,
				},
			},

			want: want{
				code: http.StatusOK,
				response: model.Metrics{
					ID:    "test1",
					MType: model.MetricTypeGauge,
					Value: &testGauge1,
				},
				contentType: "application/json",
			},
		},
		{
			name: "positive update #3",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test2",
					MType: model.MetricTypeCounter,
					Delta: &testCounter1,
				},
			},

			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive update #4",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test1",
					MType: model.MetricTypeGauge,
					Value: &testGauge2,
				},
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "positive get #3",
			request: request{
				method:      http.MethodGet,
				url:         "http://localhost:8080/value/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test2",
					MType: model.MetricTypeCounter,
				},
			},

			want: want{
				code: http.StatusOK,
				response: model.Metrics{
					ID:    "test2",
					MType: model.MetricTypeCounter,
					Delta: &testCounterResult,
				},
				contentType: "application/json",
			},
		},

		{
			name: "positive get #4",
			request: request{
				method:      http.MethodGet,
				url:         "http://localhost:8080/value/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test1",
					MType: model.MetricTypeGauge,
				},
			},

			want: want{
				code: http.StatusOK,
				response: model.Metrics{
					ID:    "test1",
					MType: model.MetricTypeGauge,
					Value: &testGauge2,
				},
				contentType: "application/json",
			},
		},
		{
			name: "negative update #1",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test1",
					MType: model.MetricTypeCounter,
					Value: &testGauge1,
				},
			},

			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative update #2",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "test1",
					MType: "non_existing_type",
					Value: &testGauge1,
				},
			},

			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "negative get #1",
			request: request{
				method:      http.MethodGet,
				url:         "http://localhost:8080/value/",
				contentType: "application/json",
				req: model.Metrics{
					ID:    "non_existing_name",
					MType: model.MetricTypeGauge,
					Value: &testGauge1,
				},
			},

			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	mapStorage := storage.NewMapStorage()
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("zap.NewDevelopment",
			zap.Error(err))
	}
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()
	buf := bytes.Buffer{}

	cfg, err2 := config.GetConfigs()
	require.NoError(t, err2)
	getMetricsService := service.NewServer(mapStorage, cfg, sugar)
	go run(getMetricsService)
	time.Sleep(1 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err22 := json.NewEncoder(&buf).Encode(tt.request.req)
			require.NoError(t, err22)
			client := &http.Client{}
			req, err3 := http.NewRequest(tt.request.method, tt.request.url, &buf)
			require.NoError(t, err3)
			response, err4 := client.Do(req)
			require.NoError(t, err4)
			require.Equal(t, tt.want.code, response.StatusCode)
			contentType := response.Header.Get("Content-Type")
			require.Equal(t, tt.want.contentType, contentType)
			if contentType == "application/json" {
				respBody := model.Metrics{}
				err44 := json.NewDecoder(response.Body).Decode(&respBody)
				require.NoError(t, err44)
				require.Equal(t, tt.want.response.MType, respBody.MType)
				require.Equal(t, tt.want.response.ID, respBody.ID)
				if tt.want.response.MType == model.MetricTypeCounter {
					require.Equal(t, *tt.want.response.Delta, *respBody.Delta)
				}
				if tt.want.response.MType == model.MetricTypeGauge {
					require.Equal(t, *tt.want.response.Value, *respBody.Value)
				}
			}
			err5 := response.Body.Close()
			require.NoError(t, err5)
		})
	}
}
