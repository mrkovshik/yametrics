package main

import (
	"github.com/mrkovshik/yametrics/internal/service"
	"github.com/mrkovshik/yametrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"testing"
)

func Test_server(t *testing.T) {
	type (
		want struct {
			code              int
			response          string
			contentType       string
			storageGaugeVal   float64
			storageCounterVal int64
		}
		request struct {
			method      string
			url         string
			contentType string
		}
	)
	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "positive test #1",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/gauge/Alloc/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusOK,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   456.00,
				storageCounterVal: 0,
			},
		},
		{
			name: "positive test #2",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/counter/PollCount/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusOK,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   456.00,
				storageCounterVal: 456,
			},
		},

		{
			name: "positive test #3",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/counter/PollCount/123",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusOK,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   456.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "positive test #4",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/gauge/Alloc/123",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusOK,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},

		{
			name: "negative test #1 (no counter name)",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/counter/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusNotFound,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "negative test #2 (invalid counter name)",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/counter/wrongname/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusNotFound,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "negative test #3 (invalid counter value)",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/counter/PollCount/45q6",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusBadRequest,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "negative test #4 (no gauge name)",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/gauge/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusNotFound,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "negative test #5 (invalid gauge name)",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/gauge/wrongname/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusNotFound,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "negative test #6 (invalid gauge value)",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/gauge/Alloc/45q6",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusBadRequest,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "negative test #7 (invalid metric type)",
			request: request{
				method:      http.MethodPost,
				url:         "http://localhost:8080/update/wrongtype/PollCount/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusBadRequest,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
		{
			name: "negative test #8 (invalid http method)",
			request: request{
				method:      http.MethodGet,
				url:         "http://localhost:8080/update/gauge/Alloc/456",
				contentType: "text/plain",
			},
			want: want{
				code:              http.StatusMethodNotAllowed,
				response:          "",
				contentType:       "text/plain; charset=utf-8",
				storageGaugeVal:   123.00,
				storageCounterVal: 579,
			},
		},
	}

	mapStorage := storage.NewMapStorage()

	getMetricsService := service.NewServiceWithMapStorage(mapStorage, log.Default())
	go run(getMetricsService)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest(tt.request.method, tt.request.url, nil)
			require.NoError(t, err)
			response, err1 := client.Do(req)
			require.NoError(t, err1)
			defer response.Body.Close()
			require.Equal(t, tt.want.code, response.StatusCode)
			require.Equal(t, tt.want.contentType, response.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.storageGaugeVal, mapStorage.Gauges["Alloc"])
			assert.Equal(t, tt.want.storageCounterVal, mapStorage.Counters["PollCount"])
		})
	}
}
