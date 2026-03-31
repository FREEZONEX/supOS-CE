package sourceflow

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

type flowAliasBindingRepo interface {
	FindAliasBindings(ctx context.Context, aliases []string, template string, excludeParentID int64) ([]*relationDB.FlowAliasBinding, error)
}

type CheckBindUNSLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckBindUNSLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckBindUNSLogic {
	return &CheckBindUNSLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckBindUNSLogic) CheckBindUNS(req *types.SourceFlowCheckBindUnsReq) error {
	if req == nil {
		return errors.Parameter.WithMsg("error.sys.parameterError")
	}
	flowID, err := parseSourceFlowEntityID(req.ID)
	if err != nil {
		return err
	}
	repo := relationDB.NewNoderedSourceFlowRepo(l.ctx)
	return validateSourceFlowUnsBindings(l.ctx, repo, flowID, req.UnsAliases, constants.FlowTypeNODERED)
}

func parseSourceFlowEntityID(id string) (int64, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return 0, errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	flowID, err := strconv.ParseInt(id, 10, 64)
	if err != nil || flowID <= 0 {
		return 0, errors.Parameter.WithMsg("nodered.flowId.empty")
	}
	return flowID, nil
}

func validateSourceFlowUnsBindings(ctx context.Context, repo flowAliasBindingRepo, currentFlowID int64, aliases []string, flowType string) error {
	cleanAliases := normalizeUnsAliases(aliases)
	if len(cleanAliases) == 0 {
		return nil
	}
	bindings, err := repo.FindAliasBindings(ctx, cleanAliases, flowType, currentFlowID)
	if err != nil {
		return err
	}
	if len(bindings) == 0 {
		return nil
	}
	return errors.Duplicate.WithMsg(buildUnsBindingConflictMsg(ctx, bindings))
}

func normalizeUnsAliases(aliases []string) []string {
	clean := make([]string, 0, len(aliases))
	seen := make(map[string]struct{}, len(aliases))
	for _, alias := range aliases {
		v := strings.TrimSpace(alias)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		clean = append(clean, v)
	}
	sort.Strings(clean)
	return clean
}

func buildUnsBindingConflictMsg(ctx context.Context, bindings []*relationDB.FlowAliasBinding) string {
	for _, binding := range bindings {
		if binding == nil {
			continue
		}
		topic := strings.TrimSpace(binding.Topic)
		flowName := strings.TrimSpace(binding.FlowName)
		if topic != "" && flowName != "" {
			return I18nUtils.GetMessageWithCtx(ctx, "nodered.topic.bound.flow", topic, flowName)
		}
		if topic != "" {
			return I18nUtils.GetMessageWithCtx(ctx, "nodered.topic.bound.flow.topicOnly", topic)
		}
	}
	return I18nUtils.GetMessageWithCtx(ctx, "nodered.topic.bound.flow.default")
}
