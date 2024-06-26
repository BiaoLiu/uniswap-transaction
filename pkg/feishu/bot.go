package fsBotAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	webhook = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"
)

type bot struct {
	webhook string

	secretKey string
}

func NewBot(key string, opts ...BotOption) Bot {
	wh := fmt.Sprintf(webhook, key)
	b := new(bot)
	if !strings.Contains(wh, "open.feishu.cn") {
		b.webhook = fmt.Sprintf(_fmtWebhook, wh)
	} else {
		b.webhook = wh
	}
	for _, fn := range opts {
		fn(b)
	}
	return b
}

type BotOption func(*bot)

func WithSecretKey(key string) BotOption {
	return func(b *bot) {
		b.secretKey = strings.TrimSpace(key)
	}
}

func (b *bot) PushText(content string) error {
	return b.pushMsg(newMsgText(content))
}

func (b *bot) PushPost(p Post, ps ...Post) error {
	return b.pushMsg(newMsgPost(p, ps...))
}

func (b *bot) PushCard(bgColor CardTitleBgColor, cfg CardConfig, c Card, more ...Card) error {
	return b.pushMsg(GenMsgCard(bgColor, cfg, c, more...))
}

func (b *bot) PushImage(imageKey string) error {
	return b.pushMsg(newMsgImage(imageKey))
}

func (b *bot) PushShareChat(chatID string) error {
	return b.pushMsg(newMsgShareChat(chatID))
}

func (b *bot) PushRawText(content map[string]interface{}) error {
	msg := map[string]interface{}{
		"msg_type": "interactive",
		"card":     content,
	}
	return b.pushMsg(msg)
}

func (b *bot) pushMsg(msg map[string]interface{}) (err error) {
	if b.secretKey != "" {
		ts := time.Now().Unix()
		signed, err := genSign(b.secretKey, ts)
		if err != nil {
			return err
		}
		msg["timestamp"] = ts
		msg["sign"] = signed
	}
	var bsJSON []byte
	if bsJSON, err = json.Marshal(msg); err != nil {
		return err
	}
	var req *http.Request
	if req, err = newRequest(http.MethodPost, b.webhook, bsJSON); err != nil {
		return err
	}

	_, err = executeHTTP(req)
	return
}
