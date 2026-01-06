package common

import (
	"github.com/bwmarrin/snowflake"
)

var snowFlake *snowflake.Node

func InitSnowflake(nodeId int64) {
	node, err := snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}
	snowFlake = node
}

// NextId 这里使用标准雪花算法，而 gitee.com/unitedrhino/share/utils.Snowflake 生成的格式不对，不能用
func NextId() int64 {
	return snowFlake.Generate().Int64()
}
