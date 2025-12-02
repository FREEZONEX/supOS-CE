package jsonstream

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
)

// StreamJsonEncoder 使用栈实现深度优先遍历的JSON编码器
type StreamJsonEncoder[Node any, ID cmp.Ordered, TreeNode any] struct {
	writer        *bufio.Writer
	childrenStart string

	stack           []stackItem[Node, TreeNode] // 栈，用于跟踪当前节点路径
	isFirst         bool                        // 是否是第一个节点
	hasStarted      bool                        // 是否已开始写入JSON数组
	nodeGetChildren func(*Node) []*Node
	nodeSetChildren func(*Node, []*Node)
	getId           func(*TreeNode) ID
	getParentId     func(*TreeNode) ID
	tree2node       func(*TreeNode) *Node
}
type stackItem[Node any, TreeNode any] struct {
	node *Node
	tree *TreeNode
}

func NewStreamJsonEncoder[Node any, ID cmp.Ordered, TreeNode any](
	jsonWriter io.Writer,
	nodeGetChildren func(*Node) []*Node,
	nodeSetChildren func(*Node, []*Node),
	getId func(*TreeNode) ID,
	getParentId func(*TreeNode) ID,
	tree2node func(*TreeNode) *Node,
) *StreamJsonEncoder[Node, ID, TreeNode] {
	var defNode Node
	childrenStart := ",\"children\":["
	childrenName := getChildrenJsonTagName(defNode)
	if len(childrenName) > 0 && childrenName != "children" {
		childrenStart = fmt.Sprintf(",\"%s\":[", childrenName)
	}
	writer := bufio.NewWriter(jsonWriter)
	return &StreamJsonEncoder[Node, ID, TreeNode]{
		writer:          writer,
		childrenStart:   childrenStart,
		stack:           make([]stackItem[Node, TreeNode], 0, 7),
		isFirst:         true,
		nodeGetChildren: nodeGetChildren,
		nodeSetChildren: nodeSetChildren,
		getId:           getId,
		getParentId:     getParentId,
		tree2node:       tree2node,
	}
}

func WriteBatch[Node any, ID cmp.Ordered, TreeNode any](
	je *StreamJsonEncoder[Node, ID, TreeNode],
	nodes []*TreeNode,
	isFinal bool,
) error {

	if !je.hasStarted {
		je.hasStarted = true
		je.writer.WriteString("[\n")
	}
	if len(nodes) > 0 {
		for _, treeNode := range nodes {
			// 创建新节点
			newNode := je.tree2node(treeNode)

			// 处理栈的状态
			if len(je.stack) == 0 {
				// 栈为空，这是一个根节点
				if !je.isFirst {
					je.writer.WriteString(",\n")
				}
				writeNodeStart(je, newNode)
				je.stack = append(je.stack, stackItem[Node, TreeNode]{node: newNode, tree: treeNode})

				//logx.Info("len(je.stack) = ", len(je.stack))
			} else {
				// 检查当前节点与栈顶节点的关系
				topNode := je.stack[len(je.stack)-1]

				if je.getParentId(treeNode) == je.getId(topNode.tree) {
					// 当前节点是栈顶节点的子节点
					// 如果这是栈顶节点的第一个子节点，需要开始children数组
					children := je.nodeGetChildren(topNode.node)
					if len(children) == 0 {
						je.writer.WriteString(je.childrenStart)
					} else {
						je.writer.WriteString(",")
					}

					writeNodeStart(je, newNode)

					je.nodeSetChildren(topNode.node, append(children, newNode))
					je.stack = append(je.stack, stackItem[Node, TreeNode]{node: newNode, tree: treeNode})

					//logx.Infof("len(stack) = %d, topChildren: %d\n", len(je.stack), len(children))
				} else {
					// 当前节点不是栈顶节点的子节点，需要回溯
					// 弹出栈直到找到父节点或栈为空
					for len(je.stack) > 0 {
						topNode = je.stack[len(je.stack)-1]

						// 闭合当前节点
						writeNodeEnd(je)
						je.stack = je.stack[:len(je.stack)-1]

						// 检查栈顶节点是否是当前节点的父节点
						if len(je.stack) > 0 && je.getParentId(treeNode) == je.getId(je.stack[len(je.stack)-1].tree) {
							break
						}
					}

					// 如果栈为空，当前节点是根节点
					if len(je.stack) == 0 {
						if !je.isFirst {
							je.writer.WriteString(",\n")
						}
						writeNodeStart(je, newNode)
						je.stack = append(je.stack, stackItem[Node, TreeNode]{node: newNode, tree: treeNode})
					} else {
						// 当前节点是栈顶节点的子节点
						topNode = je.stack[len(je.stack)-1]
						children := je.nodeGetChildren(topNode.node)
						if len(children) == 0 {
							je.writer.WriteString(je.childrenStart)
						} else {
							je.writer.WriteString(",")
						}

						writeNodeStart(je, newNode)
						je.nodeSetChildren(topNode.node, append(children, newNode))
						je.stack = append(je.stack, stackItem[Node, TreeNode]{node: newNode, tree: treeNode})
					}
				}
			}

			je.isFirst = false
		}
	}
	if isFinal {
		// 闭合所有栈中的节点
		for len(je.stack) > 0 {
			writeNodeEnd(je)
			je.stack = je.stack[:len(je.stack)-1]
		}

		je.writer.WriteString("\n]")
	}

	return je.writer.Flush()
}

// writeNodeStart 写入节点的开始部分
func writeNodeStart[Node any, ID cmp.Ordered, TreeNode any](
	je *StreamJsonEncoder[Node, ID, TreeNode], node *Node) {
	jsonBytes, _ := json.Marshal(node)
	je.writer.Write(jsonBytes[:len(jsonBytes)-1])
}

// writeNodeEnd 写入节点的结束部分
func writeNodeEnd[Node any, ID cmp.Ordered, TreeNode any](
	je *StreamJsonEncoder[Node, ID, TreeNode]) {
	topNode := je.stack[len(je.stack)-1]
	if len(je.nodeGetChildren(topNode.node)) > 0 {
		je.writer.WriteString("]")
	}
	je.writer.WriteString("}")
}
