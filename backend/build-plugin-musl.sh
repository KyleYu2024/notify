#!/bin/bash

# musl 插件构建脚本
# 使用 Docker + Alpine Linux 构建 musl 版本的插件

# 默认参数
DEFAULT_MUSL_PLATFORMS="linux/amd64,linux/arm64"
DEFAULT_LDFLAGS="-s -w"

# 帮助信息
show_help() {
    echo "Usage: $0 [options] [plugin_name_or_path]"
    echo ""
    echo "Options:"
    echo "  -p, --platforms    指定 musl 构建平台 (例如: linux/amd64,linux/arm64)"
    echo "                     默认: $DEFAULT_MUSL_PLATFORMS"
    echo "  -l, --ldflags      指定构建的 ldflags (默认: $DEFAULT_LDFLAGS)"
    echo "  -d, --debug        启用调试模式 (不压缩二进制文件)"
    echo "  -h, --help         显示帮助信息"
    echo ""
    echo "Examples:"
    echo "  $0                                         # 构建 ./plugins/ 下所有插件"
    echo "  $0 ./plugins                               # 构建指定目录下所有插件"
    echo "  $0 emby                                    # 构建 ./plugins/emby 插件"
    echo "  $0 ./plugins/emby                         # 使用完整路径构建"
    echo "  $0 -p linux/amd64 emby                    # 只构建 amd64 平台"
    echo "  $0 -d                                      # 调试模式构建所有插件"
    echo ""
    echo "说明:"
    echo "  - 不指定参数时，自动构建 ./plugins/ 下所有插件"
    echo "  - 插件名会自动在 ./plugins/ 目录下查找"
    echo "  - 也可以直接指定完整的插件路径"
    echo "  - 使用 Docker + Alpine Linux (golang:1.25.1-alpine) 构建"
    echo "  - 输出文件格式: plugin-{GOOS}-{GOARCH}-musl.so"
    echo "  - 需要 Docker 环境支持"
}

# 解析命令行参数
MUSL_PLATFORMS="$DEFAULT_MUSL_PLATFORMS"
LDFLAGS="$DEFAULT_LDFLAGS"
DEBUG_MODE=0
PLUGIN_INPUT=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--platforms)
            MUSL_PLATFORMS="$2"
            shift 2
            ;;
        -l|--ldflags)
            LDFLAGS="$2"
            shift 2
            ;;
        -d|--debug)
            DEBUG_MODE=1
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            PLUGIN_INPUT="$1"
            shift
            ;;
    esac
done

# 确定构建模式
BUILD_MODE="single"  # single 或 batch
TARGET_DIR="./plugins"

if [ -z "$PLUGIN_INPUT" ]; then
    # 没有指定参数，默认批量构建 ./plugins 目录
    BUILD_MODE="batch"
    echo "📦 批量构建模式: 将构建 $TARGET_DIR 下所有插件"
elif [ "$PLUGIN_INPUT" = "./plugins" ] || [ "$PLUGIN_INPUT" = "plugins" ]; then
    # 显式指定 plugins 目录
    BUILD_MODE="batch"
    TARGET_DIR="$PLUGIN_INPUT"
    echo "📦 批量构建模式: 将构建 $TARGET_DIR 下所有插件"
else
    # 单个插件构建模式
    BUILD_MODE="single"
    
    # 确定插件目录
    # 如果输入不包含路径分隔符，则在 ./plugins/ 下查找
    if [[ "$PLUGIN_INPUT" != *"/"* ]]; then
        PLUGIN_DIR="./plugins/$PLUGIN_INPUT"
        echo "ℹ️  使用默认插件基础目录: $PLUGIN_DIR"
    else
        PLUGIN_DIR="$PLUGIN_INPUT"
    fi
    
    # 验证插件目录
    if [ ! -d "$PLUGIN_DIR" ]; then
        echo "❌ 错误: 插件目录不存在: $PLUGIN_DIR"
        if [[ "$PLUGIN_INPUT" != *"/"* ]]; then
            echo "💡 提示: 您可以尝试使用完整路径，或确保插件在 ./plugins/ 目录下"
        fi
        exit 1
    fi
    
    if [ ! -f "$PLUGIN_DIR/plugin.go" ] || [ ! -f "$PLUGIN_DIR/go.mod" ]; then
        echo "❌ 错误: $PLUGIN_DIR 不是有效的插件目录 (缺少 plugin.go 或 go.mod)"
        exit 1
    fi
