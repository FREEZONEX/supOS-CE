package common

import "gitee.com/unitedrhino/share/utils"

var SnowFlake *utils.SnowFlake

func NextId() int64 {
	return SnowFlake.GetSnowflakeId()
}
