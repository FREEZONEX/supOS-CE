package jsonstream

import (
	"bufio"
	"cmp"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
)

// Csv2JsonStream 从 CopyTo 流式生成 JSON
func Csv2JsonStream[Node any, ID cmp.Ordered](
	csvExporter func(io.Writer) error,
	jsonWriter io.Writer,
	nodeGetChildren func(*Node) []*Node,
	nodeSetChildren func(*Node, []*Node),
	getId func(*Node) ID,
	getParentId func(*Node) ID,
	csv2node func(headers, values []string) *Node,
	fullList bool,
) (countNodes int, er error) {
	// 创建管道
	csvReader, csvWriter := io.Pipe()

	// 错误通道
	errChan := make(chan error, 2)

	// Goroutine 1: 执行 CopyTo 写入 CSV 到管道
	go func() {
		defer csvWriter.Close()
		err := csvExporter(csvWriter)
		if err != nil {
			errChan <- fmt.Errorf("CSV 导出失败: %v", err)
			return
		}
		errChan <- nil
	}()
	var defNode Node
	childrenName := getChildrenJsonTagName(defNode)
	if len(childrenName) == 0 {
		childrenName = "children"
	}
	// Goroutine 2: 从管道读取 CSV 并转换为 JSON
	go func() {
		defer csvReader.Close()
		var csvErr error
		if countNodes, csvErr = csvStreamWriteJson(csvReader, jsonWriter, nodeGetChildren, nodeSetChildren, getId, getParentId, csv2node, childrenName, fullList); csvErr != nil {
			errChan <- fmt.Errorf("CSV转JSON失败: %v", csvErr)
			return
		}
		errChan <- nil
	}()

	// 等待两个 goroutine 完成
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return countNodes, err
		}
	}

	return countNodes, nil
}

func csvStreamWriteJson[Node any, ID cmp.Ordered](
	csvData *io.PipeReader,
	jsonWriter io.Writer,

	nodeGetChildren func(*Node) []*Node,
	nodeSetChildren func(*Node, []*Node),
	getId func(*Node) ID,
	getParentId func(*Node) ID,
	csv2node func(headers, values []string) *Node,
	childrenName string,
	fullList bool,
) (int, error) {

	// 读取 CSV 表头
	csvReader := csv.NewReader(csvData)
	headers, err := csvReader.Read()
	if err != nil {
		return -1, fmt.Errorf("读取CSV表头失败: %v", err)
	}
	writer := bufio.NewWriter(jsonWriter)

	if fullList {
		// 开始写入 JSON 数组的开始标记
		if _, err = writer.Write([]byte("[\n")); err != nil {
			return -1, err
		}
	}
	var stack []*Node

	writeNodeStart := func(node *Node) {
		jsonBytes, _ := json.Marshal(node)
		writer.Write(jsonBytes[:len(jsonBytes)-1])
	}
	writeNodeEnd := func() {
		topNode := stack[len(stack)-1]
		if len(nodeGetChildren(topNode)) > 0 {
			writer.WriteString("]")
		}
		writer.WriteString("}")
	}

	childrenStart := fmt.Sprintf(",\"%s\":[", childrenName)
	countNodes := 0
	for isFirst := true; ; {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return countNodes, fmt.Errorf("读取CSV行失败: %v", err)
		}
		countNodes++
		// 转换为 JSON 对象
		newNode := csv2node(headers, record)
		// 处理栈的状态
		if len(stack) == 0 {
			// 栈为空，这是一个根节点
			if !isFirst {
				writer.WriteString(",\n")
			} else {
				isFirst = false
			}
			writeNodeStart(newNode)
			stack = append(stack, newNode)

			//logx.Info("len(stack) = ", len(stack))
		} else {
			// 检查当前节点与栈顶节点的关系
			topNode := stack[len(stack)-1]

			if getParentId(newNode) == getId(topNode) {
				// 当前节点是栈顶节点的子节点
				// 如果这是栈顶节点的第一个子节点，需要开始children数组
				children := nodeGetChildren(topNode)
				if len(children) == 0 {
					writer.WriteString(childrenStart)
				} else {
					writer.WriteString(",")
				}

				writeNodeStart(newNode)

				nodeSetChildren(topNode, append(children, newNode))
				stack = append(stack, newNode)

				//logx.Infof("len(stack) = %d, topChildren: %d\n", len(stack), len(children))
			} else {
				// 当前节点不是栈顶节点的子节点，需要回溯
				// 弹出栈直到找到父节点或栈为空
				for len(stack) > 0 {
					topNode = stack[len(stack)-1]

					// 闭合当前节点
					writeNodeEnd()
					stack = stack[:len(stack)-1]

					// 检查栈顶节点是否是当前节点的父节点
					if len(stack) > 0 && getParentId(newNode) == getId(stack[len(stack)-1]) {
						break
					}
				}

				// 如果栈为空，当前节点是根节点
				if len(stack) == 0 {
					if !isFirst {
						writer.WriteString(",\n")
					}
					writeNodeStart(newNode)
					stack = append(stack, newNode)
				} else {
					// 当前节点是栈顶节点的子节点
					topNode = stack[len(stack)-1]
					children := nodeGetChildren(topNode)
					if len(children) == 0 {
						writer.WriteString(childrenStart)
					} else {
						writer.WriteString(",")
					}

					writeNodeStart(newNode)
					nodeSetChildren(topNode, append(children, newNode))
					stack = append(stack, newNode)
				}
			}
		}
	}

	// 闭合所有栈中的节点
	for len(stack) > 0 {
		writeNodeEnd()
		stack = stack[:len(stack)-1]
	}
	if fullList {
		writer.WriteString("\n]")
	}
	return countNodes, writer.Flush()
}
