package service

import (
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	"backend/internal/types"
	"backend/share/base/buffer"
	"encoding/json"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
)

func TestDecodeStreamedJson(t *testing.T) {
	v1Alias := PathUtil.GenerateFileAlias("v1")
	t.Log("v1Alias:", v1Alias)
	countNode := 0
	decodedPath := make(map[string]bool)
	src := buffer.NewBuffer([]byte(__realJson))
	er := jsonstream.DecodeStreamedJson(src, 64, 100,
		nodeGetChildren, func(node, parent *FileData) *types.CreateTopicDto {
			countNode++
			rs := node2vo(node, parent)

			p := node.getPath()
			decodedPath["/"+p] = true
			return rs
		}, func(readSize int64, propName string, nodes []*types.CreateTopicDto) {
			jsonBs, _ := json.Marshal(nodes)
			t.Logf("readSize: %d %s[%d]: %s", readSize, propName, len(nodes), string(jsonBs))
		})
	if er != nil {
		t.Fatal(er)
	}
	t.Log("countNode:", countNode)
	vs, _, _, _ := jsonparser.Get([]byte(__realJson), "UNS")
	paths, er := GetAllPathsOptimized(string(vs))
	if er != nil {
		t.Fatal(er)
	}
	t.Log(len(paths), "paths:", paths)
	for _, path := range paths {
		if !decodedPath[path] {
			t.Log("缺少：", path)
		}
	}
}

func TestParsePath(t *testing.T) {

}

// TreeNode 定义树节点结构
type TreeNode struct {
	Name     string     `json:"name"`
	Children []TreeNode `json:"children,omitempty"`
}

// 修正后的优化版本 - 使用切片来管理路径
func GetAllPathsOptimized(jsonData string) ([]string, error) {
	var nodes []TreeNode
	if err := json.Unmarshal([]byte(jsonData), &nodes); err != nil {
		return nil, err
	}

	paths := make([]string, 0)
	currentPath := make([]string, 0)

	var traverse func(node *TreeNode)
	traverse = func(node *TreeNode) {
		// 添加当前节点到路径
		currentPath = append(currentPath, node.Name)

		// 构建完整路径并记录
		fullPath := "/" + strings.Join(currentPath, "/")
		paths = append(paths, fullPath)

		// 遍历子节点
		for i := range node.Children {
			traverse(&node.Children[i])
		}

		// 回溯：移除当前节点
		currentPath = currentPath[:len(currentPath)-1]
	}

	for i := range nodes {
		traverse(&nodes[i])
	}

	return paths, nil
}

// 使用 strings.Builder 的正确版本
func GetAllPathsWithBuilder(jsonData string) ([]string, error) {
	var nodes []TreeNode
	if err := json.Unmarshal([]byte(jsonData), &nodes); err != nil {
		return nil, err
	}

	paths := make([]string, 0)

	var traverse func(node *TreeNode, prefix string)
	traverse = func(node *TreeNode, prefix string) {
		// 构建当前路径
		var currentPath string
		if prefix == "" {
			currentPath = "/" + node.Name
		} else {
			currentPath = prefix + "/" + node.Name
		}

		// 记录路径
		paths = append(paths, currentPath)

		// 遍历子节点
		for i := range node.Children {
			traverse(&node.Children[i], currentPath)
		}
	}

	for i := range nodes {
		traverse(&nodes[i], "")
	}

	return paths, nil
}

// 性能更好的版本 - 预分配内存
func GetAllPathsEfficient(jsonData string) ([]string, error) {
	var nodes []TreeNode
	if err := json.Unmarshal([]byte(jsonData), &nodes); err != nil {
		return nil, err
	}

	// 预估计容量以减少内存分配
	paths := make([]string, 0, estimateCapacity(nodes))
	currentPath := make([]string, 0, 10) // 假设路径深度不超过10

	var traverse func(node *TreeNode)
	traverse = func(node *TreeNode) {
		currentPath = append(currentPath, node.Name)
		paths = append(paths, "/"+strings.Join(currentPath, "/"))

		for i := range node.Children {
			traverse(&node.Children[i])
		}

		currentPath = currentPath[:len(currentPath)-1]
	}

	for i := range nodes {
		traverse(&nodes[i])
	}

	return paths, nil
}

// 估计需要的容量
func estimateCapacity(nodes []TreeNode) int {
	if len(nodes) == 0 {
		return 0
	}

	count := 0
	for _, node := range nodes {
		count += countNodes(&node)
	}
	return count
}

