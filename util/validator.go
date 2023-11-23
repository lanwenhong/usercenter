package util

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func EmailValidator(f validator.FieldLevel) bool {
	value := f.Field().String()
	//ctx := context.Background()
	//logger.Debugf(ctx, "------------------------------")
	if value == "" {
		return true
	} else {
		if len(value) > 128 {
			return false
		}
		//logger.Debugf(ctx, "value: %s", value)
		return VerifyEmailFormat(value)
	}
	return true
}

func PwdValidator(f validator.FieldLevel) bool {
	val := f.Field().String()
	if len(val) < 8 || len(val) > 20 { // length需要通过验证
		fmt.Println("pwd length error")
		return false
	}

	pwdPattern := `^[0-9a-zA-Z!@#$%^&*~-_+]{8,20}$`
	reg, err := regexp.Compile(pwdPattern) // filter exclude chars
	if err != nil {
		return false
	}

	match := reg.MatchString(val)
	if !match {
		fmt.Println("not match error.")
		return false
	}

	var cnt int = 0          // 满足3中以上即可通过验证
	patternList := []string{ // 数字、大小写字母、特殊字符
		`[0-9]+`,
		`[a-z]+`,
		`[A-Z]+`,
		`[!@#$%^&*~-_+]+`,
	}
	for _, pattern := range patternList {
		match, _ = regexp.MatchString(pattern, val)
		if match {
			cnt++
		}
	}
	if cnt < 3 {
		fmt.Println("pwd should include at least 3 types.")
		return false
	}
	return true
}

func ValidatErr(u interface{}, err error) string {
	if err == nil { //如果为nil 说明校验通过
		return ""
	}
	invalid, ok := err.(*validator.InvalidValidationError) //如果是输入参数无效，则直接返回输入参数错误
	if ok {
		return "输入参数错误：" + invalid.Error()
	}
	validationErrs := err.(validator.ValidationErrors) //断言是ValidationErrors
	for _, validationErr := range validationErrs {

		fmt.Println(validationErr.Tag(), validationErr.ActualTag(), validationErr.Namespace(),
			validationErr.StructNamespace(), validationErr.Field(), validationErr.StructField(),
			validationErr.Value(), validationErr.Param(), validationErr.Kind(), validationErr.Type())
		fieldName := validationErr.Field() //获取是哪个字段不符合格式
		typeOf := reflect.TypeOf(u)
		// 如果是指针，获取其属性
		if typeOf.Kind() == reflect.Ptr {
			typeOf = typeOf.Elem()
		}
		field, ok := typeOf.FieldByName(fieldName) //通过反射获取filed
		if ok {
			errorInfo := field.Tag.Get("reg_error_info") // 获取field对应的reg_error_info tag值
			return fieldName + ":" + errorInfo           // 返回错误
		} else {
			return "缺失reg_error_info"
		}
	}
	return ""
}
