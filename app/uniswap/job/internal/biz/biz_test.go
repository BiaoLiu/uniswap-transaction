package biz

import (
	"flag"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"

	"uniswap-transaction/app/uniswap/job/internal/conf"
	zlog "uniswap-transaction/pkg/log/zap"
)

var (
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs_local", "config path, eg: -conf config.yaml")
	// registerMetrics()
}

var (
	biz *bizTest
)

func TestMain(t *testing.M) {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	zapLogger := zlog.NewZapLogger(&zlog.Config{
		Level:    bc.Log.Level,
		Filename: bc.Log.Filename,
		Options: []zap.Option{
			zap.AddCaller(),
			zap.AddCallerSkip(4)},
	})
	logger := log.With(zlog.NewLogger(zapLogger))

	var err error
	biz, _, err = initBizTest(bc.EnvMode, bc.Data, bc.Server, logger)
	if err != nil {
		panic(err)
	}
	os.Exit(t.Run())
}
