package adapter

import (
	"backend/internal/common/I18nUtils"
	"backend/share/app/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/errors"
	"gopkg.in/yaml.v3"
)

// PortainerConfig Portainer 配置
type PortainerConfig struct {
	// Portainer API 地址
	BaseURL string
	// Portainer 用户名（用于获取JWT token）
	Username string
	// Portainer 密码
	Password string
	// Portainer API 密钥（直接使用API密钥，优先级最高）
	APIKey string
	// Docker 环境ID
	EndpointID int
	// 超时时间（秒）
	Timeout int
	// StackId(tier0)
	StackId int
}

// DefaultPortainerConfig 默认 Portainer 配置
func DefaultPortainerConfig() *PortainerConfig {
	return &PortainerConfig{
		BaseURL:    getEnvOrDefault("PORTAINER_URL", "http://portainer:9000/api"),
		Username:   getEnvOrDefault("PORTAINER_USERNAME", "admin"),
		Password:   getEnvOrDefault("PORTAINER_PASSWORD", "adminpassword"),
		APIKey:     getEnvOrDefault("PORTAINER_API_KEY", ""),
		EndpointID: getEnvIntOrDefault("PORTAINER_ENDPOINT_ID", 0),
		Timeout:    getEnvIntOrDefault("PORTAINER_TIMEOUT", 300000),
	}
}

// InitializeEndpointId 初始化EndpointId，通过查询Portainer端点列表获取本地环境ID
func (p *PortainerAdapter) InitializeEndpointId() error {
	log.Printf("[app]: 开始初始化 Portainer EndpointId...")

	// 1. 首先确保认证
	if err := p.authenticate(); err != nil {
		log.Printf("[app]: 认证失败，无法初始化EndpointId: %v", err)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.failed"))
	}

	// 2. 构建请求URL
	// Portainer 端点列表API: GET /api/endpoints?search=local
	url := fmt.Sprintf("%s/endpoints?search=%s", p.config.BaseURL, "local")

	log.Printf("[app]: 查询Portainer端点列表: %s", url)

	// 3. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.endpoint.list.request.create.failed"))
	}

	// 4. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 5. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.endpoint.list.request.send.failed"))
	}
	defer resp.Body.Close()

	// 6. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[app]: Portainer API 返回错误: %s, 响应: %s", resp.Status, string(body))
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.endpoint.search.error"))
	}

	// 7. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	var endpoints []map[string]interface{}
	if err := json.Unmarshal(body, &endpoints); err != nil {
		log.Printf("[app]: 解析端点列表失败，响应内容: %s", string(body))
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.endpoint.list.parse.failed"))
	}

	log.Printf("[app]: 找到 %d 个Portainer端点", len(endpoints))

	// 8. 查找本地环境
	var localEndpointId int
	found := false

	// 首先尝试查找名称包含"local"的环境
	for _, endpoint := range endpoints {
		// 检查端点名称或类型
		name, _ := endpoint["Name"].(string)
		endpointType, _ := endpoint["Type"].(float64)
		id, _ := endpoint["Id"].(float64)

		// Type: 1=Docker, 2=Agent, 3=Azure, 4=EdgeAgent, 5=EdgeStandard
		// 我们查找Docker类型的本地环境
		if endpointType == 1 { // Docker环境
			// 检查名称是否包含"local"或"Local"（不区分大小写）
			if strings.Contains(strings.ToLower(name), "local") {
				localEndpointId = int(id)
				found = true
				log.Printf("[app]: 找到本地Docker环境: %s (ID: %d)", name, localEndpointId)
				break
			}
		}
	}

	// 9. 如果没有找到明确的本地环境，使用第一个Docker环境
	if !found && len(endpoints) > 0 {
		for _, endpoint := range endpoints {
			endpointType, _ := endpoint["Type"].(float64)
			id, _ := endpoint["Id"].(float64)
			name, _ := endpoint["Name"].(string)

			if endpointType == 1 { // Docker环境
				localEndpointId = int(id)
				found = true
				log.Printf("[app]: 使用第一个Docker环境: %s (ID: %d)", name, localEndpointId)
				break
			}
		}
	}

	// 10. 更新配置
	if found {
		oldEndpointId := p.config.EndpointID
		p.config.EndpointID = localEndpointId
		log.Printf("[app]: Portainer EndpointId 初始化成功: %d -> %d",
			oldEndpointId, p.config.EndpointID)
		return nil
	}

	log.Printf("[app]: 未找到可用的Docker环境")
	return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.endpoint.not.found"))
}

