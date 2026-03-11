package app

import (
	"backend/share/app/adapter"
	"backend/share/app/util"
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/share/app/model" // 请根据实际路径修改

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func install(featureConfig *model.NewFeatureModel) error {

	return nil
}

func createAppForZipFile(zipPath string) error {
	dir, err := util.ExtractZipToDir(zipPath, util.DefaultTempDir)
	if err != nil {
		return err
	}
	defer os.Remove(dir)
	if err := ValidateExtractedZip(dir); err != nil {
		return err
	}
	composeConfig, err := model.LoadAndValidateComposeFile(filepath.Join(dir, "compose.yaml"))
	if err != nil {
		return err
	}
	appConfig, err := model.LoadAndValidateAppConfig(filepath.Join(dir, "app.yaml"))
	if err != nil {
		return err
	}
	requireConfig, err := model.LoadAndValidateRequirementConfig(filepath.Join(dir, "requirement.yaml"))
	if err != nil {
		return err
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		log.Fatalf("初始化 Docker 客户端失败: %v", err)
		return err
	}
	defer cli.Close()

	installDir, err := installAppToDirectory(appConfig, dir)
	if err != nil {
		return err
	}
	installSuccess := false
	defer func() {
		if !installSuccess {
			// TODO 清理现场
		}
	}()
	// 1. 申请数据库资源
	if requireConfig != nil && requireConfig.Requirements.HasDatabases() {
		if err := applyDatabase(requireConfig, installDir); err != nil {
			return err
		}
	}

	for _, c := range composeConfig.Service.Containers {
		// 2. 加载镜像
		if err := loadImages(cli, c, installDir); err != nil {
			return err
		}
		// 3. 读取环境变量配置
		envs, _ := LoadConfigFromExtractDir(installDir)
		// 4. 根据 composeModel.go 配置创建并启动容器
		if err := runContainer(ctx, cli, c, requireConfig, envs); err != nil {
			log.Printf("运行容器 %s 失败: %v", c.Name, err)
			return err
		}
		// TODO 检测是否启动完成
		// 5. TODO 注册菜单
	}
	installSuccess = true
	return nil
}

func loadImages(cli *client.Client, c model.Container, dir string) error {
	imagesDir := filepath.Join(dir, "images")
	imageFoundLocal := false
	ctx := context.Background()
	// 1. 检查 images 目录是否存在并尝试加载 tar
	if _, err := os.Stat(imagesDir); err == nil {
		err = filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(info.Name(), ".tar") {
				log.Printf("尝试加载本地镜像文件: %s", path)
				if loadErr := loadLocalImage(ctx, cli, path); loadErr == nil {
					imageFoundLocal = true
				}
			}
			return nil
		})
	}

	// 2. 如果本地没加载成功，尝试在线拉取
	if !imageFoundLocal {
		log.Printf("本地未找到有效镜像，尝试在线拉取: %s", c.Image)
		out, err := cli.ImagePull(ctx, c.Image, image.PullOptions{})
		if err != nil {
			log.Printf("拉取镜像失败: %v", err)
			return err
		}
		defer out.Close()
		io.Copy(os.Stdout, out)
	}
	return nil
}

// 申请数据库并写入到db-config.ini
func applyDatabase(rc *model.RequirementConfig, extractDir string) error {
	log.Printf("开始处理数据库需求，共 %d 个数据库", len(rc.Requirements.Databases))

	// 创建PostgreSQL适配器
	pgConfig := adapter.DefaultPostgresConfig()

	pgAdapter := adapter.NewPostgresAdapter(pgConfig)

	// 连接数据库
	if err := pgAdapter.Connect(); err != nil {
		return err
	}
	defer pgAdapter.Close()

	// 测试连接
	if err := pgAdapter.TestConnection(); err != nil {
		return err
	}

	// 根据requirement创建数据库
	dbInfoList, err := pgAdapter.CreateDatabasesFromRequirements(&rc.Requirements)
	if err != nil {
		return err
	}
	return writeDatabaseConfigToFile(dbInfoList, extractDir)
}

