package service

import (
	"backend/internal/adapters/postgresql"
	"backend/internal/adapters/timescaledb"
	"backend/internal/common/I18nUtils"
	"backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"strings"
)

type HistoryQueryService struct {
	unsMapper relationDB.UnsNamespaceRepo
}

func NewHistoryQueryService() *HistoryQueryService {
	return &HistoryQueryService{}
}

func (s *HistoryQueryService) Query(ctx context.Context, req *types.HistoryValueRequest) (*types.UnsHistoryQueryData, error) {
	aliasList := normalizeHistoryAliasList(req.AliasList)
	db := relationDB.GetDb(ctx)
	unsList, err := s.unsMapper.ListByAlias(db, aliasList)
	if err != nil {
		return nil, err
	}

	unsMap := make(map[string]*relationDB.UnsNamespace, len(unsList))
	for _, item := range unsList {
		if item == nil {
			continue
		}
		unsMap[item.Alias] = item
	}

	data := &types.UnsHistoryQueryData{
		Results:     make([]types.UnsHistoryFileResult, 0, len(aliasList)),
		NotExists:   make([]string, 0),
		ErrorFields: make(map[string]string),
	}

	for _, alias := range aliasList {
		uns := unsMap[alias]
		if uns == nil {
			data.NotExists = append(data.NotExists, alias)
			continue
		}
		result, err := s.queryAlias(ctx, uns, req)
		if err != nil {
			data.ErrorFields[alias] = err.Error()
			continue
		}
		data.Results = append(data.Results, *result)
	}

	if len(data.NotExists) == 0 {
		data.NotExists = nil
	}
	if len(data.ErrorFields) == 0 {
		data.ErrorFields = nil
	}

	return data, nil
}

func (s *HistoryQueryService) queryAlias(ctx context.Context, uns *relationDB.UnsNamespace, req *types.HistoryValueRequest) (*types.UnsHistoryFileResult, error) {
	if strings.TrimSpace(uns.GetTimestampField()) == "" {
		return nil, localizedHistoryMsg(ctx, "uns.file.history.query.timestamp.not.found")
	}

	switch uns.GetSrcJdbcType() {
	case types.SrcJdbcTypePostgresql:
		svc := spring.GetBean[*postgresql.PgPersistentService]()
		if svc == nil {
			return nil, localizedHistoryMsg(ctx, "uns.file.history.query.service.unavailable")
		}
		return svc.QueryHistory(ctx, uns, req)
	case types.SrcJdbcTypeTimeScaleDB:
		svc := spring.GetBean[*timescaledb.TsdbPersistentService]()
		if svc == nil {
			return nil, localizedHistoryMsg(ctx, "uns.file.history.query.service.unavailable")
		}
		return svc.QueryHistory(ctx, uns, req)
	default:
		return nil, localizedHistoryMsg(ctx, "uns.file.history.query.datasource.unsupported")
	}
}

func normalizeHistoryAliasList(aliasList []string) []string {
	if len(aliasList) == 0 {
		return nil
	}
	result := make([]string, 0, len(aliasList))
	seen := make(map[string]struct{}, len(aliasList))
	for _, alias := range aliasList {
		alias = strings.TrimSpace(alias)
		if alias == "" {
			continue
		}
		if _, ok := seen[alias]; ok {
			continue
		}
		seen[alias] = struct{}{}
		result = append(result, alias)
	}
	return result
}

func localizedHistoryMsg(ctx context.Context, key string) error {
	return historyQueryError(I18nUtils.GetMessageWithCtx(ctx, key))
}

type historyQueryError string

func (e historyQueryError) Error() string {
	return string(e)
}