func (p *PortainerAdapter) InitializeStackId(endpointId int, composeYaml string) error {
	log.Printf("[app]: 开始初始化 Portainer StackId...")

	// 2. 构建请求URL
	// Portainer 端点列表API: GET /api/stacks?search=local
	url := fmt.Sprintf("%s/stacks?filters={\"EndpointId\":%d}", p.config.BaseURL, endpointId)

	log.Printf("[app]: 查询Stack列表: %s", url)

	// 3. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.list.request.create.failed"))
	}

	// 4. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 5. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.list.request.send.failed"))
	}
	defer resp.Body.Close()

	// 6. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[app]: Portainer API 返回错误: %s, 响应: %s", resp.Status, string(body))
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.search.error"))
	}

	// 7. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	var stacks []map[string]interface{}
	if err := json.Unmarshal(body, &stacks); err != nil {
		log.Printf("[app]: 解析stack列表失败，响应内容: %s", string(body))
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.list.parse.failed"))
	}

	log.Printf("[app]: 找到 %d 个Stack", len(stacks))

	// 8. 查找本地环境
	if len(stacks) > 0 {
		for _, stack := range stacks {
			// 检查Stack名称
			id := stack["Id"].(float64)
			p.config.StackId = int(id)
			break
		}
		return nil
	} else if composeYaml != "" {
		// 创建一个空的stack
		log.Printf("[app]: 没有找到现有Stack，开始创建新的空Stack...")
		return p.createEmptyStack(endpointId, composeYaml)
	}

	log.Printf("[app]: 未找到可用的Docker环境")
	return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.not.found"))
}

// createEmptyStack 创建一个空的Portainer Stack
func (p *PortainerAdapter) createEmptyStack(endpointId int, composeYml string) error {
	log.Printf("[app]: 开始创建空的Portainer Stack，名称: tier0")

	// 1. 构建请求体

	stackReq := map[string]interface{}{
		"Name":             "tier0",
		"StackFileContent": composeYml,
		"FromAppTemplate":  false,
	}

	reqJSON, err := json.Marshal(stackReq)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.create.request.serialize.failed"))
	}

	// 2. 构建请求URL
	// Portainer 创建Stack API: POST /api/stacks?endpointId={id}&method=string
	url := fmt.Sprintf("%s/stacks?endpointId=%d&method=string&type=2", p.config.BaseURL, endpointId)
	log.Printf("[app]: 创建Stack URL: %s", url)

	// 3. 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqJSON))
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.create.request.create.failed"))
	}

	// 4. 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 5. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 6. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.create.request.send.failed"))
	}
	defer resp.Body.Close()

	// 7. 检查响应状态
	log.Printf("[app]: Stack创建响应状态: %s", resp.Status)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[app]: 创建stack失败 返回错误: %s, 响应: %s", resp.Status, string(body))
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.create.failed"))
	}

	// 8. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	log.Printf("[app]: Stack创建响应: %s", string(body))

	var stackResp map[string]interface{}
	if err := json.Unmarshal(body, &stackResp); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.create.response.parse.failed"))
	}

	// 9. 提取 Stack ID
	stackID, ok := stackResp["Id"].(float64)
	if !ok {
		// 尝试其他可能的字段名
		if id, ok := stackResp["ID"].(float64); ok {
			stackID = id
		} else if id, ok := stackResp["id"].(float64); ok {
			stackID = id
		} else {
			log.Printf("[app]: 无法从响应中提取Stack ID，响应: %v", stackResp)
			return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.create.id.extract.failed"))
		}
	}

	// 10. 更新配置
	oldStackId := p.config.StackId
	p.config.StackId = int(stackID)
	log.Printf("[app]: Portainer StackId 创建成功: %d -> %d (名称: stack-app)",
		oldStackId, p.config.StackId)

	log.Printf("[app]: Stack创建完成，ID: %d", p.config.StackId)
	return nil
}

// verifyStackExists 验证Stack是否存在
func (p *PortainerAdapter) verifyStackExists(stackID int) error {
	log.Printf("[app]: 验证Stack是否存在，ID: %d", stackID)

	// 1. 构建请求URL
	url := fmt.Sprintf("%s/stacks/%d", p.config.BaseURL, stackID)

	// 2. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.verify.request.create.failed"))
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.verify.request.send.failed"))
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.verify.failed.status"))
	}

	// 6. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.verify.response.read.failed"))
	}

	var stackInfo map[string]interface{}
	if err := json.Unmarshal(body, &stackInfo); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.verify.response.parse.failed"))
	}

	// 7. 检查Stack信息
	if name, ok := stackInfo["Name"].(string); ok && name == "stack-app" {
		log.Printf("[app]: Stack验证成功: %s (ID: %d)", name, stackID)
		return nil
	}

	return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.verify.failed"))
}

// PortainerAdapter Portainer 适配器
type PortainerAdapter struct {
	config *PortainerConfig
	client *http.Client
	// JWT token 缓存
	token       string
	tokenExpiry time.Time
}

