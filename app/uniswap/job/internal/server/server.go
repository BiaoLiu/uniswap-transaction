package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"uniswap-transaction/app/uniswap/job/internal/conf"
	"uniswap-transaction/pkg/alert"
	"uniswap-transaction/pkg/notify"
	"uniswap-transaction/pkg/notify/feishu"
	gconf "uniswap-transaction/protobuf/conf"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(
	NewEthereumServer,
	NewHTTPServer,
	NewAlertTemplate,
	NewNotifyClient,
	NewAlertManager,
)

func NewNotifyClient(conf *conf.Alert) notify.Client {
	// c := &workwx.Config{
	//	Key: conf.WorkWx.Key,
	// }
	// return workwx.NewWorkWx(c)
	c := &feishu.Config{
		Key:       conf.Feishu.Key,
		SecretKey: conf.Feishu.SecretKey,
	}
	return feishu.NewFeiShu(c)
}

func NewAlertTemplate() alert.Template {
	return alert.NewFeishuTemplate()
}

func NewAlertManager(client notify.Client, template alert.Template, logger log.Logger) *alert.Manager {
	var opts []alert.Option
	opts = append(opts, alert.WithAlertEnvMode(
		// gconf.EnvMode_DEBUG.String(),
		gconf.EnvMode_DEVELOP.String(),
		gconf.EnvMode_TEST.String(),
		gconf.EnvMode_PRODUCTION.String(),
	))
	opts = append(opts, alert.WithIgnoreCode(505, 506, 521))
	return alert.NewAlertManager(client, template, logger, opts...)
}
