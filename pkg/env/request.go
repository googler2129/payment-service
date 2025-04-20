package env

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mercor/payment-service/constants"
)

// RequestID checks the X-Request-ID header and generate new if not found
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(constants.HeaderXMercorRequestID)

		if requestID == "" {
			requestID = NewRequestID()
		}

		// Setting Request ID in newRelic
		c.Set(constants.HeaderXMercorRequestID, requestID)
		c.Writer.Header().Set(constants.HeaderXMercorRequestID, requestID)
		c.Next()
	}
}

func NewRequestID() string {
	return uuid.New().String()
}

// GetRequestID returns the request ID from context
func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(constants.HeaderXMercorRequestID).(string)
	if !ok {
		return uuid.New().String()
	}

	return requestID
}

// GetRequestID returns the request ID from context
func GetRequestIDForPostgresqlLogging(ctx context.Context) string {
	requestID, ok := ctx.Value(constants.HeaderXMercorRequestID).(string)
	if !ok {
		return ""
	}

	return requestID
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, constants.HeaderXMercorRequestID, requestID)
}

func GetSQSMessageRequestID(ctx context.Context, attributes map[string]string) string {
	requestID, ok := attributes[constants.HeaderXMercorRequestID]
	if !ok {
		return uuid.New().String()
	}

	return requestID
}

func SetSqsMessageRequestID(ctx context.Context, attributes map[string]string) context.Context {
	requestID := GetSQSMessageRequestID(ctx, attributes)
	// Setting newrelic request ID
	return context.WithValue(ctx, constants.HeaderXMercorRequestID, requestID)
}

func GetKafkaRequestID(ctx context.Context, headers map[string]string) string {
	requestID, ok := headers[constants.HeaderXMercorRequestID]
	if !ok {
		return uuid.New().String()
	}

	return requestID
}

func SetKafkaRequestID(ctx context.Context, attributes map[string]string) context.Context {
	requestID := GetKafkaRequestID(ctx, attributes)

	// Setting newrelic request ID
	return context.WithValue(ctx, constants.HeaderXMercorRequestID, requestID)
}
