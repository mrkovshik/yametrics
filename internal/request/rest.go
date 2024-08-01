package request

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/internal/reqbuilder"
	"go.uber.org/zap"
)

type RestClient struct {
	logger *zap.SugaredLogger  // Logger for logging messages
	cfg    *config.AgentConfig // Configuration for the agent
}

func NewRestClient(logger *zap.SugaredLogger, cfg *config.AgentConfig) *RestClient {
	return &RestClient{
		logger: logger,
		cfg:    cfg,
	}
}

// Request processes metrics and sends them to the server.
func (r *RestClient) Request(id int, jobs <-chan model.Metrics) {
	for j := range jobs {
		r.logger.Debugf("worker #%v is sending %v", id, j.ID)
		metricUpdateURL := fmt.Sprintf("http://%v/update/", r.cfg.Address)

		reqBuilder := reqbuilder.NewHTTPRequestBuilder().SetURL(metricUpdateURL).AddJSONBody(j).Sign(r.cfg.Key).EncryptRSA(r.cfg.CryptoKey).Compress().SetMethod(http.MethodPost).AddIPHeader()
		if reqBuilder.Err != nil {
			r.logger.Errorf("error building request: %v\n", reqBuilder.Err)
			return
		}
		response, err := r.retryableSend(&reqBuilder.R)
		if err != nil {
			r.logger.Errorf("error sending request: %v\n", err)
			return
		}
		if response.StatusCode != http.StatusOK {
			r.logger.Errorf("status code is %v\n", response.StatusCode)
			return
		}
		if err := response.Body.Close(); err != nil {
			r.logger.Error("response.Body.Close()", err)
			return
		}
	}
}

// retryableSend sends an HTTP request with retries.
func (r *RestClient) retryableSend(req *http.Request) (*http.Response, error) {
	var (
		bodyBytes      []byte
		retryIntervals = []int{1, 3, 5} //TODO: move to config
		client         = http.Client{Timeout: 5 * time.Second}
		err            error
	)
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		// Reset the request body for retries.
		req.Body.Close() //nolint:all
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	for i := 0; i <= len(retryIntervals); i++ {
		response, err := client.Do(req)
		if err == nil {
			return response, nil
		}
		if i == len(retryIntervals) {
			return nil, err
		}
		r.logger.Errorf("failed connect to server: %v\n retry in %v seconds\n", err, retryIntervals[i])
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
		if req.Body != nil {
			req.Body.Close() //nolint:all
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}
	}
	return nil, nil
}
