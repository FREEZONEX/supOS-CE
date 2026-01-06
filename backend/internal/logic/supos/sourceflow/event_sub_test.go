package sourceflow

import (
	"backend/internal/common/constants"
	"backend/internal/common/event"
	repo "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestShouldProvisionFlow(t *testing.T) {
	trueVal := true
	falseVal := false
	tests := []struct {
		name string
		in   *types.CreateTopicDto
		exp  bool
	}{
		{"nil", nil, false},
		{"not file", &types.CreateTopicDto{PathType: constants.PathTypeDir, AddFlow: &trueVal}, false},
		{"no flag", &types.CreateTopicDto{PathType: constants.PathTypeFile}, false},
		{"flag false", &types.CreateTopicDto{PathType: constants.PathTypeFile, AddFlow: &falseVal}, false},
		{"flag true", &types.CreateTopicDto{PathType: constants.PathTypeFile, AddFlow: &trueVal}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.exp, shouldProvisionFlow(tt.in))
		})
	}
}

func TestSourceFlowService_OnEventBatchCreateTableEvent(t *testing.T) {
	svc := &SourceFlowService{
		log:    logx.WithContext(context.Background()),
		create: func(context.Context, sourceFlowRepo, string, *types.CreateTopicDto) error { return nil },
	}
	svc.repoFn = func(context.Context) sourceFlowRepo { return nil }

	var calls []string
	svc.create = func(ctx context.Context, _ sourceFlowRepo, tpl string, dto *types.CreateTopicDto) error {
		calls = append(calls, dto.GetAlias())
		require.NotEmpty(t, tpl)
		return nil
	}
	trueVal := true
	falseVal := false
	ev := &event.BatchCreateTableEvent{
		ApplicationEvent: event.ApplicationEvent{Context: context.Background()},
		Creates: map[int16][]*types.CreateTopicDto{
			constants.PathTypeFile: {
				{Alias: "mach1", Name: "mock", Path: "/path1", PathType: constants.PathTypeFile, AddFlow: &trueVal},
				{Alias: "mach2", Name: "skip", Path: "/path2", PathType: constants.PathTypeFile, AddFlow: &falseVal},
			},
			constants.PathTypeDir: {
				{Alias: "folder", Name: "folder", PathType: constants.PathTypeDir, AddFlow: &trueVal},
			},
		},
	}

	require.NoError(t, svc.OnEventBatchCreateTableEvent(ev))
	require.Equal(t, []string{"mach1"}, calls)
}

func TestSourceFlowService_OnEventAggregatesErrors(t *testing.T) {
	svc := &SourceFlowService{
		log: logx.WithContext(context.Background()),
	}
	errA := errors.New("a")
	errB := errors.New("b")
	order := 0
	svc.create = func(ctx context.Context, _ sourceFlowRepo, tpl string, dto *types.CreateTopicDto) error {
		order++
		if order == 1 {
			return errA
		}
		return errB
	}

	trueVal := true
	ev := &event.BatchCreateTableEvent{
		Creates: map[int16][]*types.CreateTopicDto{
			constants.PathTypeFile: {
				{Alias: "mach1", Name: "mock", Path: "/path1", PathType: constants.PathTypeFile, AddFlow: &trueVal},
				{Alias: "mach2", Name: "mock", Path: "/path2", PathType: constants.PathTypeFile, AddFlow: &trueVal},
			},
		},
	}
	svc.repoFn = func(context.Context) sourceFlowRepo { return nil }

	err := svc.OnEventBatchCreateTableEvent(ev)
	require.Error(t, err)
	require.ErrorIs(t, err, errA)
	require.ErrorIs(t, err, errB)
}

