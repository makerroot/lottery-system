package utils

import (
	"errors"
	"regexp"
	"unicode"
)

// 验证手机号格式（中国大陆）
func ValidatePhone(phone string) error {
	// 检查长度
	if len(phone) != 11 {
		return errors.New("手机号必须是11位数字")
	}

	// 检查是否全是数字
	matched, _ := regexp.MatchString(`^\d{11}$`, phone)
	if !matched {
		return errors.New("手机号格式不正确")
	}

	// 检查是否以1开头，第二位是3-9
	matched, _ = regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	if !matched {
		return errors.New("手机号号段不正确")
	}

	return nil
}

// 验证姓名（只允许中文、英文、空格）
func ValidateName(name string) error {
	// 检查长度
	if len(name) == 0 {
		return errors.New("姓名不能为空")
	}

	if len(name) > 50 {
		return errors.New("姓名长度不能超过50个字符")
	}

	// 检查是否包含非法字符
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.Is(unicode.Han, r) && r != ' ' {
			return errors.New("姓名只能包含中文、英文和空格")
		}
	}

	return nil
}

// 验证公司代码（只允许小写字母、数字、连字符）
func ValidateCompanyCode(code string) error {
	if code == "" {
		return errors.New("公司代码不能为空")
	}

	if len(code) > 50 {
		return errors.New("公司代码长度不能超过50个字符")
	}

	matched, _ := regexp.MatchString(`^[a-z0-9-]+$`, code)
	if !matched {
		return errors.New("公司代码只能包含小写字母、数字和连字符")
	}

	return nil
}

// 验证概率值（0-1之间）
func ValidateProbability(probability float64) error {
	if probability < 0 || probability > 1 {
		return errors.New("概率必须在0到1之间")
	}
	return nil
}

// 验证库存数量
func ValidateStock(total, used int) error {
	if total < 0 {
		return errors.New("总库存不能为负数")
	}

	if used < 0 {
		return errors.New("已使用库存不能为负数")
	}

	if used > total {
		return errors.New("已使用库存不能大于总库存")
	}

	return nil
}
