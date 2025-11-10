package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/logic/supos/flowcommon"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/clients/nodered/templates"
	"backend/share/spring"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

const mockTemplateName = "relational-emqx.json.tpl"

var (
	mockTemplateOnce sync.Once
	mockTemplate     string
	mockTemplateErr  error
)

// SourceFlowService listens to UNS lifecycle events and provisions Node-RED source flows automatically.
type SourceFlowService struct {
	log    logx.Logger
	svcCtx *svc.ServiceContext
	create func(context.Context, *dao.NoderedSourceFlowRepo, string, *types.CreateTopicDto) error
	repoFn func(context.Context) *dao.NoderedSourceFlowRepo
}

func init() {
	spring.RegisterLazy[*SourceFlowService](func() *SourceFlowService {
		svc := &SourceFlowService{
			log:    logx.WithContext(context.Background()),
			svcCtx: spring.GetBean[*svc.ServiceContext](),
		}
		svc.create = svc.createMockFlow
		svc.repoFn = func(ctx context.Context) *dao.NoderedSourceFlowRepo {
			return dao.NewNoderedSourceFlowRepo(ctx)
		}
		return svc
	})
}

// OnEventBatchCreateTableEvent consumes BatchCreateTableEvent and creates default source flows for UNS files
// that requested flow provisioning (AddFlow = true).
func (s *SourceFlowService) OnEventBatchCreateTableEvent(ev *event.BatchCreateTableEvent) error {
	if ev == nil {
		return nil
	}
	ctx := ev.Context
	if ctx == nil {
		ctx = context.Background()
	}

	files := ev.Creates[constants.PathTypeFile]
	if len(files) == 0 {
		return nil
	}

	tpl, err := loadMockTemplate()
	if err != nil {
		s.log.Errorf("load mock template failed: %v", err)
		return err
	}

	repoFactory := s.repoFn
	if repoFactory == nil {
		repoFactory = func(ctx context.Context) *dao.NoderedSourceFlowRepo {
			return dao.NewNoderedSourceFlowRepo(ctx)
		}
	}
	repo := repoFactory(ctx)
	var errs []error
	for _, dto := range files {
		if !shouldProvisionFlow(dto) {
			continue
		}
		creator := s.create
		if creator == nil {
			creator = s.createMockFlow
		}
		if err := creator(ctx, repo, tpl, dto); err != nil {
			s.log.Errorf("auto create source flow failed, alias=%s err=%v", strings.TrimSpace(dto.GetAlias()), err)
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func shouldProvisionFlow(dto *types.CreateTopicDto) bool {
	if dto == nil {
		return false
	}
	if dto.PathType != constants.PathTypeFile {
		return false
	}
	addFlow := dto.GetAddFlow()
	return addFlow != nil && *addFlow
}

func loadMockTemplate() (string, error) {
	mockTemplateOnce.Do(func() {
		mockTemplate, mockTemplateErr = templates.Load(mockTemplateName)
	})
	return mockTemplate, mockTemplateErr
}

func (s *SourceFlowService) createMockFlow(ctx context.Context, repo *dao.NoderedSourceFlowRepo, tpl string, dto *types.CreateTopicDto) error {
	if s.svcCtx == nil || s.svcCtx.SnowFlake == nil {
		return fmt.Errorf("service context not ready")
	}
	alias := strings.TrimSpace(dto.GetAlias())
	if alias == "" {
		alias = strings.TrimSpace(dto.Name)
	}
	if alias == "" {
		return fmt.Errorf("alias empty for topic id=%d", dto.GetId())
	}
	path := strings.TrimSpace(dto.GetPath())
	if path == "" {
		path = alias
	}

	// Skip if a flow already exists for this alias.
	existing, err := repo.SelectByAliases(ctx, []string{alias})
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		for _, flow := range existing {
			if flow != nil && strings.EqualFold(flow.Template, constants.FlowTypeNODERED) {
				s.log.Infof("skip auto mock flow for alias=%s, flow=%s already exists", alias, flow.FlowName)
				return nil
			}
		}
	}

	flowName, _, err := repo.FindAvailableFlowName(ctx, alias, constants.FlowTypeNODERED)
	if err != nil {
		return err
	}

	// Build payload content from fields
	payload := buildPayloadFromFields(dto.GetFields())

	// Render template with all supported $ variables
	rendered := templates.RenderDollar(tpl, map[string]string{
		"uns_path":         path,
		"model_alias":      alias,
		"alias_path_topic": path,
		"payload":          payload,
		"disabled":         "false",
		"clientid":         alias,
	}, flowcommon.GenerateNodeID)

	rec := &dao.NoderedSourceFlow{
		ID:          s.svcCtx.SnowFlake.GetSnowflakeId(),
		FlowName:    flowName,
		Description: fmt.Sprintf("auto mock flow for %s", alias),
		Template:    constants.FlowTypeNODERED,
		FlowStatus:  flowcommon.FlowStatusDraft,
		FlowData:    rendered,
	}

	if err := repo.Insert(ctx, rec); err != nil {
		return err
	}

	client := s.svcCtx.SourceNodeRed
	if client == nil {
		s.log.Infof("node-red client missing, skip deploy for flow %s", rec.FlowName)
		return nil
	}

	if _, err := flowcommon.DeployFlow(ctx, repo, rec.ID, rendered, client, flowcommon.ExtractAliases); err != nil {
		return err
	}

	// Refresh cached data after deployment.
	if _, err := repo.FindOne(ctx, rec.ID); err != nil {
		return err
	}
	return nil
}

// buildPayloadFromFields constructs the $payload snippet for the Node-RED function node
// based on topic field definitions. It maps common field types to random value
// generator functions inside the template function body.
func buildPayloadFromFields(fields []*types.FieldDefine) string {
	if len(fields) == 0 {
		return ""
	}
	parts := make([]string, 0, len(fields))
	for _, f := range fields {
		if f == nil || f.IsSystemField() {
			continue
		}
		t := strings.ToUpper(strings.TrimSpace(f.Type))
		fn := "randomString()"
		switch t {
		case "INTEGER", "LONG", "INT", "NUMBER":
			fn = "generateRandomNumber()"
		case "FLOAT", "DOUBLE":
			fn = "generateRandomFloatWithTwoDecimals()"
		case "BOOLEAN", "BOOL":
			fn = "getBool()"
		case "STRING", "TEXT":
			fn = "randomString()"
		default:
			// keep randomString()
		}
		parts = append(parts, fmt.Sprintf("'%s': %s", strings.TrimSpace(f.Name), fn))
	}
	if len(parts) == 0 {
		return ""
	}
	return "\n" + strings.Join(parts, ",\n")
}