func TestSourceFlowService_PublishEventThroughSpring(t *testing.T) {
	ctx := context.Background()
	flag := true

	service := &SourceFlowService{
		log:    logx.WithContext(ctx),
		svcCtx: &svc.ServiceContext{},
	}
	service.repoFn = func(context.Context) sourceFlowRepo { return nil }
	var aliases []string
	service.create = func(_ context.Context, _ sourceFlowRepo, _ string, dto *types.CreateTopicDto) error {
		aliases = append(aliases, dto.GetAlias())
		return nil
	}

	wrapper := &testSourceFlowWrapper{SourceFlowService: service}
	spring.RegisterBeanNamed[*testSourceFlowWrapper]("testSourceFlowServicePublish", wrapper)

	err := spring.PublishEvent(&event.BatchCreateTableEvent{
		ApplicationEvent: event.ApplicationEvent{Context: ctx},
		Creates: map[int16][]*types.CreateTopicDto{
			constants.PathTypeFile: {
				{Alias: "busFlow", PathType: constants.PathTypeFile, AddFlow: &flag},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"busFlow"}, aliases)
}

func TestSourceFlowService_OnEventRemoveTopicsEvent_Skip(t *testing.T) {
	svc := &SourceFlowService{}
	svc.repoFn = func(context.Context) sourceFlowRepo {
		require.Fail(t, "repoFn should not be called")
		return nil
	}
	require.NoError(t, svc.OnEventRemoveTopicsEvent(nil))
	require.NoError(t, svc.OnEventRemoveTopicsEvent(&event.RemoveTopicsEvent{WithFlow: false}))
	require.NoError(t, svc.OnEventRemoveTopicsEvent(&event.RemoveTopicsEvent{WithFlow: true, Topics: []*types.CreateTopicDto{}}))
}

func TestSourceFlowService_OnEventRemoveTopicsEvent_Delete(t *testing.T) {
	fakeRepo := newFakeFlowRepo()
	top := &fakeTopRepo{}
	svcInstance := &SourceFlowService{
		log:    logx.WithContext(context.Background()),
		svcCtx: &svc.ServiceContext{},
	}
	svcInstance.repoFn = func(context.Context) sourceFlowRepo { return fakeRepo }
	svcInstance.topRepoFn = func(context.Context) flowTopRepo { return top }

	err := svcInstance.OnEventRemoveTopicsEvent(&event.RemoveTopicsEvent{
		ApplicationEvent: event.ApplicationEvent{Context: context.Background()},
		WithFlow:         true,
		Topics: []*types.CreateTopicDto{
			{Alias: "alpha"},
			{Alias: "beta"},
			{Alias: "alpha"},
		},
	})
	require.NoError(t, err)
	require.ElementsMatch(t, []int64{1, 2}, fakeRepo.replaced)
	require.ElementsMatch(t, []int64{1, 2}, fakeRepo.deleted)
	require.ElementsMatch(t, []int64{1, 2}, top.deleted)
}

func TestSourceFlowService_OnEventRemoveTopicsEvent_AggregatesErrors(t *testing.T) {
	fakeRepo := newFakeFlowRepo()
	svcInstance := &SourceFlowService{
		log: logx.WithContext(context.Background()),
	}
	svcInstance.repoFn = func(context.Context) sourceFlowRepo { return fakeRepo }
	svcInstance.topRepoFn = func(context.Context) flowTopRepo { return &fakeTopRepo{} }
	errA := errors.New("a")
	errB := errors.New("b")
	order := 0
	svcInstance.delete = func(context.Context, sourceFlowRepo, flowTopRepo, *repo.NoderedSourceFlow) error {
		order++
		if order == 1 {
			return errA
		}
		return errB
	}

	err := svcInstance.OnEventRemoveTopicsEvent(&event.RemoveTopicsEvent{
		WithFlow: true,
		Topics: []*types.CreateTopicDto{
			{Alias: "alpha"},
			{Alias: "beta"},
		},
	})
	require.Error(t, err)
	require.ErrorIs(t, err, errA)
	require.ErrorIs(t, err, errB)
}

type fakeFlowRepo struct {
	flows      map[int64]*repo.NoderedSourceFlow
	aliasIndex map[string][]int64
	replaced   []int64
	deleted    []int64
}

func newFakeFlowRepo() *fakeFlowRepo {
	return &fakeFlowRepo{
		flows: map[int64]*repo.NoderedSourceFlow{
			1: {ID: 1, FlowID: "fid-1", FlowName: "f1"},
			2: {ID: 2, FlowID: "fid-2", FlowName: "f2"},
		},
		aliasIndex: map[string][]int64{
			"alpha": {1, 2},
			"beta":  {2},
		},
	}
}

func (f *fakeFlowRepo) FindAvailableFlowName(ctx context.Context, base string, flowType string) (string, int, error) {
	return base, 0, nil
}

func (f *fakeFlowRepo) Insert(ctx context.Context, data *repo.NoderedSourceFlow) error {
	f.flows[data.ID] = data
	return nil
}

func (f *fakeFlowRepo) FindOne(ctx context.Context, id int64) (*repo.NoderedSourceFlow, error) {
	return f.flows[id], nil
}

func (f *fakeFlowRepo) Update(ctx context.Context, data *repo.NoderedSourceFlow) error {
	f.flows[data.ID] = data
	return nil
}

func (f *fakeFlowRepo) ReplaceModels(ctx context.Context, parentID int64, aliases []string) error {
	f.replaced = append(f.replaced, parentID)
	return nil
}

func (f *fakeFlowRepo) SelectByAliases(ctx context.Context, aliases []string) ([]*repo.NoderedSourceFlow, error) {
	seen := make(map[int64]struct{})
	var flows []*repo.NoderedSourceFlow
	for _, alias := range aliases {
		for _, id := range f.aliasIndex[alias] {
			if _, ok := seen[id]; ok {
				continue
			}
			if flow, ok := f.flows[id]; ok {
				seen[id] = struct{}{}
				flows = append(flows, flow)
			}
		}
	}
	return flows, nil
}

func (f *fakeFlowRepo) Delete(ctx context.Context, id int64) error {
	f.deleted = append(f.deleted, id)
	delete(f.flows, id)
	return nil
}

type fakeTopRepo struct {
	deleted []int64
}

func (f *fakeTopRepo) Delete(ctx context.Context, id int64) error {
	f.deleted = append(f.deleted, id)
	return nil
}

type testSourceFlowWrapper struct {
	*SourceFlowService
}
