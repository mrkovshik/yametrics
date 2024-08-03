package grpc

import (
	"context"
	"crypto/hmac"
	"encoding/json"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/mrkovshik/yametrics/internal/signature"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func InterceptorLogger(l *zap.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		f := make([]zap.Field, 0, len(fields)/2)

		for i := 0; i < len(fields); i += 2 {
			key := fields[i]
			value := fields[i+1]

			switch v := value.(type) {
			case string:
				f = append(f, zap.String(key.(string), v))
			case int:
				f = append(f, zap.Int(key.(string), v))
			case bool:
				f = append(f, zap.Bool(key.(string), v))
			default:
				f = append(f, zap.Any(key.(string), v))
			}
		}

		logger := l.WithOptions(zap.AddCallerSkip(1)).With(f...)

		switch lvl {
		case logging.LevelDebug:
			logger.Debug(msg)
		case logging.LevelInfo:
			logger.Info(msg)
		case logging.LevelWarn:
			logger.Warn(msg)
		case logging.LevelError:
			logger.Error(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func Authenticate(key string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if key == "" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
		}
		values := md.Get(`HashSHA256`)
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing hmac-signature")
		}
		clientSig := values[0]
		messageBytes, err := json.Marshal(req)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to marshal request: %v", err)
		}

		sigSrv := signature.NewSha256Sig(key, messageBytes)
		srvSig, err := sigSrv.Generate()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to generate hash: %v", err)
		}
		if !hmac.Equal([]byte(clientSig), []byte(srvSig)) {
			return nil, status.Errorf(codes.Unauthenticated, "invalid hmac-signature")
		}
		return handler(ctx, req)
	}
}
