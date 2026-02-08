# Docker Hub 配置指南

## 📋 概述

此 GitHub Action 已配置为仅推送到 Docker Hub，不再使用 GitHub Container Registry。

## ⚙️ 必需配置

在使用此 Action 之前，您必须在 GitHub 仓库中配置以下 Secrets：

### 1. 创建 Docker Hub 访问令牌

1. 登录到 [Docker Hub](https://hub.docker.com/)
2. 点击右上角头像 → **Account Settings**
3. 选择 **Security** 标签
4. 点击 **New Access Token**
5. 填写令牌名称（如 `github-actions`）
6. 选择权限：**Read, Write, Delete**
7. 点击 **Generate** 并复制生成的令牌

### 2. 在 GitHub 仓库中添加 Secrets

1. 进入您的 GitHub 仓库
2. 点击 **Settings** → **Secrets and variables** → **Actions**
3. 点击 **New repository secret** 添加以下 Secrets：

| Secret 名称 | 值 | 说明 |
|-------------|-----|------|
| `DOCKERHUB_USERNAME` | 您的 Docker Hub 用户名 | 用于登录 Docker Hub |
| `DOCKERHUB_TOKEN` | 上面生成的访问令牌 | 用于认证推送操作 |

## 🚀 使用方法

配置完成后，您可以：

### 发布新版本
```bash
git tag v1.1.0
git push origin v1.1.0
```

### 拉取镜像
```bash
# 拉取特定版本
docker pull 您的用户名/notify:v1.1.0

# 拉取最新版本
docker pull 您的用户名/notify:latest
```

### 运行容器
```bash
docker run -d \
  --name notify-app \
  -p 8088:8088 \
  -v ./config:/config \
  -e TZ=Asia/Shanghai \
  您的用户名/notify:latest
```

## 📦 镜像标签

推送到 Docker Hub 的镜像会包含以下标签：
- `您的用户名/notify:v1.1.0` - 完整版本号
- `您的用户名/notify:1.1.0` - 不带 v 前缀
- `您的用户名/notify:1.1` - 主次版本号
- `您的用户名/notify:1` - 主版本号
- `您的用户名/notify:latest` - 最新版本

## ⚠️ 注意事项

1. **必须配置 Secrets**: 如果没有配置 `DOCKERHUB_USERNAME` 和 `DOCKERHUB_TOKEN`，Action 将会失败
2. **访问令牌权限**: 确保 Docker Hub 访问令牌有 `Read, Write, Delete` 权限
3. **仓库名称**: 镜像名称将是 `您的用户名/notify`
4. **多架构支持**: 自动构建 `linux/amd64` 和 `linux/arm64` 架构

## 🔍 故障排除

### Action 失败并提示缺少凭据
- 检查是否已正确添加 `DOCKERHUB_USERNAME` 和 `DOCKERHUB_TOKEN` Secrets
- 确认 Secret 名称拼写正确（区分大小写）

### 推送失败
- 检查 Docker Hub 访问令牌是否有效
- 确认访问令牌有足够的权限
- 检查 Docker Hub 用户名是否正确

### 镜像名称错误
- 确认 `DOCKERHUB_USERNAME` 是您的 Docker Hub 用户名
- 检查 Docker Hub 上是否已有同名仓库