// NewPortainerAdapter 创建 Portainer 适配器
func NewPortainerAdapter(config *PortainerConfig, composeYaml string) *PortainerAdapter {
	if config == nil {
		config = DefaultPortainerConfig()
	}

	adapter := &PortainerAdapter{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
	adapter.autoInitializeEndpointId(composeYaml)
	return adapter
}

// autoInitializeEndpointId 自动初始化EndpointId（内部方法）
func (p *PortainerAdapter) autoInitializeEndpointId(composeYaml string) error {
	// 只有在使用默认EndpointId且配置了认证信息时才尝试自动初始化
	if p.config.EndpointID == 0 && (p.config.APIKey != "" || (p.config.Username != "" && p.config.Password != "")) {
		log.Printf("[app]: 尝试自动初始化 Portainer EndpointId...")

		// 尝试初始化EndpointId
		if err := p.InitializeEndpointId(); err == nil {
			log.Printf("[app]: EndpointId 自动初始化成功: %d", p.config.EndpointID)
		} else {
			p.config.EndpointID = 1
			log.Printf("[app]: EndpointId 自动初始化失败，使用默认值 1: %v", err)
		}

		if err := p.InitializeStackId(p.config.EndpointID, composeYaml); err == nil {
			log.Printf("[app]: StackId 自动初始化成功: %d", p.config.StackId)
		} else {
			log.Printf("[app]: StackId 自动初始化失败: %v", err)
			return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.initialize.failed"))
		}
	}
	return nil
}

// joinTier0Network 将容器加入到 tier0_edge_network 网络
func (p *PortainerAdapter) joinTier0Network(containerName string) error {
	if containerName == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.name.empty"))
	}

	log.Printf("[app]: 开始将容器 %s 加入到 tier0_edge_network 网络", containerName)

	// 1. 首先检查容器是否存在
	container, err := p.getContainerInfo(containerName)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.info.fetch.failed"))
	}

	// 2. 检查容器是否已经在目标网络中
	if p.isContainerInNetwork(container.ID, "tier0_edge_network") {
		log.Printf("[app]: 容器 %s 已经在 tier0_edge_network 网络中", containerName)
		return nil
	}

	// 3. 检查目标网络是否存在
	networkID, err := p.getNetworkID("tier0_edge_network")
	if err != nil {
		log.Printf("[app]: 网络 tier0_edge_network 不存在")
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.not.found"))
	}

	// 4. 将容器连接到网络
	if err := p.connectContainerToNetwork(container.ID, networkID, containerName); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.connect.failed"))
	}

	log.Printf("[app]: 容器 %s 已成功加入到 tier0_edge_network 网络", containerName)
	return nil
}

// getContainerInfo 获取容器信息
func (p *PortainerAdapter) getContainerInfo(containerName string) (*ContainerInfo, error) {
	// 1. 构建请求URL
	url := fmt.Sprintf("%s/endpoints/%d/docker/containers/%s/json", p.config.BaseURL, p.config.EndpointID, containerName)

	// 2. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.create.failed"))
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.send.failed"))
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	// 6. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	var containerInfo ContainerInfo
	if err := json.Unmarshal(body, &containerInfo); err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.parse.failed"))
	}

	return &containerInfo, nil
}

// isContainerInNetwork 检查容器是否已经在指定网络中
func (p *PortainerAdapter) isContainerInNetwork(containerID, networkName string) bool {
	// 1. 构建请求URL
	url := fmt.Sprintf("%s/endpoints/%d/docker/containers/%s/json", p.config.BaseURL, p.config.EndpointID, containerID)

	// 2. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[app]: 创建容器网络检查请求失败: %v", err)
		return false
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		log.Printf("[app]: add authorize header failed: %v", err)
		return false
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		log.Printf("[app]: 发送容器网络检查请求失败: %v", err)
		return false
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		log.Printf("[app]: 容器网络检查 API 返回错误: %s", resp.Status)
		return false
	}

	// 6. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[app]: 读取容器网络检查响应失败: %v", err)
		return false
	}

	var containerInfo map[string]interface{}
	if err := json.Unmarshal(body, &containerInfo); err != nil {
		log.Printf("[app]: 解析容器网络检查响应失败: %v", err)
		return false
	}

	// 7. 检查网络配置
	if networkSettings, ok := containerInfo["NetworkSettings"].(map[string]interface{}); ok {
		if networks, ok := networkSettings["Networks"].(map[string]interface{}); ok {
			_, exists := networks[networkName]
			return exists
		}
	}

	return false
}

// getNetworkID 获取网络ID
func (p *PortainerAdapter) getNetworkID(networkName string) (string, error) {
	// 1. 构建请求URL
	url := fmt.Sprintf("%s/endpoints/%d/docker/networks", p.config.BaseURL, p.config.EndpointID)

	// 2. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.list.request.create.failed"))
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.list.request.send.failed"))
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	// 6. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	var networks []map[string]interface{}
	if err := json.Unmarshal(body, &networks); err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.list.parse.failed"))
	}

	// 7. 查找目标网络
	for _, network := range networks {
		if name, ok := network["Name"].(string); ok && name == networkName {
			if id, ok := network["Id"].(string); ok {
				return id, nil
			}
		}
	}

	return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.not.found"))
}

// connectContainerToNetwork 将容器连接到网络
func (p *PortainerAdapter) connectContainerToNetwork(containerID, networkID, containerName string) error {
	// 1. 构建连接请求体
	connectReq := map[string]interface{}{
		"Container": containerID,
		"EndpointConfig": map[string]interface{}{
			"Aliases":   []string{containerName},
			"IPAddress": "", // 自动分配IP
		},
	}

	reqJSON, err := json.Marshal(connectReq)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.connect.request.serialize.failed"))
	}

	// 2. 构建请求URL
	url := fmt.Sprintf("%s/endpoints/%d/docker/networks/%s/connect", p.config.BaseURL, p.config.EndpointID, networkID)

	// 3. 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(reqJSON))
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.connect.request.create.failed"))
	}

	// 4. 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 5. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 6. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.connect.request.send.failed"))
	}
	defer resp.Body.Close()

	// 7. 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	log.Printf("[app]: 容器 %s 已连接到网络 tier0_edge_network", containerName)
	return nil
}

