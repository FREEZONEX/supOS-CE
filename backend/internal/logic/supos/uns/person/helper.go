package person

import (
	"context"
	"strings"

	"backend/internal/common/utils/apiutil"
	"backend/internal/common/vo"
)

func currentUser(ctx context.Context) *vo.UserInfoVo {
	if ctx == nil {
		return nil
	}
	if user, ok := ctx.Value(apiutil.UserKey).(*vo.UserInfoVo); ok && user != nil {
		return user
	}
	return nil
}

func currentUserID(ctx context.Context) string {
	if user := currentUser(ctx); user != nil {
		return strings.TrimSpace(user.Sub)
	}
	return ""
}
