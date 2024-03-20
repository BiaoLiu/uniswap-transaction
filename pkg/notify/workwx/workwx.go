package workwx

import (
	"fmt"

	"github.com/valyala/fasttemplate"
	"github.com/xen0n/go-workwx"

	"uniswap-transaction/pkg/notify"
)

type Config struct {
	Key string
}

type WorkWx struct {
	*workwx.WebhookClient
}

func NewWorkWx(conf *Config) notify.Client {
	wh := workwx.NewWebhookClient(conf.Key)
	return &WorkWx{WebhookClient: wh}
}

func (a *WorkWx) SendMessage(content string) error {
	return a.SendMarkdownMessage(content)
}

func (a *WorkWx) SendTemplateMessage(template string, args map[string]interface{}) error {
	tmpl, err := fasttemplate.NewTemplate(template, "{{", "}}")
	if err != nil {
		return fmt.Errorf("create template error. err:%v", err)
	}
	content := tmpl.ExecuteString(args)
	return a.SendMarkdownMessage(content)
}
