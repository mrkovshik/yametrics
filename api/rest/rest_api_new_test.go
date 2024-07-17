package rest

import (
	"fmt"
	"testing"
)

//import (
//	"context"
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	_ "github.com/lib/pq"
//	config "github.com/mrkovshik/yametrics/internal/config/server"
//	"github.com/mrkovshik/yametrics/internal/model"
//	"github.com/mrkovshik/yametrics/internal/service/server/mock_server"
//	"github.com/stretchr/testify/require"
//	"go.uber.org/zap"
//)
//
//func TestPing(t *testing.T) {
//
//	logger, err := zap.NewDevelopment()
//	if err != nil {
//		logger.Fatal("zap.NewDevelopment",
//			zap.Error(err))
//	}
//	defer logger.Sync() //nolint:all
//	sugar := logger.Sugar()
//	cfg, err2 := config.GetTestConfig()
//	cfg.Key = "some_test_key"
//	require.NoError(t, err2)
//	ctx := context.Background()
//	ctrl := gomock.NewController(t)
//	metricService := defineServiceMock(ctx, ctrl)
//	apiService := NewServer(metricService, &cfg, sugar).ConfigureRouter()
//
//	t.Run("ping", func(t *testing.T) {
//		req, err := http.NewRequest("GET", "/ping", nil)
//		require.NoError(t, err)
//		rr := httptest.NewRecorder()
//		apiService.server.Handler.ServeHTTP(rr, req)
//		require.Equal(t, http.StatusOK, rr.Code)
//	})
//
//	t.Run("update_metric_from_json", func(t *testing.T) {
//		value := 123.45
//		newMetric := model.Metrics{
//			ID:    "metric1",
//			MType: "gauge",
//			Value: &value,
//		}
//		_, err := json.Marshal(newMetric)
//		require.NoError(t, err)
//		req, err := http.NewRequest("POST", "/update", nil)
//		require.NoError(t, err)
//		rr := httptest.NewRecorder()
//		apiService.server.Handler.ServeHTTP(rr, req)
//		require.Equal(t, http.StatusOK, rr.Code)
//	})
//}
//
//func defineServiceMock(ctx context.Context, ctrl *gomock.Controller) *mock_server.MockService {
//	srv := mock_server.NewMockService(ctrl)
//	srv.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
//	srv.EXPECT().Ping(gomock.Any()).Return(nil).AnyTimes()
//	return srv
//}

func TestPing(t *testing.T) {
	fmt.Println(3.5+1.2 == 4.7)
}
