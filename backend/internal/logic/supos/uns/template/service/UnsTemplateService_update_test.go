package service

import (
	"testing"

	"backend/internal/common/constants"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
)

func TestBuildTemplateBaseInfoUpdateDtoUsesExistingAliasAndName(t *testing.T) {
	uns := &dao.UnsNamespace{
		Id:       1001,
		Alias:    "tpl_alias",
		Name:     "Template Name",
		PathType: constants.PathTypeTemplate,
	}

	dto := buildTemplateBaseInfoUpdateDto(uns, &types.UpdateTemplateBaseInfoReq{
		ID:          uns.Id,
		Description: "updated description",
	})

	if dto.Alias != uns.Alias {
		t.Fatalf("expected alias %q, got %q", uns.Alias, dto.Alias)
	}
	if dto.Name != uns.Name {
		t.Fatalf("expected name %q, got %q", uns.Name, dto.Name)
	}
	if dto.Description == nil || *dto.Description != "updated description" {
		t.Fatalf("expected description to be updated, got %#v", dto.Description)
	}
}

func TestBuildTemplateSubscribeUpdateDtoUsesExistingAliasAndName(t *testing.T) {
	uns := &dao.UnsNamespace{
		Id:       1002,
		Alias:    "tpl_alias",
		Name:     "Template Name",
		PathType: constants.PathTypeTemplate,
	}

	dto := buildTemplateSubscribeUpdateDto(uns, &types.UpdateTemplateSubscribeReq{
		ID:                 uns.Id,
		SubscribeEnable:    "true",
		SubscribeFrequency: "5s",
	})

	if dto.Alias != uns.Alias {
		t.Fatalf("expected alias %q, got %q", uns.Alias, dto.Alias)
	}
	if dto.Name != uns.Name {
		t.Fatalf("expected name %q, got %q", uns.Name, dto.Name)
	}
	if dto.SubscribeEnable == nil || !*dto.SubscribeEnable {
		t.Fatalf("expected subscribe enable to be true, got %#v", dto.SubscribeEnable)
	}
	if dto.Frequency != "5s" {
		t.Fatalf("expected frequency to be 5s, got %q", dto.Frequency)
	}
}

func TestResolveTemplateDescriptionSupportsBothFields(t *testing.T) {
	modelDescription := "from-model-description"
	req := &types.UpdateTemplateFieldsAndDescReq{
		ModelDescription: &modelDescription,
	}

	if got := req.ResolveDescription(); got == nil || *got != modelDescription {
		t.Fatalf("expected modelDescription fallback, got %#v", got)
	}

	description := "from-description"
	req.Description = &description
	if got := req.ResolveDescription(); got == nil || *got != description {
		t.Fatalf("expected description to take precedence, got %#v", got)
	}
}