// ContainerInfo 容器信息结构体
type ContainerInfo struct {
	ID    string `json:"Id"`
	Name  string `json:"Name"`
	State struct {
		Status string `json:"Status"`
	} `json:"State"`
}

// JoinTier0Network 包装函数
func JoinTier0Network(containerName string) error {
	config := DefaultPortainerConfig()
	adapter := NewPortainerAdapter(config, "")
	return adapter.joinTier0Network(containerName)
}

// authenticate 获取认证token
func (p *PortainerAdapter) authenticate() error {
	// 如果已经有有效的token，直接使用
	if p.token != "" && time.Now().Before(p.tokenExpiry) {
		return nil
	}

	// 如果配置了API Key，直接使用
	if p.config.APIKey != "" {
		p.token = p.config.APIKey
		p.tokenExpiry = time.Now().Add(24 * time.Hour) // API Key 通常长期有效
		log.Printf("[app]: 使用配置的 Portainer API Key")
		return nil
	}

	// 否则使用用户名密码获取JWT token
	if p.config.Username == "" || p.config.Password == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.info.not.configured"))
	}

	log.Printf("[app]: 正在获取 Portainer JWT token...")

	// 构建认证请求
	authReq := map[string]string{
		"Username": p.config.Username,
		"Password": p.config.Password,
	}

	authJSON, err := json.Marshal(authReq)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.request.serialize.failed"))
	}

	// 发送认证请求
	url := p.config.BaseURL + "/auth"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(authJSON)))
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.request.create.failed"))
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.request.send.failed"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.ReadAll(resp.Body)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.failed"))
	}

	// 解析响应
	var authResp map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.response.read.failed"))
	}

	if err := json.Unmarshal(body, &authResp); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.response.parse.failed"))
	}

	// 提取JWT token
	jwt, ok := authResp["jwt"].(string)
	if !ok || jwt == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.auth.token.not.found"))
	}

	p.token = jwt
	// JWT token 通常有效期为8小时，这里设置为8小时以确保安全
	p.tokenExpiry = time.Now().Add(8 * time.Hour)

	log.Printf("[app]: Portainer JWT token 获取成功，有效期至: %s", p.tokenExpiry.Format("2006-01-02 15:04:05"))
	return nil
}

// addAuthHeader 添加认证头到请求
func (p *PortainerAdapter) addAuthHeader(req *http.Request) error {
	// 确保有有效的认证token
	if err := p.authenticate(); err != nil {
		return err
	}

	// 添加认证头
	if p.config.APIKey != "" {
		// 使用 API Key
		req.Header.Set("X-API-Key", p.token)
	} else {
		// 使用 JWT token
		req.Header.Set("Authorization", "Bearer "+p.token)
		log.Printf("[app]: add header Authorization=%s", "Bearer "+p.token)
	}

	return nil
}

