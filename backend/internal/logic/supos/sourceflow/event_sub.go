package sourceflow

import (
	"backend/internal/common"
	"backend/internal/common/constants"
	"backend/internal/common/event"
	"backend/internal/common/utils/apiutil"
	"backend/internal/logic/supos/flowcommon"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	noderedclient "backend/share/clients/nodered"
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

type sourceFlowRepo interface {
	FindAvailableFlowName(ctx context.Context, base string, flowType string) (string, int, error)
	Insert(ctx context.Context, data *dao.NoderedSourceFlow) error
	FindOne(ctx context.Context, id int64) (*dao.NoderedSourceFlow, error)
	Update(ctx context.Context, data *dao.NoderedSourceFlow) error
	ReplaceModels(ctx context.Context, parentID int64, aliases []string) error
	SelectByAliases(ctx context.Context, aliases []string) ([]*dao.NoderedSourceFlow, error)
	Delete(ctx context.Context, id int64) error
}

type flowTopRepo interface {
	Delete(ctx context.Context, id int64) error
}

// SourceFlowService listens to UNS lifecycle events and provisions Node-RED source flows automatically.
type SourceFlowService struct {
	log       logx.Logger
	svcCtx    *svc.ServiceContext
	create    func(context.Context, sourceFlowRepo, string, *types.CreateTopicDto) error
	delete    func(context.Context, sourceFlowRepo, flowTopRepo, *dao.NoderedSourceFlow) error
	repoFn    func(context.Context) sourceFlowRepo
	topRepoFn func(context.Context) flowTopRepo
}

func init() {
	fmt.Println("regSourceFlowService")
	spring.RegisterLazy[*SourceFlowService](func() *SourceFlowService {
		svc := &SourceFlowService{
			log:    logx.WithContext(context.Background()),
			svcCtx: spring.GetBean[*svc.ServiceContext](),
		}
		svc.create = svc.createMockFlow
		svc.delete = svc.deleteFlow
		svc.repoFn = func(ctx context.Context) sourceFlowRepo {
			return dao.NewNoderedSourceFlowRepo(ctx)
		}
		svc.topRepoFn = func(ctx context.Context) flowTopRepo {
			return dao.NewNoderedFlowTopRepo(ctx)
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
		repoFactory = func(ctx context.Context) sourceFlowRepo {
			return dao.NewNoderedSourceFlowRepo(ctx)
		}
	}
	repo := repoFactory(context.Background())
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

func (s *SourceFlowService) createMockFlow(ctx context.Context, repo sourceFlowRepo, tpl string, dto *types.CreateTopicDto) error {
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

	flowName, _, err := repo.FindAvailableFlowName(context.Background(), path, constants.FlowTypeNODERED)
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
	userCtx := apiutil.GetUserFromContext(ctx)
	userName := ""
	if userCtx != nil {
		userName = userCtx.PreferredUsername
	} else {
		logx.WithContext(ctx).Infof("username not found in context: %v", userCtx)
	}
	rec := &dao.NoderedSourceFlow{
		ID:         common.NextId(),
		FlowName:   flowName,
		Template:   constants.FlowTypeNODERED,
		FlowStatus: flowcommon.FlowStatusDraft,
		FlowData:   rendered,
		Creator:    userName,
	}
	ctx = context.Background()
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
		case "DATETIME":
			fn = "formatCurDate()"
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

// OnEventRemoveTopicsEvent cleans up Node-RED flows associated with deleted UNS topics.
func (s *SourceFlowService) OnEventRemoveTopicsEvent(ev *event.RemoveTopicsEvent) error {
	if ev == nil || !ev.WithFlow {
		return nil
	}
	ctx := context.Background()
	var aliases []string
	for _, t := range ev.Topics {
		if t == nil {
			continue
		}
		if alias := strings.TrimSpace(t.GetAlias()); alias != "" {
			aliases = append(aliases, alias)
		}
	}
	if len(aliases) == 0 {
		return nil
	}

	repoFactory := s.repoFn
	if repoFactory == nil {
		repoFactory = func(ctx context.Context) sourceFlowRepo {
			return dao.NewNoderedSourceFlowRepo(ctx)
		}
	}
	repo := repoFactory(ctx)
	if repo == nil {
		return fmt.Errorf("source flow repo not ready")
	}
	flows, err := repo.SelectByAliases(ctx, aliases)
	if err != nil {
		return err
	}
	if len(flows) == 0 {
		return nil
	}
	topRepoFactory := s.topRepoFn
	if topRepoFactory == nil {
		topRepoFactory = func(ctx context.Context) flowTopRepo {
			return dao.NewNoderedFlowTopRepo(ctx)
		}
	}
	topRepo := topRepoFactory(ctx)
	deleter := s.delete
	if deleter == nil {
		deleter = s.deleteFlow
	}

	seen := make(map[int64]struct{}, len(flows))
	var errs []error
	for _, f := range flows {
		if f == nil {
			continue
		}
		if _, ok := seen[f.ID]; ok {
			continue
		}
		seen[f.ID] = struct{}{}
		if err := deleter(ctx, repo, topRepo, f); err != nil {
			s.log.Errorf("delete flow on UNS removal failed, flowID=%d, name=%s, err=%v", f.ID, f.FlowName, err)
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (s *SourceFlowService) deleteFlow(ctx context.Context, repo sourceFlowRepo, topRepo flowTopRepo, flow *dao.NoderedSourceFlow) error {
	if repo == nil || flow == nil {
		return nil
	}
	if err := repo.ReplaceModels(ctx, flow.ID, nil); err != nil {
		return err
	}

	var client *noderedclient.Client
	if s != nil && s.svcCtx != nil {
		client = s.svcCtx.SourceNodeRed
	}
	if client != nil {
		if flowID := strings.TrimSpace(flow.FlowID); flowID != "" {
			var out map[string]any
			code, body, errs := client.DoJSON(ctx, "DELETE", "/flow/"+flowID, nil, &out)
			if len(errs) > 0 {
				return errors.Join(errs...)
			}
			if code != 200 && code != 204 && code != 404 {
				return fmt.Errorf("delete nodered flow %s failed, code=%d body=%s", flowID, code, string(body))
			}
		}
	} else {
		s.log.Infof("node-red client missing, skip runtime delete for flow %d", flow.ID)
	}

	if topRepo != nil {
		if err := topRepo.Delete(ctx, flow.ID); err != nil {
			return err
		}
	}
	return repo.Delete(ctx, flow.ID)
}
