package service

import (
	dao "backend/internal/repo/relationDB"
	"fmt"
	"strings"
	"testing"
)

func TestParseChildrenCount(t *testing.T) {

	data := `lay_rec                                                    |count_children|
-----------------------------------------------------------+--------------+
1992121008701575168/1992869679890173952                    |2:2           |
1992880505611096064/1992880505611096064/1992880505611096064|0:1           |
1992121008701575168/1992160505556373504                    |0:2           |
1992121008701575168/1992756402539466752                    |2:4           |
1992876268697161728                                        |2:6           |
1992121008701575168                                        |0:6           |
1992107863878668289/1992109352437157888                    |2:4           |
1992107863878668289                                        |0:4           |
1992107863878668289/1992111110131552256                    |0:2           |
1992107863878668289/1992111110131552256/1992111355687079937|0:1           |
1992107863878668289/1992108259405729792                    |2:7           |
1992121008701575168/1992162598539235328                    |0:3           |
1992121008701575168/1992162598539235328/1992162598539235329|0:2           |
1992878860768645120/1992878860768645120                    |0:1           |
1990971360238637057                                        |2:1           |`

	countChildrenList := convertToLayRecCc(data)

	// 打印结果
	for i, item := range countChildrenList {
		fmt.Printf("Item %d:,  %+v\n", i+1, item)
		fmt.Println()
	}

	ccMap := parseChildrenCount(countChildrenList)
	t.Logf("%+v\n", ccMap)
}
func convertToLayRecCc(data string) []*dao.LayRecCc {
	lines := strings.Split(data, "\n")
	var result []*dao.LayRecCc

	for _, line := range lines {
		// 跳过空行和表头行
		if strings.TrimSpace(line) == "" || strings.Contains(line, "lay_rec") || strings.Contains(line, "---+") {
			continue
		}

		// 按 | 分割行
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			layRec := strings.TrimSpace(parts[0])
			countChildren := strings.TrimSpace(parts[1])

			// 创建 LayRecCc 对象并添加到结果切片
			result = append(result, &dao.LayRecCc{
				LayRec:        layRec,
				CountChildren: countChildren,
			})
		}
	}

	return result
}
