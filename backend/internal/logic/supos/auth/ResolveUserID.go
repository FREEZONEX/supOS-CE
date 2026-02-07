package auth

import (
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/vo"
	"context"
	"strings"
)

func ResolveUserID(ctx context.Context) string {
	if ctx == nil {
		return "guest"
	}
	if v := ctx.Value(apiutil.UserKey); v != nil {
		if user, ok := v.(*vo.UserInfoVo); ok && user != nil {
			if sub := strings.TrimSpace(user.Sub); sub != "" {
				return sub
			}
		}
	}
	return "guest"
}
func ResolveUserName(ctx context.Context) string {
	if ctx == nil {
		return "guest"
	}
	if v := ctx.Value(apiutil.UserKey); v != nil {
		if user, ok := v.(*vo.UserInfoVo); ok && user != nil {
			if userName := strings.TrimSpace(user.PreferredUsername); userName != "" {
				return userName
			}
		}
	}
	return "guest"
}
