package jsonstream

import (
	"backend/internal/common/I18nUtils"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/zeromicro/go-zero/core/logx"
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

type parseContext[TreeNode any, FlatNode any] struct {
	batchSize     int
	reader        *loggedReader
	decoder       *json.Decoder
	tree2flat     func(propName string, node, parent *TreeNode) *FlatNode
	errConsumer   func(node *TreeNode)
	consumer      func(readSize int64, propName string, nodes []*FlatNode)
	fieldIndexMap map[string]int
	childrenName  string

	prevPropName string
	cacheBatch   []*FlatNode
}

// DecodeJsonTreeToFlat 把树形json流式解析为平铺结构
func DecodeJsonTreeToFlat[TreeNode any, FlatNode any](
	reader io.Reader,
	batchSize int,
	tree2flat func(propName string, node, parent *TreeNode) *FlatNode,
	consumer func(readSize int64, propName string, nodes []*FlatNode),
	errConsumer func(node *TreeNode),
) error {
	lr := &loggedReader{tar: reader}
	decoder := json.NewDecoder(lr)

	// 确保第一个 token 是对象开始
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	batchSize = max(batchSize, 1)
	ctx := &parseContext[TreeNode, FlatNode]{
		batchSize:   batchSize,
		reader:      lr,
		decoder:     decoder,
		tree2flat:   tree2flat,
		consumer:    consumer,
		errConsumer: errConsumer,
		cacheBatch:  make([]*FlatNode, 0, batchSize),
	}
	defer func() {
		finish[TreeNode, FlatNode](ctx)
	}()
	{
		var defNode TreeNode
		ctx.fieldIndexMap, ctx.childrenName = parseJsonFields(defNode)
	}
	if delim, ok := token.(json.Delim); !ok {
		return fmt.Errorf("expected json start")
	} else if delim == '[' {
		return parseChildrenArray(ctx, nil, "", false)
	} else if delim == '{' {
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
			err = parseChildrenArray(ctx, nil, propName, false)
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
	} else {
		return fmt.Errorf("expected json object or array start")
	}
}
func accept[TreeNode any, FlatNode any](ctx *parseContext[TreeNode, FlatNode], node, parent *TreeNode, propName string) bool {
	flat := ctx.tree2flat(propName, node, parent)
	if flat == nil {
		ctx.prevPropName = propName
		ctx.errConsumer(node)
		return false
	}
	if ctx.prevPropName == "" {
		ctx.prevPropName = propName
	}
	sameProp := propName == ctx.prevPropName
	prevProp := ctx.prevPropName
	ctx.prevPropName = propName
	if sameProp {
		ctx.cacheBatch = append(ctx.cacheBatch, flat)
		if batch := ctx.cacheBatch; len(batch) >= ctx.batchSize {
			ctx.cacheBatch = make([]*FlatNode, 0, ctx.batchSize)
			ctx.consumer(ctx.reader.readSize, propName, batch)
		}
	} else if batch := ctx.cacheBatch; len(batch) > 0 {
		ctx.cacheBatch = make([]*FlatNode, 0, ctx.batchSize)
		ctx.consumer(ctx.reader.readSize, prevProp, batch)

		ctx.cacheBatch = append(ctx.cacheBatch, flat)
	} else {
		ctx.cacheBatch = append(ctx.cacheBatch, flat)
	}
	if propName == "UNS" {
		logx.Debugf("Add Uns[%d]: %+v", len(ctx.cacheBatch), *flat)
	}
	return true
}
func finish[TreeNode any, FlatNode any](ctx *parseContext[TreeNode, FlatNode]) {
	logx.Info("finish: ", len(ctx.cacheBatch))
	if len(ctx.cacheBatch) > 0 {
		ctx.consumer(ctx.reader.readSize, ctx.prevPropName, ctx.cacheBatch)
		ctx.cacheBatch = ctx.cacheBatch[:0]
	}
}
func jsonErr(err error) error {
	if je, is := err.(*json.SyntaxError); is {
		return fmt.Errorf("%s: %d: %v", I18nUtils.GetMessage("uns.import.json.error"), je.Offset, je.Error())
	}
	return err
}

// 递归解析单个节点
func parseNode[TreeNode any, FlatNode any](ctx *parseContext[TreeNode, FlatNode], parent *TreeNode, propName string) error {
	hasChild := false
	countFields := 0
	var node TreeNode
	var values = reflect.ValueOf(&node).Elem()
	for ctx.decoder.More() {
		// 读取字段名
		token, err := ctx.decoder.Token()
		if err != nil {
			return err
		}
		fieldName, ok := token.(string)
		if !ok {
			_ = skipValue(ctx.decoder)
			continue
		}
		if fieldName == ctx.childrenName {
			// 手动解析 children 数组
			hasChild = true
			if accept(ctx, &node, parent, propName) {
				err = parseChildrenArray[TreeNode, FlatNode](ctx, &node, propName, true)
				if err != nil {
					return err
				}
			} else {
				err = skipValue(ctx.decoder)
			}
		} else if index, has := ctx.fieldIndexMap[fieldName]; has {
			// 根据字段名处理值
			fieldValue := values.Field(index)
			err = ctx.decoder.Decode(fieldValue.Addr().Interface())
			if err == nil {
				countFields++
			}
		} else {
			logx.Error("跳过未知字段: ", fieldName)
			_ = skipValue(ctx.decoder)
		}
	}

	// 读取对象结束标记
	token, err := ctx.decoder.Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected object end")
	}
	if !hasChild && countFields > 0 {
		accept(ctx, &node, parent, propName)
	}
	return nil
}

// 解析 children 数组
func parseChildrenArray[TreeNode any, FlatNode any](ctx *parseContext[TreeNode, FlatNode], parent *TreeNode, propName string, checkStart bool) error {
	if checkStart {
		// 检查数组开始标记
		token, err := ctx.decoder.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); !ok || delim != '[' {
			return fmt.Errorf("expected array start for children")
		}
	}

	// 遍历数组中的元素
	for ctx.decoder.More() {
		// 检查子元素是否为对象
		token, err := ctx.decoder.Token()
		if err != nil {
			return err
		}

		if delim, ok := token.(json.Delim); ok && delim == '{' {
			// 递归解析子节点
			err := parseNode(ctx, parent, propName)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("expected object in children array")
		}
	}

	// 读取数组结束标记
	token, err := ctx.decoder.Token()
	if err != nil {
		return err
	}

	if delim, ok := token.(json.Delim); !ok || delim != ']' {
		return fmt.Errorf("expected array end for children")
	}

	return nil
}
func skipValue(decoder *json.Decoder) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch token {
	case json.Delim('['), json.Delim('{'):
		// 递归跳过数组或对象
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
		if token == json.Delim('[') && endToken != json.Delim(']') ||
			token == json.Delim('{') && endToken != json.Delim('}') {
			return fmt.Errorf("mismatched delimiters")
		}
	}

	return nil
}

func parseJsonFields(ts any) (fieldIndexMap map[string]int, childrenName string) {
	t := reflect.TypeOf(ts)
	// 如果传入的是指针，获取其指向的类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 确保是结构体类型
	if t.Kind() != reflect.Struct {
		return
	}
	fieldIndexMap = make(map[string]int, 16)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		name := parseColumnName(jsonTag)
		if len(name) > 0 {
			fieldIndexMap[name] = i

			if field.Type.Kind() == reflect.Slice &&
				(field.Type.Elem() == t || (field.Type.Elem().Kind() == reflect.Ptr && field.Type.Elem().Elem() == t)) {
				childrenName = name
			}
		}
	}
	return
}
