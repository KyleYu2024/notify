# GitHub Actions - 插件构建工作流

## 📋 概述

这个工作流专门用于构建项目中的Go插件，支持多平台交叉编译。

## 🚀 触发方式

### 手动触发
工作流仅支持手动触发，在 GitHub 仓库的 Actions 页面手动运行并设置以下参数：

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| `plugin_name` | 指定要构建的插件名称 | 全部插件 | `emby` |
| `platforms` | 目标平台列表 | `all` | `darwin-amd64,linux-amd64` |
| `debug_mode` | 启用Debug模式 | `false` | `true` |

## 🏗️ 构建平台

| 运行环境 | 支持的目标平台 |
|----------|----------------|
| Ubuntu | `linux-amd64`, `linux-arm64` |
| macOS | `darwin-amd64`, `darwin-arm64` |

## 📦 构建产物

构建完成后，插件文件会被打包为 Artifacts：
- **命名格式**: `plugins-{运行环境}-{运行编号}`
- **保留期**: 30天
- **文件格式**: `plugin-{平台}.so`

### 下载构建产物
1. 进入 Actions 页面
2. 点击对应的工作流运行记录
3. 在页面底部的 "Artifacts" 部分下载

## 📁 文件结构

构建产物的目录结构：
```
artifacts/
├── emby/
│   ├── plugin-darwin-amd64.so
│   ├── plugin-darwin-arm64.so
│   ├── plugin-linux-amd64.so
│   └── plugin-linux-arm64.so
└── demo/
    ├── plugin-darwin-amd64.so
    ├── plugin-darwin-arm64.so
    ├── plugin-linux-amd64.so
    └── plugin-linux-arm64.so
```

## 🔧 使用示例

### 构建所有插件（所有平台）
```bash
# 在 GitHub Actions 页面手动触发，所有参数使用默认值
# plugin_name: (留空)
# platforms: all  
# debug_mode: false
```

### 构建特定插件
```bash
# 手动触发时设置:
# plugin_name: emby
# platforms: all
# debug_mode: false
```

### 构建特定平台
```bash
# 手动触发时设置:
# plugin_name: (留空)
# platforms: darwin-amd64,linux-amd64
# debug_mode: false
```

### Debug模式构建
```bash
# 手动触发时设置:
# plugin_name: emby
# platforms: darwin-amd64
# debug_mode: true
```

## ⚠️ 注意事项

1. **手动触发**: 工作流只能手动触发，适合按需构建插件
2. **CGO依赖**: 插件构建需要CGO支持，交叉编译会自动处理依赖问题
3. **Docker要求**: Linux平台交叉编译依赖Docker，GitHub Actions已预装
4. **并行构建**: 不同平台会并行构建以提高效率
5. **构建失败**: 如果某个平台构建失败，不会影响其他平台的构建

## 🔍 故障排除

### 常见问题
1. **构建失败**: 查看具体的构建日志，检查Go模块依赖
2. **找不到插件**: 确保插件目录包含 `plugin.go` 文件
3. **交叉编译错误**: 检查是否有平台特定的代码需要条件编译

### 调试建议
1. 先在本地使用 `build-plugin.sh` 脚本测试
2. 启用 Debug 模式获取更详细的构建信息
3. 检查插件的Go模块依赖是否完整
