package middleware

import (
	"context"
	"runtime"

	"github.com/yuldoshevgg/menu-mono/internal/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	StackTraceMaxCap = 4096
)

func GrpcGateWayRecovery(_ context.Context, p interface{}) (err error) {
	stackTrace := make([]byte, StackTraceMaxCap)
	runtime.Stack(stackTrace, false)
	// Log panic details along with stack trace
	logger.Log.Error("panic recovery", zap.Any("error", p), zap.String("panic", string(stackTrace)))
	return status.Errorf(codes.Unknown, "panic triggered: %v", p)
}
