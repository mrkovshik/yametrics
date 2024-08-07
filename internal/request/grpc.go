package request

import (
	"encoding/json"
	"time"

	"github.com/mrkovshik/yametrics/internal/reqbuilder"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	config "github.com/mrkovshik/yametrics/internal/config/agent"
	"github.com/mrkovshik/yametrics/internal/model"
	"github.com/mrkovshik/yametrics/proto"
)

// GRPCClient represents a client for interacting with a gRPC service.
type GRPCClient struct {
	logger            *zap.SugaredLogger  // Logger for logging messages.
	cfg               *config.AgentConfig // Configuration for the agent.
	proto.UsersClient                     // gRPC client for interacting with the Users service.
}

// NewGRPCClient creates a new GRPCClient instance.
//
// Parameters:
//   - logger: A SugaredLogger instance for logging.
//   - cfg: The agent configuration.
//   - conn: The gRPC client connection.
//
// Returns:
//   - *GRPCClient: A new GRPCClient instance.
func NewGRPCClient(logger *zap.SugaredLogger, cfg *config.AgentConfig, conn *grpc.ClientConn) *GRPCClient {
	return &GRPCClient{
		UsersClient: proto.NewUsersClient(conn),
		logger:      logger,
		cfg:         cfg,
	}
}

// Send processes and sends metrics to the gRPC server.
//
// Parameters:
//   - id: The worker ID.
//   - jobs: A channel of metrics to be sent.
func (r *GRPCClient) Send(id int, jobs <-chan model.Metrics) {
	for j := range jobs {
		r.logger.Debugf("worker #%v is sending %v", id, j.ID)
		var metric = proto.Metric{
			ID:    j.ID,
			MType: j.MType,
		}
		switch j.MType {
		case model.MetricTypeGauge:
			metric.Value = *j.Value
		case model.MetricTypeCounter:
			metric.Delta = *j.Delta
		}

		request := proto.UpdateMetricsRequest{
			Metrics: []*proto.Metric{
				&metric,
			},
		}
		response, err := r.retryableSend(&request)
		if err != nil {
			r.logger.Errorf("error sending request: %v\n", err)
			return
		}
		if response.GetError() != "" {
			r.logger.Errorf("status code is %v\n", response.GetError())
			return
		}
	}
}

// retryableSend sends a request with retry logic.
//
// Parameters:
//   - req: The UpdateMetricsRequest to be sent.
//
// Returns:
//   - *proto.UpdateMetricsResponse: The response from the server.
//   - error: Any error encountered during the send operation.
func (r *GRPCClient) retryableSend(req *proto.UpdateMetricsRequest) (*proto.UpdateMetricsResponse, error) {
	var retryIntervals = []int{1, 3, 5} // TODO: move to config
	messageBytes, err := json.Marshal(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
	}
	ctxBuilder := reqbuilder.NewGRPCContextBuilder().Sign(r.cfg.Key, messageBytes)
	if ctxBuilder.Err != nil {
		return nil, ctxBuilder.Err
	}
	for i := 0; i <= len(retryIntervals); i++ {
		response, err := r.UpdateMetrics(ctxBuilder.Ctx, req)
		if err == nil {
			return response, nil
		}
		if i == len(retryIntervals) {
			return nil, err
		}
		r.logger.Errorf("failed to connect to server: %v\n retrying in %v seconds\n", err, retryIntervals[i])
		time.Sleep(time.Duration(retryIntervals[i]) * time.Second)
	}
	return nil, nil
}