// writeDatabaseConfigToFile 将数据库信息写入配置文件
func writeDatabaseConfigToFile(dbInfoList []*model.DatabaseInfo, extractDir string) error {
	if len(dbInfoList) == 0 {
		return nil
	}

	// 写入主配置文件
	mainConfigPath := filepath.Join(extractDir, "db-config.ini")
	file, err := os.Create(mainConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入每个数据库的配置
	for _, dbInfo := range dbInfoList {
		sectionName := fmt.Sprintf("database:%s", dbInfo.Database)
		fmt.Fprintf(file, "[%s]\n", sectionName)
		writeDatabaseSection(file, dbInfo)
		fmt.Fprintln(file)
	}
	fmt.Fprintln(file)
	log.Printf("数据库配置文件已写入: %s", mainConfigPath)
	return nil
}

func writeDatabaseSection(writer io.Writer, dbInfo *model.DatabaseInfo) {
	fmt.Fprintf(writer, "SUPOS_%s_DBNAME=%s\n", dbInfo.Database, dbInfo.Database)
	fmt.Fprintf(writer, "SUPOS_%s_USERNAME=%s\n", dbInfo.Database, dbInfo.Username)
	fmt.Fprintf(writer, "SUPOS_%s_PASSWORD=%s\n", dbInfo.Database, dbInfo.Password)
	fmt.Fprintf(writer, "SUPOS_%s_HOST=%s\n", dbInfo.Database, dbInfo.Host)
	fmt.Fprintf(writer, "SUPOS_%s_PORT=%s\n", dbInfo.Database, dbInfo.Port)
}

// 创建应用安装目录并拷贝文件
func installAppToDirectory(appConfig *model.AppConfig, extractDir string) (string, error) {
	if appConfig == nil {
		return "", fmt.Errorf("app config is nil")
	}

	if extractDir == "" {
		return "", fmt.Errorf("extract directory is empty")
	}

	// 1. 构建安装目录路径
	installDir := filepath.Join(util.AppInstalledDir, appConfig.VendorName, appConfig.Name)

	// 2. 确保父目录存在
	parentDir := filepath.Dir(installDir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory %s: %v", parentDir, err)
	}

	// 3. 清空目录（如果已存在）
	if _, err := os.Stat(installDir); err == nil {
		log.Printf("安装目录已存在，清空目录: %s", installDir)
		if err := os.RemoveAll(installDir); err != nil {
			return "", fmt.Errorf("failed to clear existing directory %s: %v", installDir, err)
		}
	}

	// 4. 创建安装目录
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create install directory %s: %v", installDir, err)
	}

	log.Printf("创建应用安装目录: %s", installDir)

	// 5. 拷贝解压的文件到安装目录
	if err := util.CopyDirectory(extractDir, installDir); err != nil {
		// 如果拷贝失败，清理安装目录
		os.RemoveAll(installDir)
		return "", fmt.Errorf("failed to copy files to install directory: %v", err)
	}

	log.Printf("成功安装应用到: %s", installDir)

	return installDir, nil
}

// checkContainerHealth 检查容器健康状态
func checkContainerHealth(ctx context.Context, cli *client.Client, containerName string, timeout time.Duration) error {
	startTime := time.Now()

	for {
		// 检查是否超时
		if time.Since(startTime) > timeout {
			return fmt.Errorf("容器健康检查超时: %s", containerName)
		}

		// 获取容器状态
		containerJSON, err := cli.ContainerInspect(ctx, containerName)
		if err != nil {
			return fmt.Errorf("检查容器状态失败: %v", err)
		}

		// 检查容器是否在运行
		if !containerJSON.State.Running {
			return fmt.Errorf("容器未运行: %s", containerName)
		}

		// 检查健康状态（如果配置了健康检查）
		if containerJSON.State.Health != nil {
			if containerJSON.State.Health.Status == "healthy" {
				log.Printf("容器 %s 健康检查通过", containerName)
				return nil
			} else if containerJSON.State.Health.Status == "unhealthy" {
				return fmt.Errorf("容器健康检查失败: %s", containerName)
			}
			// 如果状态是 "starting" 或空，继续等待
		} else {
			// 如果没有配置健康检查，等待一段时间后认为启动成功
			if time.Since(startTime) > 30*time.Second {
				log.Printf("容器 %s 启动完成（未配置健康检查）", containerName)
				return nil
			}
		}

		// 等待一段时间再检查
		time.Sleep(2 * time.Second)
	}
}

// 加载本地 tar 镜像
func loadLocalImage(ctx context.Context, cli *client.Client, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := cli.ImageLoad(ctx, f)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 必须读取 body 以确保加载完成
	io.Copy(io.Discard, resp.Body)
	return nil
}