// PullImageFromRemote 使用 Portainer API 拉取远程镜像
func (p *PortainerAdapter) PullImageFromRemote(imageUrl string) error {

	// 验证镜像URL格式
	if !strings.Contains(imageUrl, ":") {
		// 如果没有指定标签，默认使用latest
		imageUrl = imageUrl + ":latest"
	}

	log.Printf("[app]: 开始通过 Portainer 拉取远程镜像: %s", imageUrl)

	// 1. 构建请求URL
	url := fmt.Sprintf("%s/endpoints/%d/docker/images/create", p.config.BaseURL, p.config.EndpointID)

	// 2. 构建查询参数
	queryParams := fmt.Sprintf("?fromImage=%s", imageUrl)
	fullURL := url + queryParams

	// 3. 创建请求
	req, err := http.NewRequest("POST", fullURL, nil)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.pull.request.create.failed"))
	}

	// 4. 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 5. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 6. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.pull.request.send.failed"))
	}
	defer resp.Body.Close()

	// 7. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[app]: Portainer API 返回错误: %s, 响应: %s, url: %s", resp.Status, string(body), fullURL)
		return nil
	}

	// 8. 读取拉取进度
	buf := make([]byte, 1024)
	var lastLogTime time.Time
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("[app]: 读取拉取进度失败: %v", err)
			break
		}

		if n > 0 {
			// 解析JSON进度信息
			progress := string(buf[:n])
			lines := strings.Split(progress, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				// 尝试解析JSON进度
				var progressObj map[string]interface{}
				if err := json.Unmarshal([]byte(line), &progressObj); err == nil {
					// 提取进度信息
					if status, ok := progressObj["status"].(string); ok {
						// 每2秒打印一次进度，避免日志过多
						if time.Since(lastLogTime) > 2*time.Second {
							if id, ok := progressObj["id"].(string); ok && id != "" {
								log.Printf("[app]: Portainer 拉取进度: %s [%s]", status, id)
							} else {
								log.Printf("[app]: Portainer 拉取进度: %s", status)
							}
							lastLogTime = time.Now()
						}
					}
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	// 9. 验证镜像是否拉取成功
	if err := p.verifyImageExists(imageUrl); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.pull.verify.failed"))
	}

	log.Printf("[app]: Portainer 远程镜像拉取完成: %s", imageUrl)
	return nil
}

// verifyImageExists 验证镜像是否存在
func (p *PortainerAdapter) verifyImageExists(imageUrl string) error {
	// 1. 构建请求URL
	url := fmt.Sprintf("%s/endpoints/%d/docker/images/json", p.config.BaseURL, p.config.EndpointID)

	// 2. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.pull.verify.request.create.failed"))
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.pull.verify.request.send.failed"))
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	// 6. 解析响应
	var images []map[string]interface{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	if err := json.Unmarshal(body, &images); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.pull.verify.parse.failed"))
	}

	// 7. 查找镜像
	found := false
	for _, img := range images {
		if repoTags, ok := img["RepoTags"].([]interface{}); ok {
			for _, tag := range repoTags {
				if tagStr, ok := tag.(string); ok && tagStr == imageUrl {
					found = true
					log.Printf("[app]: 镜像验证成功: %s", imageUrl)
					break
				}
			}
		}
		if found {
			break
		}
	}

	if !found {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.not.found"))
	}

	return nil
}

// LoadImageFromLocal 使用 Portainer API 加载本地镜像文件
func (p *PortainerAdapter) LoadImageFromLocal(imageId string) error {
	if imageId == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.load.path.empty"))
	}
	filePath := util.FindFileByID(util.ATTACHMENT_DIR, imageId)
	if filePath == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.load.file.not.found") + ": " + imageId)
	}
	// 验证文件格式
	ext := filepath.Ext(filePath)
	if ext != ".tar" && ext != ".tar.gz" && ext != ".tgz" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.load.invalid.format"))
	}

	log.Printf("[app]: 开始通过 Portainer 加载本地镜像文件: %s", filePath)
	// 1. 读取镜像文件
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.load.file.read.failed"))
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.load.file.info.failed"))
	}
	log.Printf("[app]: 镜像文件大小: %.2f MB", float64(fileInfo.Size())/1024/1024)

	// 2. 构建请求URL
	// Portainer 的镜像加载API: POST /api/endpoints/{id}/docker/images/load
	url := fmt.Sprintf("%s/endpoints/%d/docker/images/load", p.config.BaseURL, p.config.EndpointID)

	// 3. 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewReader(fileData))
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.create.failed"))
	}

	// 4. 设置请求头
	// Portainer 的镜像加载接口需要 multipart/form-data 或 application/x-tar
	req.Header.Set("Content-Type", "application/x-tar")
	if strings.HasSuffix(filePath, ".gz") || strings.HasSuffix(filePath, ".tgz") {
		req.Header.Set("Content-Type", "application/x-gzip")
	}

	// 5. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 6. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.load.request.send.failed"))
	}
	defer resp.Body.Close()

	// 7. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Portainer API 返回错误: %s, 响应: %s", resp.Status, string(body))
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.load.failed"))
	}

	// 8. 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	// 9. 解析响应，获取加载的镜像信息
	var loadResult []map[string]interface{}
	if err := json.Unmarshal(body, &loadResult); err != nil {
		// 如果解析失败，可能是简单的成功消息
		respStr := string(body)
		if strings.Contains(respStr, "Loaded image:") {
			log.Printf("[app]: 镜像加载成功: %s", respStr)
			// 提取镜像名称
			lines := strings.Split(respStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Loaded image:") {
					parts := strings.Split(line, "Loaded image:")
					if len(parts) > 1 {
						imageName := strings.TrimSpace(parts[1])
						log.Printf("[app]: 成功加载镜像: %s", imageName)

						// 验证镜像是否已加载
						if err := p.verifyImageExists(imageName); err != nil {
							log.Printf("[app]: 警告: 镜像验证失败: %v", err)
						} else {
							log.Printf("[app]: 镜像验证成功: %s", imageName)
						}
					}
				}
			}
			log.Printf("[app]: Portainer 本地镜像加载完成: %s", filePath)
			return nil
		}
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.parse.failed"))
	}

	// 10. 处理JSON格式的响应
	if len(loadResult) > 0 {
		for _, img := range loadResult {
			if stream, ok := img["stream"].(string); ok && stream != "" {
				stream = strings.TrimSpace(stream)
				if stream != "" {
					log.Printf("[app]: 镜像加载输出: %s", stream)

					// 提取镜像名称
					if strings.Contains(stream, "Loaded image:") {
						parts := strings.Split(stream, "Loaded image:")
						if len(parts) > 1 {
							imageName := strings.TrimSpace(parts[1])
							log.Printf("[app]: 成功加载镜像: %s", imageName)

							// 验证镜像是否已加载
							if err := p.verifyImageExists(imageName); err != nil {
								log.Printf("[app]: 警告: 镜像验证失败: %v", err)
							} else {
								log.Printf("[app]: 镜像验证成功: %s", imageName)
							}
						}
					}
				}
			}
		}
	}

	log.Printf("[app]: Portainer 本地镜像加载完成: %s", filePath)
	return nil
}

