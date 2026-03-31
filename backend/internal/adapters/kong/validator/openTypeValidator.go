package validator

import (
	"backend/internal/common/errors"
	"context"
)

// ValidateOpenType 校验菜单打开方式
func ValidateOpenType(ctx context.Context, openType int) error {
	// 0: iframe, 1: 新页面
	if openType == 0 || openType == 1 {
		return nil
	}
	return errors.BadRequest(ctx, "menu.opentype.invalid")
}
