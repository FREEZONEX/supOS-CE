package service

import (
	"backend/internal/logic/supos/uns/uns/bo"
	"backend/internal/types"
)

type UnsMountService interface {
	ParseMountDetail(po bo.UnsInfo, simple bool) *types.MountDetailVo
}
