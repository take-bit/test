package grpcAuthController

import (
	"context"
	grpcGen "github.com/take-bit/proto/gen/auth/go/v1"
)

func (c *Controller) Login(ctx context.Context, req *grpcGen.LoginRequest) (*grpcGen.LoginResponse, error) {
	panic("todo")
}
