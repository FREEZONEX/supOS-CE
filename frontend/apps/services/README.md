# Services App

这是一个基于Hono框架的Node.js服务器应用。

## 功能特性

- ✅ 优雅的错误处理和全局异常捕获
- ✅ 环境变量配置支持
- ✅ 健康检查端点 (`/health`)
- ✅ 请求日志记录
- ✅ CORS跨域支持
- ✅ 优雅关闭处理
- ✅ TypeScript支持

## 快速开始

### 安装依赖

```bash
pnpm install
```

### 开发模式

```bash
pnpm dev
```

### 生产模式

```bash
# 构建
pnpm build

# 启动
pnpm start
```

## 环境配置

复制 `.env.example` 为 `.env` 并配置相应环境变量：

```bash
cp .env.example .env
```

### 环境变量说明

- `PORT`: 服务器端口 (默认: 3001)
- `NODE_ENV`: 运行环境 (development/production)
- `ALLOWED_ORIGINS`: 允许的跨域源
- `LOG_LEVEL`: 日志级别

## API端点

### 根路径

- **GET** `/`
  - 返回欢迎信息和服务器状态

### 健康检查

- **GET** `/health`
  - 返回服务器健康状态
  - 包含环境信息和时间戳

## 项目结构

```
src/
  index.ts          # 主入口文件
dist/               # 编译输出目录
.env.example        # 环境变量示例
package.json        # 项目配置
README.md          # 项目说明
```

## 开发说明

- 使用 TypeScript 进行类型安全开发
- 支持热重载开发模式
- 内置错误处理和日志记录
- 支持优雅关闭，确保服务稳定运行

## 部署说明

1. 设置生产环境变量
2. 构建项目: `pnpm build`
3. 启动服务: `pnpm start`

服务器将在配置的端口上运行，并自动处理跨域请求和错误。
