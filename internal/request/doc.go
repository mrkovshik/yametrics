// Package request provides clients for interacting with remote services via gRPC and HTTP.
//
// The package defines two main clients: GRPCClient and RestClient. These clients
// handle sending metrics to remote servers using gRPC and HTTP respectively, with
// retry mechanisms in case of failures.
//
// The clients use a structured logging approach with zap.SugaredLogger and support
// configurable retry intervals.
//
// # GRPCClient
//
// The GRPCClient struct represents a client for interacting with a gRPC service.
// It includes methods to send metrics and handle retries.
//
// Example usage:
//
//	logger := zap.NewExample().Sugar()
//	cfg := &config.AgentConfig{
//		// Populate configuration fields
//	}
//	conn, err := grpc.Dial("address", grpc.WithInsecure())
//	if err != nil {
//		logger.Fatal(err)
//	}
//	client := request.NewGRPCClient(logger, cfg, conn)
//	jobs := make(chan model.Metrics)
//	go client.Send(1, jobs)
//
// # RestClient
//
// The RestClient struct represents a client for interacting with a REST service.
// It includes methods to send metrics using HTTP and handle retries.
//
// Example usage:
//
//	logger := zap.NewExample().Sugar()
//	cfg := &config.AgentConfig{
//		// Populate configuration fields
//	}
//	client := request.NewRestClient(logger, cfg)
//	jobs := make(chan model.Metrics)
//	go client.Send(1, jobs)
//
// Both clients provide retry mechanisms that attempt to resend requests based on
// configurable intervals.
package request
