# Demo MCP Server

一个基于 Model Context Protocol (MCP) 的演示服务器，提供天气查询功能。

## 功能特性

- 🌤️ **天气查询**：获取指定城市的实时天气信息
- 🔄 **多种传输协议**：支持 stdio、SSE 和 HTTP 流式传输
- 🛠️ **标准化接口**：基于 MCP 协议，兼容各种 MCP 客户端
- 📦 **易于部署**：提供 npm 包和 CLI 工具

## 安装

### 作为全局工具安装

```bash
npm install -g @supos-os-edge/demo-mcp-server
```

### 作为项目依赖安装

```bash
npm install @supos-os-edge/demo-mcp-server
```

## 使用方法

### 1. 标准输入输出模式 (stdio)

```bash
# 使用全局安装
npx @supos-os-edge/demo-mcp-server stdio

# 或使用本地安装
node_modules/.bin/@supos-os-edge/demo-mcp-server stdio
```

### 2. SSE 服务器模式

```bash
# 启动 SSE 服务器
npx @supos-os-edge/demo-mcp-server sse
```

### 3. HTTP 流式传输模式

```bash
# 启动 HTTP 流式传输服务器
npx @supos-os-edge/demo-mcp-server streamableHttp
```

## 配置

### 在 Claude Desktop 中配置

在 Claude Desktop 的配置文件中添加：

```json
{
  "mcpServers": {
    "demo-mcp-server": {
      "command": "npx",
      "args": ["-y", "@supos-os-edge/demo-mcp-server"]
    }
  }
}
```

### 环境变量

- `PORT`：服务器端口（默认为 3000）
- `HOST`：服务器主机（默认为 localhost）

## 可用工具

### get_weather

获取指定城市的天气信息。

**参数：**

- `city`：城市名称（例如：北京、上海、杭州）

**返回示例：**

```
🌤️ 天气信息：

📍 位置：浙江省 - 杭州市
🌡️ 温度：25°C
🌤️ 天气：多云 (多云)
💧 湿度：65%
💨 风向：东南风 3级
🌊 气压：1013hPa
💦 降水量：0mm
🕐 更新时间：2024-01-01 12:00:00

天气代码：01
```

## 开发

### 项目结构

```
src/
├── index.ts              # 主入口文件
├── server/
│   └── index.ts          # MCP 服务器实现
├── tools/
│   └── weather.ts        # 天气工具实现
├── transport/
│   ├── stdio.ts          # 标准输入输出传输
│   ├── sse.ts            # SSE 传输
│   └── streamableHttp.ts # HTTP 流式传输
└── utils/
    └── log.ts            # 日志工具
```

### 构建项目

```bash
# 安装依赖
pnpm install

# 开发模式（监听文件变化）
pnpm run dev

# 构建项目
pnpm run build

# 启动服务器
pnpm start
```

### 测试

```bash
# 使用 MCP 检查器测试
pnpm run inspector
```

## API 文档

### 天气数据来源

本服务使用腾讯天气 API 获取实时天气数据。

### 错误处理

- 城市不存在：返回友好的错误提示
- 网络错误：提供重试建议
- 参数错误：详细的参数验证信息

## 许可证

ISC License

## 贡献

欢迎提交 Issue 和 Pull Request！

## 更新日志

### v1.0.1

- 初始版本发布
- 实现基本的天气查询功能
- 支持多种传输协议

## 技术支持

如有问题，请提交 Issue 或联系开发团队。
