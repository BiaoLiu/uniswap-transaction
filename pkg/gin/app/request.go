package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	"uniswap-transaction/pkg/validator"
)

const (
	bind = "bind"
)

type Paginate struct {
	PageNum  int
	PageSize int
}

// BindAndValidate 参数绑定与校验，根据Content-Type进行绑定
func BindAndValidate(c *gin.Context, dst interface{}) error {
	if err := Bind(c, dst); err != nil {
		return err
	}
	return validator.Validate(dst)
}

// ShouldBindAndValidate 参数绑定与校验，指定参数类型
func ShouldBindAndValidate(c *gin.Context, dst interface{}, binding binding.Binding) error {
	if err := shouldBindWith(c, dst, binding); err != nil {
		return err
	}
	return validator.Validate(dst)
}

// Bind 参数绑定，根据Content-Type进行参数绑定
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	if err := shouldBindWith(c, obj, b); err != nil {
		return err
	}
	return nil
}

// shouldBindWith 参数绑定，指定参数类型
func shouldBindWith(c *gin.Context, obj interface{}, binding binding.Binding) error {
	if err := c.ShouldBindWith(obj, binding); err != nil {
		return errors.WithMessage(err, "参数绑定失败")
	}
	SetBindContext(c, obj)
	return nil
}

// SetBindContext 设置绑定参数上下文
func SetBindContext(c *gin.Context, obj interface{}) {
	c.Set(bind, obj)
}

// FromBindContext 获取绑定参数
func FromBindContext(c *gin.Context) interface{} {
	obj, _ := c.Get(bind)
	return obj
}
