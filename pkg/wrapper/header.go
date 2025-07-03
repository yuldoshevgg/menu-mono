package wrapper

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "X-Venue-Id":
		return key, true
	case "X-Table-Number":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
