#!/bin/bash

# 默认构建平台
DEFAULT_PLATFORMS="linux/amd64,linux/arm64,darwin/amd64,darwin/arm64"

# 帮助信息
show_help() {
    echo "Usage: $0 [options] [plugin_name]"
    echo ""
    echo "Options:"
    echo "  -p, --platforms    指定构建平台 (例如: linux/amd64,darwin/arm64)"
    echo "                     默认: $DEFAULT_PLATFORMS"
    echo "  -d, --debug        启用调试模式 (不压缩二进制文件)"
    echo "  -h, --help         显示帮助信息"
    echo ""
    echo "Examples:"
    echo "  $0                 # 构建所有插件"
    echo "  $0 emby            # 只构建 emby 插件"
    echo "  $0 -p linux/amd64 emby  # 自定义平台构建 emby 插件"
    echo "  $0 -d              # 调试模式构建所有插件"
    echo ""
    echo "说明:"
    echo "  - 使用 xgo 进行交叉编译，支持 Linux、macOS 等平台"
    echo "  - musl 版本通过 GitHub Actions 自动构建"
}

# 解析命令行参数
PLATFORMS="$DEFAULT_PLATFORMS"
DEBUG_MODE=0
PLUGIN_NAME=""

while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--platforms)
            PLATFORMS="$2"
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
            PLUGIN_NAME="$1"
            shift
            ;;
    esac
done

# 设置编译参数
LDFLAGS="-s -w"
if [ "$DEBUG_MODE" -eq 1 ]; then
    LDFLAGS=""
    echo "🔍 调试模式已启用"
fi

# 构建单个插件（使用 xgo）
build_plugin_impl() {
    local plugin_dir="$1"
    echo "🔨 构建插件: $plugin_dir"
    
    # 进入插件目录
    cd "$plugin_dir" || exit 1
    
    # 检查必要文件
    if [ ! -f "plugin.go" ] || [ ! -f "go.mod" ]; then
        echo "❌ 错误: plugin.go 或 go.mod 文件不存在"
        cd ..
        return 1
    fi
    
    echo "📦 目标平台: $PLATFORMS"
    
    # 使用 xgo 构建
    xgo --targets="$PLATFORMS" \
        --buildmode=plugin \
        --ldflags="$LDFLAGS" \
        --out="plugin" \
        .
    
    # 检查构建结果
    if [ $? -ne 0 ]; then
        echo "❌ 构建失败: $plugin_dir"
        cd ..
        return 1
    fi
    
    # 重命名输出文件
    for file in plugin-*; do
        if [[ -f "$file" ]]; then
            # 提取平台信息
            platform=${file#plugin-}
            
            # 移除各种可能的扩展名，只保留平台信息
            platform=${platform%.so}     # 移除 .so
            platform=${platform%.exe}    # 移除 .exe
            platform=${platform%.dll}    # 移除 .dll
            
            # 生成最终文件名
            final_name="plugin-${platform}.so"
            
            # 重命名（如果文件名不同的话）
            if [[ "$file" != "$final_name" ]]; then
                mv "$file" "$final_name"
            fi
            
            echo "✅ 生成: $final_name"
            
            # 如果不是调试模式，使用 upx 压缩（如果可用）
            if [ "$DEBUG_MODE" -eq 0 ] && command -v upx &> /dev/null; then
                upx -q "$final_name" || true
            fi
        fi
    done
    
    cd ..
    return 0
}


# 主构建逻辑
echo "🚀 开始构建..."
echo "🎯 目标平台: $PLATFORMS"

if [ -n "$PLUGIN_NAME" ]; then
    echo "📍 指定插件: $PLUGIN_NAME"
fi

echo "ℹ️  musl 版本将通过 GitHub Actions 自动构建"

# 确保在 plugins 目录下
if [[ $(basename "$PWD") != "plugins" ]]; then
    if [[ -d "plugins" ]]; then
        cd plugins || exit 1
    elif [[ -d "backend/plugins" ]]; then
        cd backend/plugins || exit 1
    else
        echo "❌ 错误: 无法找到 plugins 目录"
        exit 1
    fi
fi

# 构建插件函数
build_plugin() {
    local plugin_dir="$1"
    build_plugin_impl "$plugin_dir"
}

# 构建指定插件或所有插件
if [ -n "$PLUGIN_NAME" ]; then
    if [ -d "$PLUGIN_NAME" ]; then
        build_plugin "$PLUGIN_NAME"
    else
        echo "❌ 错误: 插件目录 '$PLUGIN_NAME' 不存在"
        exit 1
    fi
else
    # 构建所有插件
    for plugin_dir in */; do
        if [ -f "${plugin_dir}plugin.go" ]; then
            build_plugin "${plugin_dir%/}"
        fi
    done
fi

echo "✨ 构建完成!"
echo "📊 构建摘要:"
echo "  ✅ 目标平台: $PLATFORMS"
