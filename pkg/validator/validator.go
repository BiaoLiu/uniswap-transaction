package validator

import (
	zhLocales "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslation "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
)

var (
	zh       = zhLocales.New()
	uni      = ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
)

func init() {
	_ = zhTranslation.RegisterDefaultTranslations(validate, trans)
}

// Validate 参数校验
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		errs := err.(validator.ValidationErrors)[0]
		return errors.New("参数校验失败: " + errs.Translate(trans))
	}
	return nil
}