fi

# 检查 Docker 是否可用
if ! command -v docker &> /dev/null; then
    echo "❌ 错误: Docker 未安装或不可用"
    echo "ℹ️  请安装 Docker 以支持 musl 构建"
    exit 1
fi

# 检查 Docker 服务是否运行
if ! docker info &> /dev/null; then
    echo "❌ 错误: Docker 服务未运行"
    echo "ℹ️  请启动 Docker 服务"
    exit 1
fi

# 调试模式下的 ldflags 设置
if [ "$DEBUG_MODE" -eq 1 ]; then
    LDFLAGS=""
    echo "🔍 调试模式已启用"
fi

# 显示构建信息
echo "🎯 目标平台: $MUSL_PLATFORMS"
echo "🏗️  LDFLAGS: $LDFLAGS"

# 构建单个插件的函数
build_single_plugin() {
    local plugin_dir="$1"
    local plugin_name=$(basename "$plugin_dir")
    
    echo ""
    echo "🔨 构建插件: $plugin_name ($plugin_dir)"
    
    # 验证插件目录
    if [ ! -d "$plugin_dir" ]; then
        echo "❌ 插件目录不存在: $plugin_dir"
        return 1
    fi
    
    if [ ! -f "$plugin_dir/plugin.go" ] || [ ! -f "$plugin_dir/go.mod" ]; then
        echo "❌ $plugin_dir 不是有效的插件目录 (缺少 plugin.go 或 go.mod)"
        return 1
    fi
    
    # 解析平台列表并构建
    IFS=',' read -ra PLATFORM_ARRAY <<< "$MUSL_PLATFORMS"
    local build_success=0
    local build_total=0
    
    for platform in "${PLATFORM_ARRAY[@]}"; do
        # 去除前后空格
        platform=$(echo "$platform" | xargs)
        
        # 解析 GOOS 和 GOARCH
        IFS='/' read -ra PARTS <<< "$platform"
        if [ ${#PARTS[@]} -ne 2 ]; then
            echo "⚠️  跳过无效平台格式: $platform"
            continue
        fi
        
        GOOS="${PARTS[0]}"
        GOARCH="${PARTS[1]}"
        OUTPUT_FILE="plugin-${GOOS}-${GOARCH}-musl.so"
        
        echo "  🏗️  构建 $GOOS/$GOARCH (musl) [平台: $docker_platform]..."
        build_total=$((build_total + 1))
        
        # 确定 Docker 平台
        local docker_platform="linux/amd64"  # 默认平台
        if [ "$GOARCH" = "arm64" ]; then
            docker_platform="linux/arm64"
        elif [ "$GOARCH" = "amd64" ]; then
            docker_platform="linux/amd64"
        fi
        
        # 使用 Docker 构建（支持多架构）
        if docker run --rm \
            --platform="$docker_platform" \
            -v "$PWD/$plugin_dir:/workspace" \
            -w /workspace \
            -e GOOS="$GOOS" \
            -e GOARCH="$GOARCH" \
            -e CGO_ENABLED=1 \
            golang:1.25.1-alpine \
            sh -c "apk add --no-cache gcc musl-dev && go build -buildmode=plugin -ldflags='$LDFLAGS' -o '$OUTPUT_FILE' ." 2>/dev/null; then
            
            # 检查文件是否成功生成
            if [ -f "$plugin_dir/$OUTPUT_FILE" ]; then
                echo "    ✅ 生成: $OUTPUT_FILE"
                build_success=$((build_success + 1))
                
                # 如果不是调试模式，使用 upx 压缩
                if [ "$DEBUG_MODE" -eq 0 ] && command -v upx &> /dev/null; then
                    upx -q "$plugin_dir/$OUTPUT_FILE" 2>/dev/null || true
                fi
            else
                echo "    ❌ 构建失败: 输出文件未生成"
            fi
        else
            echo "    ❌ 构建失败: $GOOS/$GOARCH"
        fi
    done
    
    echo "  📊 $plugin_name 构建结果: $build_success/$build_total 成功"
    return 0
}

# 批量构建插件的函数
build_all_plugins() {
    local plugins_dir="$1"
    
    if [ ! -d "$plugins_dir" ]; then
        echo "❌ 错误: 插件目录不存在: $plugins_dir"
        exit 1
    fi
    
    echo "📂 扫描插件目录: $plugins_dir"
    
    local total_plugins=0
    local built_plugins=0
    
    # 扫描插件目录
    for plugin_path in "$plugins_dir"/*/; do
        if [ -d "$plugin_path" ] && [ -f "${plugin_path}plugin.go" ]; then
            plugin_name=$(basename "$plugin_path")
            echo "  📦 发现插件: $plugin_name"
            total_plugins=$((total_plugins + 1))
        fi
    done
    
    if [ "$total_plugins" -eq 0 ]; then
        echo "❌ 在 $plugins_dir 中没有找到有效的插件目录"
        exit 1
    fi
    
    echo "📊 找到 $total_plugins 个插件，开始构建..."
    
    # 构建每个插件
    for plugin_path in "$plugins_dir"/*/; do
        if [ -d "$plugin_path" ] && [ -f "${plugin_path}plugin.go" ]; then
            if build_single_plugin "${plugin_path%/}"; then
                built_plugins=$((built_plugins + 1))
            fi
        fi
    done
    
    echo ""
    echo "✨ 批量构建完成!"
    echo "📊 总体构建摘要:"
    echo "  🎯 目标平台: $MUSL_PLATFORMS"
    echo "  📦 插件总数: $total_plugins"
    echo "  ✅ 成功构建: $built_plugins"
    
    # 显示生成的文件
    echo ""
    echo "📁 生成的文件:"
    find "$plugins_dir" -name "*-musl.so" -type f | sort
    
    if [ "$built_plugins" -eq 0 ]; then
        echo ""
        echo "❌ 所有插件构建都失败了"
        exit 1
    elif [ "$built_plugins" -lt "$total_plugins" ]; then
        echo ""
        echo "⚠️  部分插件构建失败"
        exit 0
    else
        echo ""
        echo "🎉 所有插件构建都成功了!"
        exit 0
    fi
}

# 检查多架构支持
echo "📦 检查 Docker 多架构支持..."
if command -v docker &> /dev/null; then
    if docker buildx version &> /dev/null; then
        echo "  ✅ Docker Buildx 可用"
    else
        echo "  ⚠️  Docker Buildx 不可用，可能影响 ARM 构建"
    fi
    
    # 检查 QEMU 支持
    if docker run --rm --privileged --platform linux/arm64 alpine:latest uname -m 2>/dev/null | grep -q "aarch64"; then
        echo "  ✅ ARM64 模拟可用"
    else
        echo "  ⚠️  ARM64 模拟不可用，可能需要设置 QEMU"
    fi
fi

# 拉取 Alpine Linux 镜像（支持多架构）
echo "📦 检查 Alpine Linux 镜像..."
if ! docker pull golang:1.25.1-alpine 2>/dev/null; then
    echo "❌ 错误: 无法拉取 golang:1.25.1-alpine 镜像"
    exit 1
fi

# 根据构建模式执行相应逻辑
if [ "$BUILD_MODE" = "batch" ]; then
    build_all_plugins "$TARGET_DIR"
else
    build_single_plugin "$PLUGIN_DIR"
    
    # 单个插件构建的摘要
    echo ""
    echo "✨ 插件构建完成!"
    echo "📊 构建摘要:"
    echo "  🎯 目标平台: $MUSL_PLATFORMS"
    
    # 显示生成的文件
    echo ""
    echo "📁 生成的文件:"
    find "$PLUGIN_DIR" -name "*-musl.so" -type f | sort
fi
