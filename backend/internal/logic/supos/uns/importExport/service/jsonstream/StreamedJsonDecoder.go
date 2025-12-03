package jsonstream

import (
	"backend/internal/common/I18nUtils"
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
			return jsonErr(err)
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
			return jsonErr(er)
		}

		propName, isString := fieldName.(string)
		if !isString {
			// 跳过未知字段的值
			continue
		}

		// 读取数组开始标记
		t, err := decoder.Token()
		if err != nil {
			return jsonErr(err)
		}
		if t != json.Delim('[') {
			continue // 跳过未知字段的值
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
		return jsonErr(err)
	}
	if t != json.Delim('}') {
		return fmt.Errorf("expected object end, got %v", t)
	}

	return nil
}
func jsonErr(err error) error {
	if je, is := err.(*json.SyntaxError); is {
		return fmt.Errorf("%s: %d: %v", I18nUtils.GetMessage("uns.import.json.error"), je.Offset, je.Error())
	}
	return err
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
