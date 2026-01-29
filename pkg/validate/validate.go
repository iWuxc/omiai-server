package validate

import (
	"regexp"

	validator2 "github.com/go-playground/validator/v10"
	"github.com/iWuxc/go-wit/validator"
)

type validate struct{}

func (v validate) Validators() []validator.CustomRule {
	return []validator.CustomRule{
		{
			Name: "base64Format",
			Rule: base64Format,
			Msg:  "{0}格式不正确",
		},
		{
			Name: "checkTaxId",
			Msg:  "{0}必须是一个有效的纳税人识别号",
			Rule: checkTaxId,
		},
		{
			Name: "mobile",
			Rule: mobile,
			Msg:  "{0}必须是一个有效的手机号",
		},
		{
			Name: "date",
			Rule: date,
			Msg:  "{0}格式不正确",
		},
	}
}

func NewValidator() validator.Validate {
	return validate{}
}

func base64Format(fl validator2.FieldLevel) bool {
	re := `^data:\s*image\/(\w+);base64,`
	r := regexp.MustCompile(re)
	return r.MatchString(fl.Field().String())
}

// 纳税人识别号校验
func checkTaxId(fl validator2.FieldLevel) bool {
	var regArr = []string{
		`^[\da-zA-Z]{10,15}$`,
		`^\d{6}[\da-zA-Z]{10,12}$`,
		`^[a-zA-Z]\d{6}[\da-zA-Z]{9,11}$`,
		`^[a-zA-Z]{2}\d{6}[\da-zA-Z]{8,10}$`,
		`^\d{14}[\d][\da-zA-Z]{4,5}$`,
		`^\d{17}[\d][\da-zA-Z]{1,2}$`,
		`^[a-zA-Z]\d{14}[\d][\da-zA-Z]{3,4}$`,
		`^[a-zA-Z]\d{17}[\d][\da-zA-Z]{0,1}$`,
		`^[\d]{6}[\da-zA-Z]{13,14}$`,
		`^[0-9A-HJ-NPQRTUWXY]{2}\d{6}[0-9A-HJ-NPQRTUWXY]{10}$`,
	}
	for _, re := range regArr {
		r := regexp.MustCompile(re)
		if r.MatchString(fl.Field().String()) {
			return true
		}
	}
	return false
}

// 手机号校验
func mobile(fl validator2.FieldLevel) bool {
	re := `^1[123456789]\d{9}$`
	r := regexp.MustCompile(re)
	return r.MatchString(fl.Field().String())
}

// 时间校验
func date(fl validator2.FieldLevel) bool {
	timeFormat := fl.Field().String()
	timeRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return timeRegex.MatchString(timeFormat)
}
