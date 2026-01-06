package common

import (
	"fmt"
	"testing"
)
import "github.com/bwmarrin/snowflake"

func TestSnowId(t *testing.T) {
	fmt.Println("=== bwmarrin/snowflake 示例 ===")

	// 设置节点ID (0-1023)
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Printf("创建节点失败: %v\n", err)
		return
	}

	// 生成多个ID
	for i := 0; i < 5; i++ {
		id := node.Generate()
		fmt.Printf("Id: %d\n", id)
		fmt.Printf("  时间戳: %d\n", id.Time())
		fmt.Printf("  节点ID: %d\n", id.Node())
		fmt.Printf("  序列号: %d\n", id.Step())
		fmt.Printf("  字符串: %s\n", id.String())
		fmt.Println("---")
	}
}
