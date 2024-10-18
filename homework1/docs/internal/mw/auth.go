package mw

import (
	"context"
	cliserver "homework1/pkg/cli/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	switch req.(type) {
	case
		*cliserver.AcceptOrderRequest,
		*cliserver.AcceptReturnRequest,
		*cliserver.GiveOrderRequest,
		*cliserver.ReturnOrderRequest:

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "can not parse metada")
		}
		v := md.Get("x-api-token")
		if len(v) == 0 {
			return nil, status.Error(codes.Unauthenticated, "can not parse token")
		}
	}
	return handler(ctx, req)

}
