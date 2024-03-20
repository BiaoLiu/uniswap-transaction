package feishu

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasttemplate"

	fsBotAPI "uniswap-transaction/pkg/feishu"
	"uniswap-transaction/pkg/notify"
	"uniswap-transaction/pkg/utils"
)

type Config struct {
	Key       string
	SecretKey string
}

type FeiShu struct {
	fsBotAPI.Bot
}

func NewFeiShu(conf *Config) notify.Client {
	bot := fsBotAPI.NewBot(conf.Key, fsBotAPI.WithSecretKey(conf.SecretKey))
	return &FeiShu{
		Bot: bot,
	}
}

func (a *FeiShu) SendMessage(content string) error {
	var raw map[string]interface{}
	err := json.Unmarshal(utils.StringToBytes(content), &raw)
	if err != nil {
		return err
	}
	return a.Bot.PushRawText(raw)
}

func (a *FeiShu) SendTemplateMessage(template string, args map[string]interface{}) error {
	tmpl, err := fasttemplate.NewTemplate(template, "{{", "}}")
	if err != nil {
		return fmt.Errorf("create template error. err:%v", err)
	}
	content := tmpl.ExecuteString(args)
	return a.Bot.PushText(content)
}
