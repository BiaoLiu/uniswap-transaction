package zap

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	zlogger := NewZapLogger(&Config{
		Level:    "debug",
		Filename: "",
		Options:  []zap.Option{zap.AddCaller(), zap.AddCallerSkip(4)},
	})
	logger := log.With(NewLogger(zlogger))
	log := log.NewHelper(logger)
	log.Debug("test")
}
