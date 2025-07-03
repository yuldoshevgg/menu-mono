package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/yuldoshevgg/menu-mono/internal/config"
	"google.golang.org/grpc/metadata"
)

var publicEndpoints = map[string]string{
	"/v1/menu":        "GET",
	"/v1/access-menu": "POST",
}

func AuthFn(ctx context.Context) (context.Context, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		httpPath := getStringOrEmptyItem(md, config.GrpcGatewayHTTPPath)
		httpMethod := getStringOrEmptyItem(md, config.GrpcGatewayHTTPMethod)

		if key := md.Get(string(config.VenueID)); len(key) > 0 {
			ctx = context.WithValue(ctx, config.VenueID, key[0])
		}

		if key := md.Get(string(config.TableNumber)); len(key) > 0 {
			ctx = context.WithValue(ctx, config.TableNumber, key[0])
		}

		if publicEndpoints[httpPath] == httpMethod {
			return ctx, nil
		}

		if key := md.Get(string(config.Authorization)); len(key) > 0 {
			token := strings.Split(key[0], " ")
			if (len(token) == 2 && token[0] != "Bearer") || len(token) != 2 {
				return nil, fmt.Errorf("invalid token provided")
			}
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	return ctx, nil
}

func getStringOrEmptyItem(md metadata.MD, key string) string {
	if len(md[key]) > 0 {
		return md[key][0]
	}

	return ""
}
