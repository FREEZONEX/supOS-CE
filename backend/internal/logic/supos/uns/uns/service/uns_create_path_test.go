package service

import (
	"backend/internal/common/constants"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
	"strings"
	"testing"
)

func TestIsAbsoluteCreatePathName(t *testing.T) {
	if !isAbsoluteCreatePathName("a/b") {
		t.Fatal("expected multi-segment name to be treated as absolute path")
	}
	if isAbsoluteCreatePathName("single") {
		t.Fatal("expected single segment name to stay relative")
	}
}

func TestValidateCreatePathSegments(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name     string
		pathType int16
		wantErr  bool
	}{
		{name: "device/Metric/temp", pathType: constants.PathTypeFile, wantErr: false},
		{name: "/device/Metric/temp", pathType: constants.PathTypeFile, wantErr: true},
		{name: "device/Metric/temp/", pathType: constants.PathTypeFile, wantErr: true},
		{name: "device//Metric/temp", pathType: constants.PathTypeFile, wantErr: true},
		{name: "device/template/temp", pathType: constants.PathTypeFile, wantErr: true},
		{name: strings.Repeat("a", 64) + "/Metric/temp", pathType: constants.PathTypeFile, wantErr: true},
		{name: "device/folder/path", pathType: constants.PathTypeDir, wantErr: false},
		{name: "device/label/path", pathType: constants.PathTypeDir, wantErr: true},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			segments, msg := validateCreatePathSegments(ctx, tt.name, tt.pathType)
			if tt.wantErr {
				if msg == "" {
					t.Fatalf("expected validation error, got success: %#v", segments)
				}
				return
			}
			if msg != "" {
				t.Fatalf("expected success, got error: %s", msg)
			}
			if len(segments) == 0 {
				t.Fatal("expected non-empty segments")
			}
		})
	}
}

func TestValidateCreateParent(t *testing.T) {
	ctx := context.Background()
	categoryType := int16(3)

	msg := validateCreateParent(ctx, &types.CreateTopicDto{PathType: constants.PathTypeDir}, &dao.UnsNamespace{
		PathType: constants.PathTypeDir,
		DataType: &categoryType,
	})
	if msg == "" {
		t.Fatal("expected folder creation under category folder to be rejected")
	}

	msg = validateCreateParent(ctx, &types.CreateTopicDto{PathType: constants.PathTypeFile}, &dao.UnsNamespace{
		PathType: constants.PathTypeDir,
		DataType: &categoryType,
	})
	if msg != "" {
		t.Fatalf("expected file creation under category folder to be allowed, got %s", msg)
	}
}

func TestTryReuseExistingCategoryDirLeaf(t *testing.T) {
	metricType := int16(3)
	parentID := int64(101)
	existing := &dao.UnsNamespace{
		Id:       202,
		Name:     "Metric",
		Alias:    "metric_parent",
		ParentId: &parentID,
		PathType: constants.PathTypeDir,
		DataType: &metricType,
	}
	dto := &types.CreateTopicDto{
		Name:     "Metric",
		PathType: constants.PathTypeDir,
	}
	if !tryReuseExistingCategoryDirLeaf(dto, existing) {
		t.Fatal("expected existing category dir to be reused")
	}
	if dto.Id != existing.Id || dto.Alias != existing.Alias {
		t.Fatalf("expected dto to reuse existing dir identity, got id=%d alias=%s", dto.Id, dto.Alias)
	}

	normalDir := &dao.UnsNamespace{
		Id:       303,
		Name:     "Metric",
		Alias:    "metric_normal",
		PathType: constants.PathTypeDir,
	}
	dto = &types.CreateTopicDto{
		Name:     "Metric",
		PathType: constants.PathTypeDir,
	}
	if tryReuseExistingCategoryDirLeaf(dto, normalDir) {
		t.Fatal("did not expect normal folder named Metric to be reused as category dir")
	}
}

func TestValidateMultiSegmentFilePath(t *testing.T) {
	ctx := context.Background()
	jsonbType := constants.JsonbType
	stateType := int16(1)

	parentDataType, msg := validateMultiSegmentFilePath(
		ctx,
		[]string{"device", "State", "temperature"},
		nil,
		&jsonbType,
	)
	if msg != "" {
		t.Fatalf("expected success, got error: %s", msg)
	}
	if parentDataType == nil || *parentDataType != 1 {
		t.Fatalf("expected parentDataType 1, got %#v", parentDataType)
	}

	_, msg = validateMultiSegmentFilePath(
		ctx,
		[]string{"device", "Metric", "temperature"},
		nil,
		&jsonbType,
	)
	if msg == "" {
		t.Fatal("expected mismatch error for Metric + JsonbType")
	}

	_, msg = validateMultiSegmentFilePath(
		ctx,
		[]string{"device", "custom", "temperature"},
		nil,
		&jsonbType,
	)
	if msg == "" {
		t.Fatal("expected category segment validation error")
	}

	_, msg = validateMultiSegmentFilePath(
		ctx,
		[]string{"device", "Action", "temperature"},
		&stateType,
		&jsonbType,
	)
	if msg == "" {
		t.Fatal("expected parentDataType mismatch error")
	}
}
