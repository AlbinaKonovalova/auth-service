package httpinfra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"

	pb "github.com/AlbinaKonovalova/auth-service/pkg/authservice/v1"
)

func NewGatewayMux(ctx context.Context, svc pb.AuthServiceServer) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "Cookie", "cookie":
				return key, true
			}
			return runtime.DefaultHeaderMatcher(key)
		}),
		runtime.WithForwardResponseOption(forwardCookieHeader),
	)

	if err := pb.RegisterAuthServiceHandlerServer(ctx, mux, svc); err != nil {
		return nil, fmt.Errorf("register auth service gateway: %w", err)
	}

	return mux, nil
}

func forwardCookieHeader(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}
	for _, v := range md.HeaderMD.Get("set-cookie") {
		w.Header().Add("Set-Cookie", v)
	}
	return nil
}
