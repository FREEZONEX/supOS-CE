package noderedutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/internal/common/enums"

	"github.com/zeromicro/go-zero/core/logx"
)

// NodeRedUtils provides utility functions for Node-RED operations.
//
// GetTag retrieves tags from Node-RED by node ID.
func GetTag(nodeID, nodeRedHost, nodeRedPort string) (string, error) {
	var url string
	if nodeRedPort != "" {
		url = fmt.Sprintf("http://%s:%s/nodered-api/load/tags?nodeId=%s", nodeRedHost, nodeRedPort, nodeID)
	} else {
		url = fmt.Sprintf("%s/nodered-api/load/tags?nodeId=%s", nodeRedHost, nodeID)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// SaveTags saves tags to Node-RED.
func SaveTags(nodeID string, tags [][]string, nodeRedHost, nodeRedPort string) error {
	var url string
	url = fmt.Sprintf("http://%s:%s/nodered-api/save/tags", nodeRedHost, nodeRedPort)

	requestBody := map[string]any{
		"nodeId": nodeID,
		"tags":   tags,
	}

	reqJSON, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Create creates or updates a Node-RED flow.
func Create(globalExportModule enums.GlobalExportModuleEnum, flowID, flowName, description string, nodes any, nodeRedHost, nodeRedPort string) (string, error) {
	if nodes == nil {
		nodes = []any{}
	}

	requestBody := map[string]any{
		"id":       flowID,
		"nodes":    nodes,
		"disabled": false,
		"label":    flowName,
		"info":     description,
	}

	reqJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	logx.Debugf(globalExportModule.Code+" 创建 nodered 请求: %s", string(reqJSON))

	var url string
	var method string

	if flowID != "" {
		// Update existing flow
		method = http.MethodPut
		if nodeRedPort != "" {
			url = fmt.Sprintf("http://%s:%s/flow/%s", nodeRedHost, nodeRedPort, flowID)
		} else {
			url = fmt.Sprintf("%s/flow/%s", nodeRedHost, flowID)
		}
	} else {
		// Create new flow
		method = http.MethodPost
		if nodeRedPort != "" {
			url = fmt.Sprintf("http://%s:%s/flow", nodeRedHost, nodeRedPort)
		} else {
			url = fmt.Sprintf("%s/flow", nodeRedHost)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: 10 * time.Minute, // 10 minutes timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf(globalExportModule.Code+" 创建 nodered 返回结果: %s", string(body))

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("no id in response")
}

// Get retrieves a Node-RED flow by ID.
func Get(nodeRedHost, nodeRedPort, flowID string) (string, error) {
	url := fmt.Sprintf("http://%s:%s/flow/%s", nodeRedHost, nodeRedPort, flowID)
	logx.Debugf("eventFlow 查询 nodered 请求: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logx.Debugf("eventFlow 查询 nodered 返回结果: %s", string(body))

	if resp.StatusCode != 200 {
		return "", nil
	}
	return string(body), nil
}