// Helper function for local image loading
func LoadImageFromLocal(imagePath string) error {
	config := DefaultPortainerConfig()

	adapter := NewPortainerAdapter(config, "")
	return adapter.LoadImageFromLocal(imagePath)
}

// Helper function to create a simple wrapper for backward compatibility
func PullImageFromRemote(imageUrl string) error {
	// 使用默认配置创建适配器
	config := DefaultPortainerConfig()

	adapter := NewPortainerAdapter(config, "")
	return adapter.PullImageFromRemote(imageUrl)
}

// DeployComposeYaml 使用 Portainer API 部署 Docker Compose 配置
func (p *PortainerAdapter) deployComposeYaml(composeYaml string, network string) error {
	if composeYaml == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.compose.empty"))
	}

	log.Printf("[app]: 开始通过 Portainer 部署 Docker Compose 配置")

	// 1. 构建请求体
	composeReq := map[string]interface{}{
		"Env":              []map[string]string{},
		"StackFileContent": composeYaml,
		"Prune":            true,
		"PullImage":        false,
	}

	reqJSON, err := json.Marshal(composeReq)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.compose.request.serialize.failed"))
	}

	// 2. 构建请求URL
	// Portainer 的 Stack 创建API: POST /api/stacks?endpointId={id}&method=string
	url := fmt.Sprintf("%s/stacks/%d?endpointId=%d&method=string&type=2", p.config.BaseURL, p.config.StackId, p.config.EndpointID)

	// 3. 创建请求
	req, err := http.NewRequest("PUT", url, bytes.NewReader(reqJSON))
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.compose.deploy.request.create.failed"))
	}

	// 4. 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 5. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 6. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.compose.deploy.request.send.failed"))
	}
	defer resp.Body.Close()

	// 7. 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("[app]: 部署yaml到portainer失败， url=%s, composeYaml=%s", url, composeYaml)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.compose.deploy.failed", string(body)))
	}

	// 8. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.compose.response.read.failed"))
	}

	var stackResp map[string]interface{}
	if err := json.Unmarshal(body, &stackResp); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.compose.response.parse.failed"))
	}

	// 9. 提取 Stack ID 和名称
	stackID, ok := stackResp["Id"].(float64)
	if !ok {
		stackID = 0
	}

	stackName, _ := stackResp["Name"].(string)
	if stackName == "" {
		stackName = "tier0"
	}

	log.Printf("[app]: Docker Compose Stack 部署成功: %s (ID: %.0f)", stackName, stackID)

	// 10. 等待 Stack 启动完成
	if err := p.waitForStackStart(int(stackID), stackName); err != nil {
		log.Printf("[app]: 警告: Stack 启动状态检查失败: %v", err)
		// 不因为状态检查失败而返回错误，Stack 可能仍在启动中
	}

	// 11. 将部署完的容器加入到指定 Docker 网络
	if network != "" {
		log.Printf("[app]: 开始将 Stack %s 中的容器加入到网络 %s", stackName, network)
		containerName, _ := extractContainerNameFromCompose(composeYaml)
		if err := p.joinStackContainersToNetwork(containerName, stackName, network); err != nil {
			log.Printf("[app]: 警告: 将容器加入到网络失败: %v", err)
			// 不因为网络加入失败而返回错误，记录日志但继续执行
		}
	}

	log.Printf("[app]: Portainer Docker Compose 部署完成")
	return nil
}

// extractContainerNameFromCompose 从 Compose YAML 中提取服务名称
func extractContainerNameFromCompose(composeYaml string) (string, error) {

	// 解析 YAML
	var composeConfig map[string]interface{}
	if err := yaml.Unmarshal([]byte(composeYaml), &composeConfig); err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.yaml.parse.failed"))
	}

	// 检查 services 部分
	services, ok := composeConfig["services"].(map[string]interface{})
	if !ok || len(services) == 0 {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.yaml.no.services"))
	}

	// 获取第一个服务名称
	for _, serviceConfig := range services {
		serviceMap, _ := serviceConfig.(map[string]interface{})
		if containerName, ok := serviceMap["container_name"].(string); ok && containerName != "" {
			return containerName, nil
		}
	}

	return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.yaml.no.service.name"))
}

// joinStackContainersToNetwork 将 Stack 中的所有容器加入到指定网络
func (p *PortainerAdapter) joinStackContainersToNetwork(containerName, stackName, networkName string) error {
	// 1. 获取 Stack 中的所有容器
	containers, err := p.getStackContainers(containerName)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.containers.fetch.failed"))
	}

	if len(containers) == 0 {
		log.Printf("[app]: Stack %s 中没有找到容器", stackName)
		return nil
	}

	log.Printf("[app]: Stack %s 中共有 %d 个容器需要加入到网络 %s", stackName, len(containers), networkName)

	// 2. 检查目标网络是否存在，如果不存在则创建
	networkID, err := p.getNetworkID(networkName)
	if err != nil {
		log.Printf("[app]: 网络 %s 不存在，尝试创建", networkName)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.create.failed"))
	}

	// 3. 将每个容器加入到网络
	successCount := 0
	for _, container := range containers {
		containerName := strings.TrimPrefix(container.Names[0], "/")
		log.Printf("[app]: 正在将容器 %s 加入到网络 %s", containerName, networkName)

		// 检查容器是否已经在网络中
		if p.isContainerInNetwork(container.ID, networkName) {
			log.Printf("[app]: 容器 %s 已经在网络 %s 中，跳过", containerName, networkName)
			successCount++
			continue
		}

		// 将容器连接到网络
		if err := p.connectContainerToNetwork(container.ID, networkID, containerName); err != nil {
			log.Printf("[app]: 容器 %s 加入网络失败: %v", containerName, err)
		} else {
			log.Printf("[app]: 容器 %s 已成功加入到网络 %s", containerName, networkName)
			successCount++
		}

		// 避免请求过于频繁
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("[app]: Stack %s 容器网络加入完成: %d/%d 个容器成功加入网络 %s",
		stackName, successCount, len(containers), networkName)

	if successCount == 0 {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.containers.join.failed"))
	}

	return nil
}

