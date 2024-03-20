package orm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type loggerOptions struct {
	SlowThreshold time.Duration
	LogLevel      string
	Logger        log.Logger
}

type LoggerOption func(*loggerOptions)

func WithSlowThreshold(slowThreshold time.Duration) LoggerOption {
	return func(o *loggerOptions) {
		o.SlowThreshold = slowThreshold
	}
}

func WithLogLevel(level string) LoggerOption {
	return func(o *loggerOptions) {
		o.LogLevel = level
	}
}

func WithLogger(logger log.Logger) LoggerOption {
	return func(o *loggerOptions) {
		o.Logger = logger
	}
}

type Logger struct {
	gormlogger.Config
	log *log.Helper
}

func NewLogger(opts ...LoggerOption) gormlogger.Interface {
	var opt loggerOptions
	for _, o := range opts {
		o(&opt)
	}
	if opt.Logger == nil {
		return gormlogger.Default.LogMode(getLogLevel(opt.LogLevel))
	}
	if opt.SlowThreshold <= 0 {
		opt.SlowThreshold = 200 * time.Millisecond
	}
	return &Logger{
		Config: gormlogger.Config{
			SlowThreshold:             opt.SlowThreshold,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  getLogLevel(opt.LogLevel),
		},
		log: log.NewHelper(opt.Logger),
	}
}

// LogMode log mode
func (l *Logger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.log.WithContext(ctx).Infof(msg, data)
}

// Warn print warn messages
func (l Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.log.WithContext(ctx).Warnf(msg, data)
}

// Error print error messages
func (l Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.log.WithContext(ctx).Errorf(msg, data)
}

// Trace print sql message
func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, _ := fc()
		l.log.WithContext(ctx).Errorf("elapsed:%+vms sql:%v err:%v", float64(elapsed.Nanoseconds())/1e6, sql, err)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.log.WithContext(ctx).Warnf("elapsed:%+vms slowLog:%v sql:%v rows:%v", float64(elapsed.Nanoseconds())/1e6, slowLog, sql, rows)
	case l.LogLevel == gormlogger.Info:
		sql, rows := fc()
		l.log.WithContext(ctx).Infof("elapsed:%+vms sql:%v rows:%v", float64(elapsed.Nanoseconds())/1e6, sql, rows)
	}
}

func getLogLevel(level string) gormlogger.LogLevel {
	var logLevel gormlogger.LogLevel
	switch level {
	case "error":
		logLevel = gormlogger.Error
	case "warn":
		logLevel = gormlogger.Warn
	case "info":
		logLevel = gormlogger.Info
	default:
		logLevel = gormlogger.Silent
	}
	return logLevel
}
