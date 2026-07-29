package grpcutil

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/art-es/b2b-usage-based-billing-platform/services/api-gateway/internal/auth"
)

func CallOpts(ctx context.Context) []grpc.CallOption {
	md := metadata.New(map[string]string{})

	if val, ok := auth.Get(ctx); ok {
		md.Set("authorization", val)
	}

	return []grpc.CallOption{
		grpc.Header(&md),
	}
}