// getStackContainers 获取 Stack 中的所有容器
func (p *PortainerAdapter) getStackContainers(containerName string) ([]ContainerListItem, error) {
	// 1. 获取所有容器
	url := fmt.Sprintf("%s/endpoints/%d/docker/containers/json", p.config.BaseURL, p.config.EndpointID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.list.request.create.failed"))
	}

	if err := p.addAuthHeader(req); err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.list.request.send.failed"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	var allContainers []ContainerListItem
	if err := json.Unmarshal(body, &allContainers); err != nil {
		return nil, errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.list.parse.failed"))
	}

	// 2. 过滤出属于该 Stack 的容器
	// Docker Stack 容器名称格式: stackname_servicename.number.hash
	var stackContainers []ContainerListItem
	for _, container := range allContainers {
		for _, name := range container.Names {
			// 移除开头的斜杠
			cleanName := strings.TrimPrefix(name, "/")
			if containerName == cleanName {
				stackContainers = append(stackContainers, container)
			}
		}
	}

	return stackContainers, nil
}

// createCustomNetwork 创建自定义网络
func (p *PortainerAdapter) createCustomNetwork(networkName string) (string, error) {
	log.Printf("[app]: 正在创建网络: %s", networkName)

	// 构建网络创建请求体
	networkReq := map[string]interface{}{
		"Name":           networkName,
		"Driver":         "bridge",
		"CheckDuplicate": true,
		"Internal":       false,
		"Attachable":     true,
		"Ingress":        false,
		"IPAM": map[string]interface{}{
			"Driver": "default",
			"Config": []map[string]string{
				{
					"Subnet":  "172.21.0.0/16",
					"Gateway": "172.21.0.1",
				},
			},
		},
		"EnableIPv6": false,
		"Options":    map[string]string{},
		"Labels":     map[string]string{},
	}

	reqJSON, err := json.Marshal(networkReq)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.create.request_serialize_failed"))
	}

	// 构建请求URL
	url := fmt.Sprintf("%s/endpoints/%d/docker/networks/create", p.config.BaseURL, p.config.EndpointID)

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqJSON))
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.create.request.create.failed"))
	}

	req.Header.Set("Content-Type", "application/json")

	if err := p.addAuthHeader(req); err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.create.request.send.failed"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	var networkResp map[string]interface{}
	if err := json.Unmarshal(body, &networkResp); err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.network.create.response.parse.failed"))
	}

	// 提取网络ID
	if networkID, ok := networkResp["Id"].(string); ok && networkID != "" {
		log.Printf("[app]: 网络 %s 创建成功，ID: %s", networkName, networkID)
		return networkID, nil
	}

	// 如果响应中没有ID，尝试通过名称获取
	return p.getNetworkID(networkName)
}

// ContainerListItem 容器列表项结构体
type ContainerListItem struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	Command string   `json:"Command"`
	Created int64    `json:"Created"`
	State   string   `json:"State"`
	Status  string   `json:"Status"`
}

// waitForStackStart 等待 Stack 启动完成
func (p *PortainerAdapter) waitForStackStart(stackID int, stackName string) error {
	log.Printf("[app]: 等待 Stack 启动: %s (ID: %d)", stackName, stackID)

	timeout := time.Duration(p.config.Timeout) * time.Second
	startTime := time.Now()
	checkInterval := 5 * time.Second

	for {
		// 检查是否超时
		if time.Since(startTime) > timeout {
			return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.wait.start.timeout"))
		}

		// 获取 Stack 状态
		status, err := p.getStackStatus(stackID)
		if err != nil {
			log.Printf("[app]: 获取 Stack 状态失败: %v", err)
			time.Sleep(checkInterval)
			continue
		}

		log.Printf("[app]: Stack 状态: %s", status)

		// 检查状态
		switch status {
		case "active":
			log.Printf("[app]: Stack 启动成功: %s", stackName)
			return nil
		case "inactive":
			// 等待一段时间再检查
			time.Sleep(checkInterval)
			continue
		case "error":
			return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.wait.start.failed"))
		default:
			log.Printf("[app]: 未知的 Stack 状态: %s", status)
			time.Sleep(checkInterval)
		}

		// 如果已经等待了30秒，认为启动基本完成
		if time.Since(startTime) > 30*time.Second {
			log.Printf("[app]: Stack 启动基本完成（已等待30秒）: %s", stackName)
			return nil
		}
	}
}

