package rmq

import (
	"context"
	"strings"

	"github.com/streadway/amqp"
)

// ----------------------------------------------------------------------------
// context section
// ----------------------------------------------------------------------------

type metadataContextKeyType string
type metadata map[string]string

const (
	metadataContextKey metadataContextKeyType = "rmq.metadata"
)

func ContextMetadata(ctx context.Context) metadata {
	meta := ctx.Value(metadataContextKey)
	if meta == nil {
		return nil
	}
	return meta.(metadata)
}

func ContextWithMetaValue(ctx context.Context, key, value string) context.Context {
	meta := ContextMetadata(ctx)
	if meta == nil { // create a new map if not exist
		meta = make(metadata)
	}
	meta[key] = value // upsert the value
	return context.WithValue(ctx, metadataContextKey, meta)
}

func ContextMetaValue(ctx context.Context, key string) string {
	meta := ContextMetadata(ctx)
	if meta == nil {
		return ""
	}
	return meta[key]
}

// ----------------------------------------------------------------------------
// amqp headers section
// ----------------------------------------------------------------------------

const (
	responseCodeHeaderKey   = "rmq.responseCode"
	metadataHeaderPrefixKey = "rmq.metadata."

	responseOK    int16 = 0
	responseError int16 = 1
)

func HasError(message *amqp.Delivery) bool {
	headers := message.Headers
	code, exist := headers[responseCodeHeaderKey]
	if !exist { // if there is no such header, we are ok
		return false
	}
	if code != responseError {
		return false
	}
	return true
}

func WithError(message *amqp.Publishing) *amqp.Publishing {
	if message.Headers == nil {
		message.Headers = make(amqp.Table)
	}
	message.Headers[responseCodeHeaderKey] = responseError
	return message
}

func Metadata(message *amqp.Delivery) metadata {
	if message.Headers == nil {
		return nil
	}

	meta := make(metadata)
	for key, val := range message.Headers {
		if strings.HasPrefix(key, metadataHeaderPrefixKey) {
			meta[strings.TrimPrefix(key, metadataHeaderPrefixKey)] = val.(string)
		}
	}
	return meta
}

func WithMetadata(message *amqp.Publishing, metadata metadata) *amqp.Publishing {
	if message.Headers == nil {
		message.Headers = make(amqp.Table)
	}
	for key, val := range metadata {
		message.Headers[metadataHeaderPrefixKey+key] = val
	}
	return message
}
