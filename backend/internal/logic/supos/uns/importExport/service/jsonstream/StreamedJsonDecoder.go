package jsonstream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

type loggedReader struct {
	tar      io.Reader
	readSize int64
}

func (l *loggedReader) Read(p []byte) (n int, err error) {
	n, err = l.tar.Read(p)
	l.readSize += int64(n)
	return n, err
}
func DecodeStreamedJson[Node any, TreeNode any](
	jsonFile io.Reader,
	bufSize, batchSize int,
	getChildren func(*Node) []*Node,
	node2tree func(node, parent *Node) *TreeNode,
	batchConsumer func(readSize int64, prop string, list []*TreeNode)) error {

	reader := &loggedReader{tar: bufio.NewReaderSize(jsonFile, bufSize)}
	decoder := json.NewDecoder(reader)
	// 读取整个对象的开始标记
	{
		startChar, err := decoder.Token()
		if err != nil {
			return err
		}
		if startChar == json.Delim('[') {
			return decodeArray(decoder, reader, batchSize, getChildren, node2tree, "", batchConsumer)
		}
		if startChar != json.Delim('{') {
			return fmt.Errorf("expected object start, got %v", startChar)
		}
	}
	// 处理对象中的字段
	for decoder.More() {
		// 读取字段名
		fieldName, er := decoder.Token()
		if er != nil {
			return er
		}

		propName, isString := fieldName.(string)
		if !isString {
			// 跳过未知字段的值
			if err := skipValue(decoder); err != nil {
				return err
			}
		}

		// 读取数组开始标记
		t, err := decoder.Token()
		if err != nil {
			return err
		}
		if t != json.Delim('[') {
			// 跳过未知字段的值
			if err := skipValue(decoder); err != nil {
				return err
			}
		}

		// 解析数组
		err = decodeArray(decoder, reader, batchSize, getChildren, node2tree, propName, batchConsumer)
		if err != nil {
			return err
		}
	}

	// 读取对象结束标记
	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('}') {
		return fmt.Errorf("expected object end, got %v", t)
	}

	return nil
}

func decodeArray[Node any, TreeNode any](
	decoder *json.Decoder,
	loggedReader *loggedReader,
	batchSize int,
	getChildren func(*Node) []*Node,
	node2tree func(node, parent *Node) *TreeNode,
	propName string,
	batchConsumer func(readSize int64, prop string, list []*TreeNode)) error {

	var batch []*TreeNode
	type StackItem struct {
		node   *Node
		parent *Node
	}
	// 使用栈来避免递归
	var stack []StackItem

	// 解析根节点
	for decoder.More() {
		var node Node
		if err := decoder.Decode(&node); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error decoding JSONArray: %v", err)
		}

		// 将根节点入栈
		stack = append(stack, StackItem{node: &node})

		// 处理栈中的节点
		for len(stack) > 0 {
			// 弹出栈顶元素
			item := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if item.node == nil {
				continue
			}

			// 处理当前节点
			treeNode := node2tree(item.node, item.parent)
			batch = append(batch, treeNode)

			// 检查批次大小
			if len(batch) >= batchSize {
				batchConsumer(loggedReader.readSize, propName, batch)
				batch = make([]*TreeNode, 0, batchSize)
			}

			// 将子节点逆序入栈（保证处理顺序）
			children := getChildren(item.node)
			for i := len(children) - 1; i >= 0; i-- {
				stack = append(stack, StackItem{
					node:   children[i],
					parent: item.node,
				})
			}
		}
	}

	if len(batch) > 0 { // 处理剩余的批次
		batchConsumer(loggedReader.readSize, propName, batch)
	} else if batch == nil && propName != "" {
		batchConsumer(loggedReader.readSize, propName, []*TreeNode{})
	}

	// 读取数组结束标记
	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim(']') {
		return fmt.Errorf("expected array end, got %v", t)
	}
	return nil
}

// 跳过未知字段的值
func skipValue(decoder *json.Decoder) error {
	t, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t {
	case json.Delim('['), json.Delim('{'):
		// 对于数组或对象，需要递归跳过所有内容
		for {
			if !decoder.More() {
				break
			}
			if err := skipValue(decoder); err != nil {
				return err
			}
		}
		// 读取结束标记
		endToken, err := decoder.Token()
		if err != nil {
			return err
		}
		if (t == json.Delim('[') && endToken != json.Delim(']')) ||
			(t == json.Delim('{') && endToken != json.Delim('}')) {
			return fmt.Errorf("unexpected end token %v", endToken)
		}
	}
	// 对于其他类型（字符串、数字、布尔、null），已经读取了值，无需额外处理

	return nil
}
