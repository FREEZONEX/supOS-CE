# @supos-os-edge/scripts

平台执行脚本包

## 功能

- **intl:once**: 国际化处理
- **intl:watch**: 国际化处理
- **json-to-properties**: JSON文件转换为Properties格式

## 使用方法

### 安装依赖

```bash
npm install
```

### 构建项目

```bash
npm run build
```

### 开发模式（监听文件变化自动编译）

```bash
npm run dev
```

### 运行脚本

```bash
# 使用编译后的脚本
supos-scripts <script-name>

# 示例
supos-scripts intl:once
supos-scripts intl:watch
supos-scripts json-to-properties --prefix=myapp --config=config.json
supos-scripts json-to-properties --prefix=myapp
supos-scripts json-to-properties
```
