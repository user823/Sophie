package validators

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func GetValidatorOr() *validator.Validate {
	once.Do(func() {
		validate = validator.New(validator.WithRequiredStructEnabled())
		registerValidators()
	})
	return validate
}

func registerValidators() {
	validate.RegisterValidation("xss", xssValidator)

}

func BuildErrMsgHelper(err validator.FieldError, fieldname string, args ...any) error {
	switch err.ActualTag() {
	case "min":
		return fmt.Errorf("%s的长度必须不能小于%v", fieldname, err.Param())
	case "max":
		return fmt.Errorf("%s的长度必须不能超过%v", fieldname, err.Param())
	case "xss":
		return fmt.Errorf("%s不能包含脚本字符")
	case "email":
		return fmt.Errorf("邮箱格式不正确")
	}
	return fmt.Errorf("")
}
