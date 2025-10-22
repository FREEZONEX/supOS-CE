# SupOS Edge Frontend

现代化的前端项目，基于 pnpm monorepo 架构，使用 React 19、Vite、TypeScript 等现代技术栈。

## 🚀 技术栈

- **框架**: React + React Router
- **构建工具**: Vite + Turbo + tsdown
- **语言**: TypeScript
- **UI组件**: Ant Design
- **状态管理**: Zustand
- **包管理**: pnpm
- **代码质量**: ESLint + Prettier + Stylelint
- **图表**: AntV X6
- **AI集成**: CopilotKit + OpenAI

## 📁 项目结构

```
frontend/
├── apps/                    # 应用目录
│   ├── web/                # 主Web应用
│   ├── services/           # 服务应用
│   └── services-express/   # Express服务应用
├── packages/               # 共享包
│   ├── scripts/            # 构建脚本
│   └── typescript-config/  # TypeScript配置
├── plugins/               # 插件 - 模块联邦实现
│   └──alert/            # 报警模块
├── scripts/                # 项目脚本
└── 配置文件
```

## 🛠️ 开发环境要求

- Node.js >= 22.20.0
- pnpm >= 10.13.1

## 📦 安装依赖

```bash
# 切换目录
cd ./frontend
```

```bash
# 安装项目依赖
pnpm install
```

## 🏃‍♂️ 开发命令

### 仅启动Web应用

```bash
pnpm dev:web
```

### 仅启动services应用

```bash
pnpm dev:servicesExpress
```

### 构建项目

```bash

# 仅构建Web应用
pnpm build:web

# 仅构建services应用
pnpm build:servicesExpress
```

### 代码质量

```bash
# 代码检查
pnpm lint

# 国际化相关
pnpm intl:once      # 一次性国际化处理
pnpm intl:watch     # 监听模式国际化处理
pnpm properties:convert  # 属性文件转换
```

### 清理

```bash
pnpm clean  # 清理node_modules
```

## 🔧 配置说明

### 包管理配置

- `pnpm-workspace.yaml`: 定义 monorepo 工作区
- `turbo.json`: Turbo 构建配置

### 代码规范

- `eslint.config.mjs`: ESLint 配置
- `.prettierrc`: Prettier 配置
- `stylelint.config.js`: Stylelint 配置
- `commitlint.config.cjs`: Git 提交规范

## 🌐 应用说明

### Web 应用 (`apps/web`)

主前端应用，包含：

- 现代化的用户界面
- 模块联邦支持
- 国际化支持

### 服务应用 - Express (`apps/services-express`)

后端服务应用，目前提供copilotkit接口和健康检查接口。

### 服务应用 (`apps/services`)

基于hono.js实现的版本

## 📊 特性

### 核心特性

- 🏗️ **Monorepo 架构**: 使用 pnpm workspace 管理多包
- ⚡ **快速构建**: 基于 Vite 的快速开发体验
- 🎯 **TypeScript**: 完整的类型支持
- 🎨 **现代化 UI**: 基于 Ant Design 5.x
- 🌍 **国际化**: 完整的国际化支持
- 🤖 **AI 集成**: CopilotKit AI 助手

### 开发体验

- 🔍 **热重载**: 开发时热重载支持
- 📏 **代码规范**: 统一的代码规范配置
- 🧪 **测试就绪**: 完善的测试环境配置
- 🐛 **调试友好**: 完整的调试支持

## 🚀 部署

项目支持多种部署方式：

### Docker 部署

使用提供的 Dockerfile-CN 进行容器化部署。

### 传统部署

构建后部署静态文件到 Web 服务器。

## 🤝 贡献指南

1. 遵循项目的代码规范
2. 提交前运行 `pnpm lint` 确保代码质量
3. 使用规范的 Git 提交信息
4. 确保所有测试通过

## 📄 许可证

本项目基于 SupOS Edge 项目的许可证。
