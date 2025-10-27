package validator

import (
	"errors"
)

// ValidateOpenType 校验菜单打开方式
func ValidateOpenType(openType int) error {
	// 0: iframe, 1: 新页面
	if openType == 0 || openType == 1 {
		return nil
	}
	return errors.New("menu.opentype.invalid")
}
