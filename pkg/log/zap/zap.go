package zap

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"uniswap-transaction/pkg/file"
)

type Config struct {
	Level    string
	Filename string
	Options  []zap.Option
}

func NewZapLogger(c *Config) *zap.Logger {
	level := zapcore.Level(log.ParseLevel(c.Level))
	encoder := getEncoder()
	writeSyncer := getLogWriter(c.Filename)
	core := zapcore.NewCore(encoder, writeSyncer, level)
	logger := zap.New(core, c.Options...)
	return logger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:   "ts",
		LevelKey:  "level",
		NameKey:   "logger",
		CallerKey: "caller",
		//MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string) zapcore.WriteSyncer {
	if filename == "" {
		return os.Stdout
	}
	if dir := path.Dir(filename); dir != "" {
		if err := file.CheckCreateDir(dir); err != nil {
			panic(fmt.Sprintf("config local file system logger error. %+v", errors.WithStack(err)))
		}
	}
	// 保存日志30天，每24小时分割一次日志
	filenameWithoutExt := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	linkName := filepath.Join(filepath.Dir(filename), filenameWithoutExt+"_%Y%m%d"+filepath.Ext(filename))
	rotateLogs, err := rotatelogs.New(
		linkName,
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(fmt.Sprintf("config local file system logger error. %+v", errors.WithStack(err)))
	}
	return zapcore.AddSync(rotateLogs)
}
