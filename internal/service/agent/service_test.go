package service

import (
	"log"
	"testing"
	"time"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/metrics"
	"github.com/mrkovshik/yametrics/internal/request"
	storage2 "github.com/mrkovshik/yametrics/internal/storage"
	"go.uber.org/zap"
)

func TestAgent_Poll(t *testing.T) {
	src := metrics.NewMockMetrics()
	strg := storage2.NewInMemoryStorage()
	cfg, _ := config.GetConfigs()
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("zap.NewDevelopment", zap.Error(err))
	}

	// Flushes buffered log entries before program exits
	defer logger.Sync() //nolint:all
	sugar := logger.Sugar()
	client := request.NewRestClient(sugar, &cfg)
	a := NewAgent(src, &cfg, strg, sugar, client)

	t.Run("util_metrics", func(t *testing.T) {
		ch := make(chan time.Time)
		done := make(chan struct{}, 1)
		go func(ch chan time.Time) {
			ch <- time.Now()
			close(ch)
		}(ch)
		a.PollUtilMetrics(ch, done)
		<-done

	})
	t.Run("metrics", func(t *testing.T) {
		done := make(chan struct{}, 1)
		ch := make(chan time.Time)
		go func(ch chan time.Time) {
			ch <- time.Now()
			close(ch)
		}(ch)
		a.PollMetrics(ch, done)
		<-done
	})
}
