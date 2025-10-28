package service

import (
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
)

type UnsMountService interface {
	ParseMountDetail(po *dao.UnsNamespace, simple bool) *types.MountDetailVo
}
