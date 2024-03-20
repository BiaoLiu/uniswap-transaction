package server

import (
	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/transport/http/pprof"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"uniswap-transaction/app/uniswap/job/internal/conf"
	"uniswap-transaction/app/uniswap/job/internal/service"
	"uniswap-transaction/pkg/alert"
	"uniswap-transaction/pkg/encoder"
	"uniswap-transaction/pkg/metric"
	"uniswap-transaction/pkg/middleware/alerting"
	"uniswap-transaction/pkg/middleware/logging"
	gconf "uniswap-transaction/protobuf/conf"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(em gconf.EnvMode, c *conf.Server, data *conf.Data, logger log.Logger, alertManager *alert.Manager,
	ethereumService *service.EthereumService) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			metadata.Server(metadata.WithPropagatedPrefix("")),
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(metric.Seconds)),
				metrics.WithRequests(prom.NewCounter(metric.Requests)),
			),
			validate.Validator(),
			logging.Server(logger, false),
			alerting.Server(alertManager),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.ResponseEncoder(encoder.ResponseEncoder), http.ErrorEncoder(encoder.ErrorEncoder))
	srv := http.NewServer(opts...)
	srv.Handle("/", pprof.NewHandler())
	srv.Handle("/metrics", promhttp.Handler())
	return srv
}
