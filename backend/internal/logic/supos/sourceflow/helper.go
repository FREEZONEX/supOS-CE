package sourceflow

import (
	"backend/internal/logic/supos/auth"
	"context"
	"strings"

	"backend/internal/common/utils/apiutil"
	"backend/internal/common/vo"
	"backend/internal/repo/relationDB"

	"gitee.com/unitedrhino/share/errors"
)

func LoadFlowByType(ctx context.Context, repo *relationDB.NoderedSourceFlowRepo, flowID int64, flowType string) (*relationDB.NoderedSourceFlow, error) {
	rec, err := repo.FindOne(ctx, flowID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, errors.NotFind.WithMsg("nodered.flow.not.exist")
	}
	ft := strings.TrimSpace(flowType)
	if ft != "" && !strings.EqualFold(strings.TrimSpace(rec.Template), ft) {
		return nil, errors.NotFind.WithMsg("nodered.flow.not.exist")
	}
	return rec, nil
}

func resolveUserID(ctx context.Context) string {
	return auth.ResolveUserID(ctx)
}

func resolveUser(ctx context.Context) *vo.UserInfoVo {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(apiutil.UserKey); v != nil {
		if user, ok := v.(*vo.UserInfoVo); ok && user != nil {
			return user
		}
	}
	return nil
}
