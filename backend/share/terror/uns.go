package terror

import "gitee.com/unitedrhino/share/errors"

const SysError = 3000000

var (
	// OK               = errors.NewCodeError(200, "成功")
	// Default          = errors.NewCodeError(SysError+1, "其他错误")
	UndoDuplicate     = errors.NewCodeError(SysError+10, "撤销时名称冲突")
	UnsNameDuplicate  = errors.NewCodeError(SysError+11, "名称冲突")
	FlowNameDuplicate = errors.NewCodeError(SysError+20, "Flow名称冲突")
)
