package log

import (
	"time"

	"go.uber.org/zap"
)

type Field = zap.Field

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}
