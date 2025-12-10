package service

import (
	"backend/internal/types"
)

type UnsMountService interface {
	ParseMountDetail(po types.UnsInfo, simple bool) *types.MountDetailVo
}
