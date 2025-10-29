package common

/**

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"backend/internal/common/constants"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	shareerrors "gitee.com/unitedrhino/share/errors"
)

const (
	templateRootID    int64  = 1
	templateRootAlias string = "__templates__"
	templateRootPath  string = "tmplt"
)

type CreateUNSResult struct {
	Namespace *relationDB.UnsNamespace
}

// CreateTemplate 如果需要模板，需先调用该函数，返回模板ID写回 dto.ModelId
func CreateTemplate(ctx context.Context, svcCtx *svc.ServiceContext, req *types.UnsCreateTopicDTO) (*relationDB.UnsNamespace, error) {
	templateName := strings.TrimSpace(req.Name)
	if templateName == "" {
		return nil, shareerrors.Parameter.WithMsg("模板名称不能为空")
	}
	if len(req.Fields) == 0 {
		return nil, shareerrors.Parameter.WithMsg("模板字段不能为空")
	}

	repo := relationDB.NewUnsNamespaceRepo()

	aliasSeed := strings.TrimSpace(req.Alias)
	pathForAlias := fmt.Sprintf("%s/%s", templateRootPath, templateName)
	if aliasSeed == "" {
		aliasSeed = generateAliasSeed(svcCtx, templateName, pathForAlias, constants.PathTypeTemplate)
	}
	templateAlias, err := ensureUniqueAlias(ctx, repo, aliasSeed)
	if err != nil {
		return nil, err
	}

	template := &relationDB.UnsNamespace{
		Id:           svcCtx.SnowFlake.GetSnowflakeId(),
		Name:         templateName,
		DisplayName:  templateName,
		Alias_:       templateAlias,
		ParentAlias:  templateRootAlias,
		ParentId:     templateRootID,
		PathType:     constants.PathTypeTemplate,
		DataType:     0,
		Path:         fmt.Sprintf("%s/%s", templateRootPath, templateAlias),
		LayRec:       fmt.Sprintf("%d/%s", templateRootID, templateAlias),
		Description:  strings.TrimSpace(req.Description),
		Extend:       req.Extend,
		NumberFields: int16(len(req.Fields)),
		Status:       1,
	}

	if err = repo.Insert(ctx, template); err != nil {
		return nil, err
	}

	return template, nil
}

// CreateUNSBatch 批量创建 UNS 实体（不包含模板创建）
func CreateUNSBatch(ctx context.Context, svcCtx *svc.ServiceContext, reqs []*types.UnsCreateTopicDTO) ([]*relationDB.UnsNamespace, error) {
	if len(reqs) == 0 {
		return nil, nil
	}
	repo := relationDB.NewUnsNamespaceRepo(ctx)
	results := make([]*relationDB.UnsNamespace, 0, len(reqs))
	for _, req := range reqs {
		ns, err := createSingle(ctx, svcCtx, repo, req)
		if err != nil {
			return nil, err
		}
		results = append(results, ns)
	}
	return results, nil
}

func createSingle(ctx context.Context, svcCtx *svc.ServiceContext, repo *relationDB.UnsNamespaceRepo, req *types.UnsCreateTopicDTO) (*relationDB.UnsNamespace, error) {
	if req == nil {
		return nil, shareerrors.Parameter.WithMsg("请求参数为空")
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, shareerrors.Parameter.WithMsg("名称不能为空")
	}

	if req.PathType != int64(constants.PathTypeDir) && req.PathType != int64(constants.PathTypeFile) {
		return nil, shareerrors.Parameter.WithMsg("pathType仅支持目录或文件")
	}

	if req.DataType != 0 &&
		req.DataType != int64(constants.TimeSequenceType) &&
		req.DataType != int64(constants.RelationType) {
		return nil, shareerrors.Parameter.WithMsg("dataType只支持1或2")
	}

	var (
		parent      *relationDB.UnsNamespace
		parentAlias = strings.TrimSpace(req.ParentAlias)
		err         error
	)

	if req.ParentId != 0 {
		parent, err = repo.FindOne(ctx, req.ParentId)
		if err != nil {
			if shareerrors.Is(err, shareerrors.NotFind) {
				return nil, shareerrors.Parameter.WithMsg("父节点不存在")
			}
			return nil, err
		}
		parentAlias = parent.Alias_
	} else if parentAlias != "" {
		parent, err = repo.FindOneByAlias(ctx, parentAlias)
		if err != nil {
			if shareerrors.Is(err, shareerrors.NotFind) {
				return nil, shareerrors.Parameter.WithMsg("父节点不存在")
			}
			return nil, err
		}
		if parent != nil {
			req.ParentId = parent.Id
		}
	}

	path := buildPath(name, parent, req.PathType)

	aliasSeed := strings.TrimSpace(req.Alias)
	if aliasSeed == "" {
		aliasSeed = generateAliasSeed(svcCtx, name, path, req.PathType)
	}
	alias, err := ensureUniqueAlias(ctx, repo, aliasSeed)
	if err != nil {
		return nil, err
	}
	layRec := buildLayRec(alias, parent)

	dataType16, err := toInt16(req.DataType)
	if err != nil {
		return nil, shareerrors.Parameter.WithMsg("dataType超出范围")
	}
	dataSrcID16, err := toInt16(req.DataSrcId)
	if err != nil {
		return nil, shareerrors.Parameter.WithMsg("dataSrcId越界")
	}
	numberFields16, err := toInt16(int64(len(req.Fields)))
	if err != nil {
		return nil, shareerrors.Parameter.WithMsg("字段数量过多")
	}

	withFlags := calcUnsFlags(req)
	entity := &relationDB.UnsNamespace{
		Id:           svcCtx.SnowFlake.GetSnowflakeId(),
		Name:         name,
		DisplayName:  name,
		Alias_:       alias,
		ParentAlias:  parentAlias,
		ParentId:     req.ParentId,
		PathType:     int16(req.PathType),
		DataType:     int16(dataType16),
		Path:         path,
		LayRec:       layRec,
		Description:  strings.TrimSpace(req.Description),
		Protocol:     strings.TrimSpace(req.Protocol),
		ProtocolType: strings.TrimSpace(req.ProtocolType),
		DataSrcId:    int16(dataSrcID16),
		Extend:       req.Extend,
		NumberFields: int16(numberFields16),
		WithFlags:    withFlags,
		Status:       1,
		ModelId:      req.ModelId,
	}

	if parent != nil {
		entity.ParentId = parent.Id
		entity.ParentAlias = parent.Alias_
	}

	if err = repo.Insert(ctx, entity); err != nil {
		return nil, err
	}

	return entity, nil
}

func generateAliasSeed(svcCtx *svc.ServiceContext, name, path string, pathType int64) string {
	source := strings.TrimSpace(path)
	if source == "" {
		source = strings.TrimSpace(name)
	}
	source = strings.Trim(source, "/")
	if source == "" {
		source = "node"
	}
	if pathType == int64(constants.PathTypeDir) || pathType == int64(constants.PathTypeTemplate) {
		if idx := strings.LastIndex(source, "/"); idx >= 0 {
			source = source[idx+1:]
		}
	}

	source = strings.ReplaceAll(source, "/", "_")
	source = strings.ReplaceAll(source, "-", "_")
	base := normalizeAliasBase(source)
	if base == "" {
		base = "node"
	}
	if len(base) > 20 {
		base = base[:20]
	}
	suffix := strings.ToLower(strconv.FormatInt(svcCtx.SnowFlake.GetSnowflakeId(), 36))
	if len(suffix) > 20 {
		suffix = suffix[:20]
	}
	alias := base + "_" + suffix
	if alias[0] < 'a' || alias[0] > 'z' {
		alias = "a_" + alias
	}
	return alias
}

func ensureUniqueAlias(ctx context.Context, repo *relationDB.UnsNamespaceRepo, seed string) (string, error) {
	candidate := sanitizeAlias(seed)
	if candidate == "" {
		candidate = "topic"
	}

	base, suffix, hasNumeric := splitAliasSuffix(candidate)
	searchBase := candidate
	if hasNumeric {
		searchBase = base
	}

	aliases, err := repo.ListAliasByBase(ctx, searchBase)
	if err != nil {
		return "", err
	}

	if len(aliases) == 0 {
		return candidate, nil
	}

	aliasExists := false
	maxSuffix := -1
	lowerCandidate := strings.ToLower(candidate)
	lowerBase := strings.ToLower(searchBase)

	for _, item := range aliases {
		l := strings.ToLower(strings.TrimSpace(item))
		if l == lowerCandidate {
			aliasExists = true
			if hasNumeric && suffix > maxSuffix {
				maxSuffix = suffix
			} else if !hasNumeric && maxSuffix < 0 {
				maxSuffix = 0
			}
		}
		if b, idx, ok := splitAliasSuffix(l); ok && b == lowerBase {
			if idx > maxSuffix {
				maxSuffix = idx
			}
		} else if l == lowerBase && maxSuffix < 0 {
			maxSuffix = 0
		}
	}

	if !aliasExists {
		return candidate, nil
	}

	if hasNumeric && suffix > maxSuffix {
		maxSuffix = suffix
	}
	if maxSuffix < 0 {
		maxSuffix = 0
	}

	return fmt.Sprintf("%s-%d", searchBase, maxSuffix+1), nil
}

func sanitizeAlias(seed string) string {
	base := normalizeAliasBase(seed)
	if base == "" {
		base = "topic"
	}
	if base[0] >= '0' && base[0] <= '9' {
		base = "a_" + base
	}
	return base
}

func normalizeAliasBase(seed string) string {
	seed = strings.TrimSpace(seed)
	if seed == "" {
		return ""
	}
	var builder strings.Builder
	for _, r := range seed {
		switch {
		case unicode.IsLetter(r):
			builder.WriteRune(unicode.ToLower(r))
		case unicode.IsDigit(r):
			builder.WriteRune(r)
		case r == '-' || r == '_':
			builder.WriteRune('_')
		}
	}
	result := builder.String()
	return strings.Trim(result, "_")
}

func splitAliasSuffix(alias string) (string, int, bool) {
	idx := strings.LastIndex(alias, "-")
	if idx <= 0 || idx == len(alias)-1 {
		return alias, -1, false
	}
	numPart := alias[idx+1:]
	n, err := strconv.Atoi(numPart)
	if err != nil {
		return alias, -1, false
	}
	return alias[:idx], n, true
}

func buildPath(name string, parent *relationDB.UnsNamespace, pathType int64) string {
	cleanName := strings.Trim(strings.TrimSpace(name), "/")
	parentPath := ""
	if parent != nil {
		parentPath = strings.Trim(parent.Path, "/")
	}

	var path string
	switch {
	case parentPath != "" && cleanName != "":
		path = parentPath + "/" + cleanName
	case parentPath != "":
		path = parentPath
	default:
		path = cleanName
	}
	return path
}

func buildLayRec(alias string, parent *relationDB.UnsNamespace) string {
	cleanAlias := strings.Trim(alias, "/")
	if parent == nil {
		return cleanAlias
	}
	parentLayRec := strings.Trim(parent.LayRec, "/")
	if parentLayRec == "" {
		return cleanAlias
	}
	return parentLayRec + "/" + cleanAlias
}

func toInt16(id int64) (int16, error) {
	if id == 0 {
		return 0, nil
	}
	const (
		maxInt16 = 1<<15 - 1
		minInt16 = -1 << 15
	)
	if id > maxInt16 || id < minInt16 {
		return 0, fmt.Errorf("out of int16 range")
	}
	return int16(id), nil
}

func calcUnsFlags(req *types.UnsCreateTopicDTO) int32 {
	flags := 0
	if req.AddFlow {
		flags |= constants.UnsFlagWithFlow
	}
	if req.AddDashBoard {
		flags |= constants.UnsFlagWithDashboard
	}
	if req.Save2db {
		flags |= constants.UnsFlagWithSave2DB
	}
	return int32(flags)
}
*/
