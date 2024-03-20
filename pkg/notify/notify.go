package notify

type Client interface {
	SendMessage(content string) error
	SendTemplateMessage(template string, args map[string]interface{}) error
}