// getStackStatus 获取 Stack 状态
func (p *PortainerAdapter) getStackStatus(stackID int) (string, error) {
	// 1. 构建请求URL
	url := fmt.Sprintf("%s/stacks/%d", p.config.BaseURL, stackID)

	// 2. 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.status.request.create.failed"))
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.stack.status.request.send.failed"))
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	// 6. 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.read.failed"))
	}

	var stackInfo map[string]interface{}
	if err := json.Unmarshal(body, &stackInfo); err != nil {
		return "", errors.Server.WithMsg(I18nUtils.GetMessage("portainer.response.parse.failed"))
	}

	// 7. 提取状态
	// Portainer 的 Stack 状态通常在 Status 字段
	if status, ok := stackInfo["Status"].(float64); ok {
		// Status: 1=active, 2=inactive
		if status == 1 {
			return "active", nil
		} else if status == 2 {
			return "inactive", nil
		}
	}

	// 或者检查是否有错误信息
	if errMsg, ok := stackInfo["ErrorMessage"].(string); ok && errMsg != "" {
		return "error", nil
	}

	// 默认返回 inactive
	return "inactive", nil
}

// Helper function for compose yaml loading
func DeployComposeYaml(yaml string, network string) error {
	config := DefaultPortainerConfig()

	adapter := NewPortainerAdapter(config, yaml)
	return adapter.deployComposeYaml(yaml, network)
}

// 环境变量辅助函数
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
	}
	return defaultValue
}

// DeleteContainer 删除容器
func (p *PortainerAdapter) DeleteContainer(containerNameOrID string, force bool, removeVolumes bool) error {
	if containerNameOrID == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.nameOrId.empty"))
	}

	log.Printf("[app]: 开始删除容器: %s (force: %v, removeVolumes: %v)", containerNameOrID, force, removeVolumes)

	// 1. 构建请求URL
	// Portainer 删除容器API: DELETE /api/endpoints/{id}/docker/containers/{id}?v={removeVolumes}&force={force}
	url := fmt.Sprintf("%s/endpoints/%d/docker/containers/%s", p.config.BaseURL, p.config.EndpointID, containerNameOrID)

	// 添加查询参数
	queryParams := fmt.Sprintf("?v=%v&force=%v", removeVolumes, force)
	fullURL := url + queryParams

	// 2. 创建请求
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.delete.request.create.failed"))
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.container.delete.request.send.failed"))
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	// 6. 处理响应
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[app]: 容器 %s 不存在或已被删除", containerNameOrID)
		return nil
	}

	log.Printf("[app]: 容器 %s 删除成功", containerNameOrID)
	return nil
}

// DeleteImage 删除镜像
func (p *PortainerAdapter) DeleteImage(imageNameOrID string, force bool, pruneChildren bool) error {
	if imageNameOrID == "" {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.nameOrId.empty"))
	}

	log.Printf("[app]: 开始删除镜像: %s (force: %v, pruneChildren: %v)", imageNameOrID, force, pruneChildren)

	// 1. 构建请求URL
	// Portainer 删除镜像API: DELETE /api/endpoints/{id}/docker/images/{id}?force={force}&noprune={!pruneChildren}
	url := fmt.Sprintf("%s/endpoints/%d/docker/images/%s", p.config.BaseURL, p.config.EndpointID, imageNameOrID)

	// 添加查询参数
	// noprune=true 表示不删除未标记的父镜像，noprune=false 表示删除未标记的父镜像
	noprune := !pruneChildren
	queryParams := fmt.Sprintf("?force=%v&noprune=%v", force, noprune)
	fullURL := url + queryParams

	// 2. 创建请求
	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.delete.request.create.failed"))
	}

	// 3. 添加认证头
	if err := p.addAuthHeader(req); err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.request.auth.failed"))
	}

	// 4. 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.image.delete.request.send.failed"))
	}
	defer resp.Body.Close()

	// 5. 检查响应状态
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		body, _ := io.ReadAll(resp.Body)
		return errors.Server.WithMsg(I18nUtils.GetMessage("portainer.api.error", string(body)))
	}

	// 6. 处理响应
	if resp.StatusCode == http.StatusNotFound {
		log.Printf("[app]: 镜像 %s 不存在或已被删除", imageNameOrID)
		return nil
	}

	// 7. 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[app]: 读取删除镜像响应失败: %v", err)
	} else if len(body) > 0 {
		// 解析删除结果
		var deleteResult []map[string]interface{}
		if err := json.Unmarshal(body, &deleteResult); err == nil {
			for _, result := range deleteResult {
				if untagged, ok := result["Untagged"].(string); ok && untagged != "" {
					log.Printf("[app]: 镜像已取消标记: %s", untagged)
				}
				if deleted, ok := result["Deleted"].(string); ok && deleted != "" {
					log.Printf("[app]: 镜像已删除: %s", deleted)
				}
			}
		}
	}

	log.Printf("[app]: 镜像 %s 删除成功", imageNameOrID)
	return nil
}
