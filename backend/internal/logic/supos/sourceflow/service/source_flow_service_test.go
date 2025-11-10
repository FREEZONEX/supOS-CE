package service

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
		create: func(context.Context, *repo.NoderedSourceFlowRepo, string, *types.CreateTopicDto) error { return nil },
	}
	svc.repoFn = func(context.Context) *repo.NoderedSourceFlowRepo { return nil }

	var calls []string
	svc.create = func(ctx context.Context, _ *repo.NoderedSourceFlowRepo, tpl string, dto *types.CreateTopicDto) error {
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
	svc.create = func(ctx context.Context, _ *repo.NoderedSourceFlowRepo, tpl string, dto *types.CreateTopicDto) error {
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
	svc.repoFn = func(context.Context) *repo.NoderedSourceFlowRepo { return nil }

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
	service.repoFn = func(context.Context) *repo.NoderedSourceFlowRepo { return nil }
	var aliases []string
	service.create = func(_ context.Context, _ *repo.NoderedSourceFlowRepo, _ string, dto *types.CreateTopicDto) error {
		aliases = append(aliases, dto.GetAlias())
		return nil
	}

	spring.RegisterBeanNamed[*SourceFlowService]("testSourceFlowServicePublish", service)

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