func countNodes(node *TreeNode) int {
	if node == nil {
		return 0
	}

	count := 1 // 当前节点
	for i := range node.Children {
		count += countNodes(&node.Children[i])
	}
	return count
}

var __realJson = `
{
  "Template": [],
  "Label": [],
  "UNS": [
    {
      "name": "v1",
      "type": "folder",
      "children": [
        {
          "name": "Plant_Name",
          "type": "folder",
          "children": [
            {
              "name": "SMT-Area-1",
              "type": "folder",
              "children": [
                {
                  "name": "SMT-Line-1",
                  "type": "folder",
                  "children": [
                    {
                      "name": "Printer-Cell",
                      "type": "folder",
                      "children": [
                        {
                          "name": "Printer01",
                          "type": "folder",
                          "children": [
                            {
                              "name": "State",
                              "type": "folder",
                              "topicType": "STATE",
                              "children": [
                                {
                                  "name": "current_job",
                                  "type": "file",
                                  "topicType": "STATE",
                                  "dataType": "RELATION_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "job_id",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "product_id",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "planned_quantity",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "completed_quantity",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "status",
                                      "type": "LONG"
                                    }
                                  ]
                                },
                                {
                                  "name": "alarm_status",
                                  "type": "file",
                                  "topicType": "STATE",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE"
                                }
                              ]
                            },
                            {
                              "name": "Action",
                              "type": "folder",
                              "topicType": "ACTION",
                              "children": [
                                {
                                  "name": "start_job",
                                  "type": "file",
                                  "topicType": "ACTION",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "FALSE",
                                  "mockData": "FALSE"
                                },
                                {
                                  "name": "stop_job",
                                  "type": "file",
                                  "topicType": "ACTION",
                                  "dataType": "JSONB_TYPE",
                                  "generateDashboard": "FALSE",
                                  "enableHistory": "FALSE",
                                  "mockData": "FALSE"
                                }
                              ]
                            },
                            {
                              "name": "Metric",
                              "type": "folder",
                              "topicType": "METRIC",
                              "children": [
                                {
                                  "name": "board_cycle_time",
                                  "type": "file",
                                  "topicType": "METRIC",
                                  "dataType": "TIME_SEQUENCE_TYPE",
                                  "generateDashboard": "TRUE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "cycle_time_ms",
                                      "type": "LONG"
                                    }
                                  ]
                                },
                                {
                                  "name": "boards_count",
                                  "type": "file",
                                  "topicType": "METRIC",
                                  "dataType": "TIME_SEQUENCE_TYPE",
                                  "generateDashboard": "TRUE",
                                  "enableHistory": "TRUE",
                                  "mockData": "FALSE",
                                  "fields": [
                                    {
                                      "name": "good_count",
                                      "type": "LONG"
                                    },
                                    {
                                      "name": "ng_count",
                                      "type": "LONG"
                                    }
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
`
var __unsJson = `[
                    {
                      "alias": "wenjianjiaA_686e7fad643b4db0b11e",
                      "displayName": "指标",
                      "type": "folder",
                      "name": "指标",
                      "children": [
                        {
                          "alias": "wenjianjiaA_96a9ab047699462faa27",
                          "fields": [
                            {
                              "name": "aaa",
                              "type": "INTEGER"
                            },
                            {
                              "name": "status",
                              "type": "LONG"
                            }
                          ],
                          "dataType": "TIME_SEQUENCE_TYPE",
                          "generateDashboard": "FALSE",
                          "enableHistory": "FALSE",
                          "mockData": "FALSE",
                          "type": "file",
                          "name": "时序1",
                          "topicType": "METRIC"
                        },
                        {
                          "alias": "wenjianjiaA_c515c3d4e61a4a8f9024",
                          "fields": [
                            {
                              "name": "a",
                              "type": "LONG"
                            }
                          ],
                          "dataType": "TIME_SEQUENCE_TYPE",
                          "generateDashboard": "TRUE",
                          "enableHistory": "FALSE",
                          "mockData": "FALSE",
                          "type": "file",
                          "name": "asdfa",
                          "topicType": "METRIC"
                        },
                        {
                          "alias": "wenjianjiaA_ce18b74fd6c74ef2b5a4",
                          "fields": [
                            {
                              "name": "a",
                              "type": "LONG"
                            }
                          ],
                          "dataType": "TIME_SEQUENCE_TYPE",
                          "generateDashboard": "TRUE",
                          "enableHistory": "FALSE",
                          "mockData": "FALSE",
                          "type": "file",
                          "name": "asdfa",
                          "topicType": "METRIC"
                        }
                      ],
                      "topicType": "METRIC"
                    }
]
`
