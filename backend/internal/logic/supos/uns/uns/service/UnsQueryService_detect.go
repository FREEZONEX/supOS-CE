package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/dto"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

var httClient = &http.Client{
	Timeout: time.Second * 3,
}

func (l *UnsQueryService) DetectIfFieldReferenced(ctx context.Context, req *types.UpdateModeRequestVo) (resp *types.ResultVO, err error) {
	dataMap := map[string]any{"referred": false}
	resp = &types.ResultVO{Code: 200, Msg: "ok", Data: dataMap}

	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.GetByAlias(db, req.Alias)
	if er != nil || uns == nil {
		err = er
		return
	}
	if uns.PathType == 0 {
		return
	}
	files, _ := l.unsMapper.ListFileByTemplateId(db, uns.Id)
	if len(files) == 0 {
		return
	}
	instanceTopics := base.Map[*dao.UnsNamespace, string](files, func(e *dao.UnsNamespace) string {
		return e.Path
	})
	nodeRedRef, er := detectReferencedNodeRed(instanceTopics)
	if er != nil {
		err = er
	} else if nodeRedRef {
		dataMap["referred"] = true
		dataMap["tips"] = I18nUtils.GetMessage("uns.update.field.tips1")
	}
	return
}

// 检查 NodeRed 引用
func detectReferencedNodeRed(instanceTopics []string) (bool, error) {
	// 构建请求 URL
	url := "http://localhost:8080/service-api/supos/flow/by/topics"

	// 准备请求体
	requestBody := struct {
		Topics []string `json:"topics"`
	}{
		Topics: instanceTopics,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return false, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return false, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := httClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("执行请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("读取响应体失败: %w", err)
	}

	// 解析 JSON 响应
	var result dto.ResultDTO[[]*dto.NodeFlowDTO]
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("解析响应 JSON 失败: %w", err)
	}

	// 检查数据是否非空
	hasReferences := result.Data != nil && len(result.Data) > 0
	return hasReferences, nil
}
