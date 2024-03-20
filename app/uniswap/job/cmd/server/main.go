package main

import (
	"flag"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"

	"uniswap-transaction/app/uniswap/job/internal/conf"
	"uniswap-transaction/app/uniswap/job/internal/server"
	zlog "uniswap-transaction/pkg/log/zap"
	"uniswap-transaction/pkg/metadata"
	"uniswap-transaction/pkg/metric"
	gconf "uniswap-transaction/protobuf/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string = "uniswap.job"
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
	pid   = os.Getpid()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	registerMetrics()
}

func registerMetrics() {
	prometheus.MustRegister(metric.Seconds, metric.Requests)
}

func setTracerProvider(conf *gconf.Trace) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.Endpoint)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(Name),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func newApp(logger log.Logger, em gconf.EnvMode, hs *http.Server, es *server.EthereumServer) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{
			metadata.EnvMode: em.String(),
		}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			es,
		),
	)
}

func main() {
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
			zap.AddCallerSkip(4),
		},
	})
	logger := log.With(zlog.NewLogger(zapLogger),
		// "ts", log.DefaultTimestamp,
		// "caller", log.Caller(4),
		"hostname", id,
		"pid", pid,
		"service", Name,
		"version", Version,
		"traceid", tracing.TraceID(),
		"spanid", tracing.SpanID(),
	)

	if err := setTracerProvider(bc.Trace); err != nil {
		panic(err)
	}

	app, cleanup, err := initApp(bc.EnvMode, bc.Server, bc.Data, bc.Alert, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