// 根据配置运行容器
func runContainer(ctx context.Context, cli *client.Client, c model.Container, rc *model.RequirementConfig, configEnvVars []string) error {
	// 映射端口配置
	portSet := nat.PortSet{}
	portBindings := nat.PortMap{}
	for _, p := range c.Ports {
		port := nat.Port(fmt.Sprintf("%d/tcp", p.ContainerPort))
		portSet[port] = struct{}{}
		// 默认将容器端口映射到宿主机同端口，可根据 model 调整
		portBindings[port] = []nat.PortBinding{{HostPort: fmt.Sprintf("%d", p.ContainerPort)}}
	}

	networkName := "tier0_edge_network"

	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {
				// 可以设置网络别名等
				Aliases: []string{c.Name},
			},
		},
	}

	// 创建主机配置
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Resources: container.Resources{
			Memory:     parseMemory(c.Resources.Limits.Memory),
			MemorySwap: parseMemory(c.Resources.Limits.Memory), // 设置swap内存
		},
		RestartPolicy: container.RestartPolicy{
			Name:              "unless-stopped", // 自动重启策略
			MaximumRetryCount: 0,
		},
	}
	// 挂载目录
	mounts, err2 := buildVolumeMount(c, rc)
	if err2 != nil {
		return err2
	}
	// 设置挂载
	if len(mounts) > 0 {
		hostConfig.Mounts = mounts
	}

	// 创建容器
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        c.Image,
		ExposedPorts: portSet,
		Env:          configEnvVars,
	}, hostConfig, networkingConfig, nil, c.Name)
	if err != nil {
		return err
	}

	// 启动容器
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("启动容器失败: %v", err)
	}
	log.Printf("容器 %s 启动成功，容器ID: %s，网络: %s",
		c.Name, resp.ID[:12], networkName)

	return nil
}

func buildVolumeMount(c model.Container, rc *model.RequirementConfig) ([]mount.Mount, error) {
	var mounts []mount.Mount
	// 处理VolumeMounts，需要和requirement匹配
	for _, volume := range c.VolumeMounts {
		if volume.Name != "" && volume.MountPath != "" {
			// 检查volume.name是否能在requirement的volumes中找到
			found := false
			var requirementVolume *model.VolumeRequirement

			// 遍历requirement中的volumes
			for i, reqVolume := range rc.Requirements.Volumes {
				if reqVolume.Name == volume.Name {
					found = true
					requirementVolume = &rc.Requirements.Volumes[i] // 使用指针
					break
				}
			}

			// 如果没找到，记录警告并继续下一个循环
			if !found {
				log.Printf("警告: 容器 %s 的卷 %s 在requirement.yaml中未定义，跳过挂载",
					c.Name, volume.Name)
				continue
			}

			// 创建本地目录（如果不存在）
			createVolumeDirectory(requirementVolume.LocalPath)

			// 创建挂载配置
			mountConfig := mount.Mount{
				Type:        mount.TypeBind,
				Source:      requirementVolume.LocalPath,
				Target:      volume.MountPath,
				ReadOnly:    false, // 默认读写，可以根据需要调整
				Consistency: mount.ConsistencyDefault,
			}

			mounts = append(mounts, mountConfig)

			// 记录详细信息
			log.Printf("设置卷挂载: %s -> %s", requirementVolume.LocalPath, volume.MountPath)
			log.Printf("  卷名称: %s", volume.Name)
			log.Printf("  卷大小: %s", requirementVolume.Size)
			log.Printf("  资源类型: %s", requirementVolume.ResourceType)
		}
	}
	return mounts, nil
}

// createVolumeDirectory 创建卷目录
func createVolumeDirectory(localPath string) {

	// 创建目录
	if err := os.MkdirAll(localPath, 0755); err != nil {
		log.Printf("警告: 创建卷目录失败 %s: %v", localPath, err)
	}

	// 设置权限
	os.Chmod(localPath, 0755)

}

// createLogDirectory 创建日志目录
func createLogDirectory(containerName string) string {
	// 构建日志目录路径
	logDir := filepath.Join("/var/log/docker", containerName)

	// 创建目录
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("警告: 创建日志目录失败 %s: %v", logDir, err)
	}

	return logDir
}

