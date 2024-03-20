package alert

type Template interface {
	ServiceMessageContent() string
}

type WorkWxTemplate struct {
}

func NewWorkWxTemplate() Template {
	return &WorkWxTemplate{}
}

func (t *WorkWxTemplate) ServiceMessageContent() string {
	return `
**服务告警，请注意：**
>**environment**: <font color="warning">{{environment}}</font>
>**ts**: {{ts}}
>**hostname**: {{hostname}}
>**service**: <font color="warning">{{service}}</font>
>**traceid**: {{traceid}}
>**spanid**: {{spanid}}
>**kind**: {{kind}}
>**component**: {{component}}
>**operation**: <font color="warning">{{operation}}</font>
>**userid**: {{userid}}
>**username**: {{username}}
>**args**: {{args}}
>**code**: <font color="warning">{{code}}</font>
>**reason**: <font color="warning">{{reason}}</font>
>**stack**: {{stack}}
`
}

type FeishuTemplate struct {
}

func NewFeishuTemplate() Template {
	return &FeishuTemplate{}
}

func (t *FeishuTemplate) ServiceMessageContent() string {
	return `
服务告警，请注意：
>environment: {{environment}}
>ts: {{ts}}
>hostname: {{hostname}}
>service: {{service}}
>traceid: {{traceid}}
>spanid: {{spanid}}
>kind: {{kind}}
>component: {{component}}
>operation: {{operation}}
>userid: {{userid}}
>username: {{username}}
>args: {{args}}
>code: {{code}}
>reason: {{reason}}
>stack: {{stack}}
`
}
