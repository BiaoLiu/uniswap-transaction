package xxl

import (
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type Logger struct {
	log *log.Helper
}

func NewLogger(logger log.Logger) *Logger {
	return &Logger{
		log: log.NewHelper(logger),
	}
}

func (l *Logger) Info(format string, a ...interface{}) {
	if strings.Contains(format, "执行器注册成功") {
		return
	}
	if a == nil {
		l.log.Infof(format)
	} else {
		l.log.Infof(format, a)
	}
}

func (l *Logger) Error(format string, a ...interface{}) {
	if a == nil {
		l.log.Errorf(format)
	} else {
		l.log.Errorf(format, a)
	}
}