// parseMemory 解析内存字符串为字节数
// 支持格式: 512Mi, 512MB, 1Gi, 1GB, 512M, 1G 等（忽略大小写）
func parseMemory(mem string) int64 {
	if mem == "" {
		mem = "512MB"
	}

	// 转换为小写并移除空格
	mem = strings.ToLower(strings.TrimSpace(mem))

	// 定义单位映射
	unitMap := map[string]int64{
		"k":  1024,
		"kb": 1024,
		"m":  1024 * 1024,
		"mb": 1024 * 1024,
		"g":  1024 * 1024 * 1024,
		"gb": 1024 * 1024 * 1024,
		"t":  1024 * 1024 * 1024 * 1024,
		"tb": 1024 * 1024 * 1024 * 1024,
		"ki": 1024,
		"mi": 1024 * 1024,
		"gi": 1024 * 1024 * 1024,
		"ti": 1024 * 1024 * 1024 * 1024,
	}

	// 提取数字部分
	var numStr string
	var unitStr string

	// 遍历字符串，分离数字和单位
	for i, r := range mem {
		if (r >= '0' && r <= '9') || r == '.' {
			numStr += string(r)
		} else {
			unitStr = mem[i:]
			break
		}
	}

	// 如果没有单位，假设是字节
	if unitStr == "" {
		unitStr = "m"
	}

	// 解析数字
	var num float64
	_, err := fmt.Sscanf(numStr, "%f", &num)
	if err != nil {
		log.Printf("无法解析内存值: %s, 使用默认值 512MB", mem)
		return 512 * 1024 * 1024 // 默认 512MB
	}

	// 查找单位
	multiplier, ok := unitMap[unitStr]
	if !ok {
		log.Printf("未知的内存单位: %s, 使用默认单位 MB", unitStr)
		multiplier = 1024 * 1024 // 默认 MB
	}

	// 计算字节数
	result := int64(num * float64(multiplier))

	return result
}

// ConfigEntry 配置文件条目
type ConfigEntry struct {
	Key   string
	Value string
}

// ReadConfigINI 读取config.ini文件
// filePath: config.ini文件路径
// 返回: 配置条目列表和错误信息
func ReadConfigINI(filePath string) ([]ConfigEntry, error) {
	var entries []ConfigEntry

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return entries, nil // 文件不存在，返回空列表
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	// 读取文件内容
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// 解析键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			log.Printf("警告: 第%d行格式错误，跳过: %s", lineNumber, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 移除值中的引号（如果存在）
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') ||
			(value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		entries = append(entries, ConfigEntry{
			Key:   key,
			Value: value,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	log.Printf("从 %s 读取了 %d 个配置项", filePath, len(entries))
	return entries, nil
}

// ConfigEntriesToEnv 将配置条目转换为环境变量格式
func ConfigEntriesToEnv(entries []ConfigEntry) []string {
	var envVars []string

	for _, entry := range entries {
		envVar := fmt.Sprintf("%s=%s", strings.ToUpper(entry.Key), entry.Value)
		envVars = append(envVars, envVar)
	}

	return envVars
}

// LoadConfigFromExtractDir 从解压目录加载config.ini和db-config.ini
func LoadConfigFromExtractDir(extractDir string) ([]string, error) {
	configPath := filepath.Join(extractDir, "config.ini")
	dbConfigPath := filepath.Join(extractDir, "db-config.ini")

	var allEntries []ConfigEntry

	// 1. 读取config.ini文件（如果存在）
	if _, err := os.Stat(configPath); err == nil {
		entries, err := ReadConfigINI(configPath)
		if err != nil {
			log.Printf("读取config.ini失败: %v", err)
			// 不因为读取失败而中断，继续处理
		} else {
			allEntries = append(allEntries, entries...)
			log.Printf("从config.ini读取了 %d 个配置项", len(entries))
		}
	} else {
		log.Printf("config.ini文件不存在: %s", configPath)
	}

	// 2. 读取db-config.ini文件（如果存在）
	if _, err := os.Stat(dbConfigPath); err == nil {
		entries, err := ReadConfigINI(dbConfigPath)
		if err != nil {
			log.Printf("读取db-config.ini失败: %v", err)
			// 不因为读取失败而中断，继续处理
		} else {
			allEntries = append(allEntries, entries...)
			log.Printf("从db-config.ini读取了 %d 个配置项", len(entries))
		}
	} else {
		log.Printf("db-config.ini文件不存在: %s", dbConfigPath)
	}

	// 3. 如果没有读取到任何配置，返回空列表
	if len(allEntries) == 0 {
		log.Println("未找到任何配置文件")
		return []string{}, nil
	}

	// 4. 将配置条目转换为环境变量
	envVars := ConfigEntriesToEnv(allEntries)
	log.Printf("总共生成了 %d 个环境变量", len(envVars))

	return envVars, nil
}
