package middleware

import (
	"context"
	"encoding/json"

	"github.com/yuldoshevgg/menu-mono/internal/config"
	"github.com/yuldoshevgg/menu-mono/internal/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	headerKeyAdminUserID  = "X-Admin-User-Id"
	headerKeyMerchant     = "X-Merchant"
	headerKeyPlatform     = "X-Platform-Hash"
	headerKeyOperationMsg = "X-Operation-Msg"
	headerKeyUserID       = "X-User-Id"
	headerKeyTaskID       = "X-Task-Id"
	headerKeyRequestedID  = "X-Requested-Id"
	headerAPIKey          = "x-api-key"
	headerSellerID        = "x-seller-id"
)

type Metadata struct {
	Platform     string
	AdminID      string
	UserID       string
	TaskID       string
	Merchant     string
	OperationMsg string
	Path         string
	Body         map[string]interface{}
	Time         string
	RequestBasic string
}

func ConvertToMap(data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, err
	}

	redactSensitiveInfo(jsonMap)
	return jsonMap, nil
}

func redactSensitiveInfo(data map[string]interface{}) {
	var sensitiveInfo = []string{
		"password",
		"front",
		"card_number",
		"pass_data",
		"number",
		"pinfl",
		"client_photo",
		"doc_number",
		"doc_series",
	}

	for _, key := range sensitiveInfo {
		redactKey(data, key)
	}
}

func redactKey(data map[string]interface{}, key string) {
	for k, v := range data {
		if k == key {
			data[k] = "[REDACTED]"
		}
		if nestedMap, ok := v.(map[string]interface{}); ok {
			redactKey(nestedMap, key)
		}
	}
}

func extractMetadata(ctx context.Context, req interface{}) *Metadata {
	var (
		mtdt = &Metadata{}
		err  error
	)

	mtdt.Body, err = ConvertToMap(req)
	if err != nil {
		logger.Log.Error("failed to convert to map: ", zap.Error(err))
		return nil
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if key := md.Get(headerKeyPlatform); len(key) > 0 {
			mtdt.Platform = key[0]
		}
		if key := md.Get(headerKeyAdminUserID); len(key) > 0 {
			mtdt.AdminID = key[0]
		}
		if key := md.Get(headerKeyUserID); len(key) > 0 {
			mtdt.UserID = key[0]
		}
		if key := md.Get(headerKeyTaskID); len(key) > 0 {
			mtdt.TaskID = key[0]
		}
		if key := md.Get(headerKeyMerchant); len(key) > 0 {
			mtdt.Merchant = key[0]
		}
		if key := md.Get(headerKeyOperationMsg); len(key) > 0 {
			mtdt.OperationMsg = key[0]
		}
		if key := md.Get(config.GrpcGatewayHTTPPath); len(key) > 0 {
			mtdt.Path = key[0]
		}
	}

	return mtdt
}
