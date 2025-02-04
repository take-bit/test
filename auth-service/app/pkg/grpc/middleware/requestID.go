package grpcMiddleware

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func RequestID(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = context.WithValue(ctx, "request_id", uuid.NewString())

	return handler(ctx, req)
}
