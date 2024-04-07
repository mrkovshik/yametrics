package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/mrkovshik/yametrics/internal/model"
	service2 "github.com/mrkovshik/yametrics/internal/service/agent"
	"github.com/mrkovshik/yametrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	config "github.com/mrkovshik/yametrics/internal/config/server"
	service "github.com/mrkovshik/yametrics/internal/service/server"
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
			method        string
			url           string
			contentType   string
			contentEncode string
			//acceptEncoding string
			req model.Metrics
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
				method:        http.MethodPost,
				url:           "http://localhost:8080/update/",
				contentType:   "application/json",
				contentEncode: "gzip",
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
				method:        http.MethodPost,
				url:           "http://localhost:8080/update/",
				contentType:   "application/json",
				contentEncode: "gzip",
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
				method:      http.MethodPost,
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
				method:      http.MethodPost,
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
				url:         "http://localhost:8080/update/counter/test2/3",
				contentType: "text/plain; charset=utf-8",
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
				url:         "http://localhost:8080/update/gauge/test1/3.5",
				contentType: "text/plain; charset=utf-8",
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
				url:         "http://localhost:8080/value/counter/test2",
				contentType: "text/plain; charset=utf-8",
			},

			want: want{
				code: http.StatusOK,
				response: model.Metrics{
					ID:    "test2",
					MType: model.MetricTypeCounter,
					Delta: &testCounterResult,
				},
				contentType: "text/plain; charset=utf-8",
			},
		},

		{
			name: "positive get #4",
			request: request{
				method:      http.MethodGet,
				url:         "http://localhost:8080/value/gauge/test1",
				contentType: "text/plain; charset=utf-8",
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
				contentType: "text/plain; charset=utf-8",
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
				method:      http.MethodPost,
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
	cfg, err2 := config.GetConfigs()
	require.NoError(t, err2)
	getMetricsService := service.NewServer(mapStorage, cfg, sugar, nil)
	go run(getMetricsService, sugar, cfg)
	time.Sleep(1 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			req := *service2.NewRequestBuilder().SetURL(tt.request.url).SetMethod(tt.request.method)
			if tt.request.contentType == "application/json" {
				req.AddJSONBody(tt.request.req)
			}
			if tt.request.contentEncode == "gzip" {
				req.Compress()
			}
			require.NoError(t, req.Err)
			response, err4 := client.Do(&req.R)
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
			} else {
				if tt.want.response.MType == model.MetricTypeCounter {
					body, err55 := io.ReadAll(response.Body)
					assert.NoError(t, err55)
					val, err66 := strconv.ParseInt(string(body), 10, 64)
					assert.NoError(t, err66)
					require.Equal(t, *tt.want.response.Delta, val)
				}
				if tt.want.response.MType == model.MetricTypeGauge {
					body, err55 := io.ReadAll(response.Body)
					assert.NoError(t, err55)
					val, err66 := strconv.ParseFloat(string(body), 64)
					assert.NoError(t, err66)
					require.Equal(t, *tt.want.response.Value, val)
				}
			}

			err5 := response.Body.Close()
			require.NoError(t, err5)
		})
	}
}
