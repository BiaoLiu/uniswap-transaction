package zap

import "testing"

func TestNewZapLogger(t *testing.T) {
	logger := NewZapLogger(&Config{
		Level:    "debug",
		Filename: "",
	})
	logger.Info("test stdout...")

	logger = NewZapLogger(&Config{
		Level:    "debug",
		Filename: "log",
	})
	logger.Info("test file")
}
