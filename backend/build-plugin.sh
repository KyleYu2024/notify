
export LOG_LEVEL="debug"
export LOG_FORMAT="text"
export CGO_ENABLED=1
export PORT="8088"
export CONFIG_FILE=".config/config.yaml"

# Go版本配置 - 用于Docker编译环境
GO_VERSION="1.25.1"

plugins_dir="./plugins"
fail=0

# 支持的平台配置
get_platform_info() {
    case "$1" in
        "darwin-amd64")
            echo "darwin amd64"
            ;;
        "darwin-arm64")
            echo "darwin arm64"
            ;;
        "linux-amd64")
            echo "linux amd64"
            ;;
        "linux-arm64")
            echo "linux arm64"
            ;;
        *)
            echo ""
            ;;
    esac
}

# 所有支持的平台列表
ALL_PLATFORMS="darwin-amd64 darwin-arm64 linux-amd64 linux-arm64"

# 默认值
DEBUG_MODE=0
SPECIFIC_PLUGIN=""
TARGET_PLATFORMS="all"
SHOW_HELP=0

# 帮助信息
show_help() {
    echo "用法: $0 [选项] [插件名]"
    echo ""
    echo "选项:"
    echo "  -d, --debug           使用debug模式构建"
    echo "  -p, --platform PLAT   指定目标平台，支持："
    echo "                        - all (默认，构建所有平台)"
    echo "                        - darwin-amd64 (macOS Intel)"
    echo "                        - darwin-arm64 (macOS Apple Silicon)"
    echo "                        - linux-amd64 (Linux x86_64)"
    echo "                        - linux-arm64 (Linux ARM64)"
    echo "                        - 也可以用逗号分隔多个平台，如: darwin-amd64,linux-amd64"
    echo "  -h, --help            显示此帮助信息"
    echo ""
    echo "参数:"
    echo "  插件名                 指定要构建的插件名称（可选，不指定则构建所有插件）"
    echo ""
    echo "交叉编译要求:"
    echo "  由于插件需要CGO支持，交叉编译需要安装相应的工具链："
    echo "  Linux AMD64:  需要 x86_64-linux-gnu-gcc 或 Docker"
    echo "  Linux ARM64:  需要 aarch64-linux-gnu-gcc 或 Docker"
    echo "  Darwin目标:   只能在 macOS 系统上编译"
    echo ""
    echo "  macOS 安装交叉编译工具链："
    echo "    brew install x86_64-elf-gcc      # Linux AMD64"
    echo "    brew install aarch64-elf-gcc     # Linux ARM64"
    echo "  或者确保安装了 Docker 来使用容器编译"
    echo ""
    echo "配置说明:"
    echo "  可以通过修改脚本顶部的 GO_VERSION 变量来更改Docker镜像的Go版本"
    echo "  当前设置: golang:${GO_VERSION}-alpine"
    echo ""
    echo "示例:"
    echo "  $0                    # 构建所有插件的所有平台版本"
    echo "  $0 emby               # 构建emby插件的所有平台版本"
    echo "  $0 -p darwin-amd64    # 构建所有插件的macOS Intel版本"
    echo "  $0 -p linux-arm64 emby # 构建emby插件的Linux ARM64版本"
    echo "  $0 -d -p darwin-amd64,linux-amd64 emby # debug模式构建emby插件的指定平台版本"
}

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -d|--debug)
            DEBUG_MODE=1
            shift
            ;;
        -p|--platform)
            TARGET_PLATFORMS="$2"
            shift 2
            ;;
        -h|--help)
            SHOW_HELP=1
            shift
            ;;
        debug) # 兼容旧的debug参数
            DEBUG_MODE=1
            shift
            ;;
        -*)
            echo "未知选项: $1"
            SHOW_HELP=1
            shift
            ;;
        *)
            if [ -z "$SPECIFIC_PLUGIN" ]; then
                SPECIFIC_PLUGIN="$1"
            else
                echo "错误: 只能指定一个插件名"
                SHOW_HELP=1
            fi
            shift
            ;;
    esac
done

if [ "$SHOW_HELP" -eq 1 ]; then
    show_help
    exit 0
fi

# 解析目标平台
if [ "$TARGET_PLATFORMS" = "all" ]; then
    build_platforms_str="$ALL_PLATFORMS"
else
    build_platforms_str=$(echo "$TARGET_PLATFORMS" | tr ',' ' ')
    # 验证平台是否支持
    for platform in $build_platforms_str; do
        platform_info=$(get_platform_info "$platform")
        if [ -z "$platform_info" ]; then
            echo "错误: 不支持的平台 '$platform'"
            echo "支持的平台: $ALL_PLATFORMS"
            exit 1
        fi
    done
fi

echo "目标平台: $build_platforms_str"
if [ "$DEBUG_MODE" -eq 1 ]; then
    echo "构建模式: Debug"
else
    echo "构建模式: Release"
fi

# 构建函数
build_plugin() {
    local plugin_path="$1"
    local plugin_name="$2"
    local platform="$3"
    local goos="$4"
    local goarch="$5"
    
    echo "  --> 构建 $plugin_name for $platform ($goos/$goarch)"
    
    (
        cd "$plugin_path" && \
        # 在设置环境变量之前获取当前平台信息
        local current_goos=$(go env GOOS)
        local current_goarch=$(go env GOARCH)
        export GOOS="$goos" && \
        export GOARCH="$goarch" && \
        local output_name="plugin-${platform}.so"
        rm -f "./$output_name" && \
        
        # 检查是否为交叉编译
        if [ "$goos" != "$current_goos" ] || [ "$goarch" != "$current_goarch" ]; then
            echo "    检测到交叉编译 ($current_goos/$current_goarch -> $goos/$goarch)"
            
            # 根据目标平台设置编译器
            case "$goos-$goarch" in
                "linux-amd64")
                    if command -v x86_64-linux-gnu-gcc >/dev/null 2>&1; then
                        export CC=x86_64-linux-gnu-gcc
                        export CXX=x86_64-linux-gnu-g++
                        echo "    使用 x86_64-linux-gnu-gcc 编译器"
                    else
                        echo "    警告: 未找到 x86_64-linux-gnu-gcc，尝试使用 Docker 编译..."
                        # 使用Docker编译
                        if command -v docker >/dev/null 2>&1; then
                            docker run --rm \
                                -v "$(pwd)":/workspace \
                                -w /workspace \
                                golang:${GO_VERSION}-alpine \
                                sh -c "apk add --no-cache git ca-certificates build-base pkgconfig && \
                                       go build -buildmode=plugin -o $output_name ./plugin.go"
                            return $?
                        else
                            echo "    错误: 需要安装 Docker 或 x86_64-linux-gnu-gcc 来进行交叉编译"
                            return 1
                        fi
                    fi
                    ;;
                "linux-arm64")
                    if command -v aarch64-linux-gnu-gcc >/dev/null 2>&1; then
                        export CC=aarch64-linux-gnu-gcc
                        export CXX=aarch64-linux-gnu-g++
                        echo "    使用 aarch64-linux-gnu-gcc 编译器"
                    else
                        echo "    警告: 未找到 aarch64-linux-gnu-gcc，尝试使用 Docker 编译..."
                        if command -v docker >/dev/null 2>&1; then
                            docker run --rm \
                                -v "$(pwd)":/workspace \
                                -w /workspace \
                                --platform linux/arm64 \
                                golang:${GO_VERSION}-alpine \
                                sh -c " git ca-certificates build-base pkgconfig && \
                                       go build -buildmode=plugin -o $output_name ./plugin.go"
                            return $?
                        else
                            echo "    错误: 需要安装 Docker 或 aarch64-linux-gnu-gcc 来进行交叉编译"
                            return 1
                        fi
                    fi
                    ;;
                "darwin-arm64")
                    if [ "$current_goos" = "darwin" ]; then
                        echo "    在 macOS 上交叉编译 ARM64"
                    else
                        echo "    错误: 无法在非 macOS 系统上编译 Darwin 目标"
                        return 1
                    fi
                    ;;
                "darwin-amd64")
                    if [ "$current_goos" = "darwin" ]; then
                        echo "    在 macOS 上交叉编译 AMD64"
                    else
                        echo "    错误: 无法在非 macOS 系统上编译 Darwin 目标"
                        return 1
                    fi
                    ;;
                *)
                    echo "    警告: 未配置的平台组合 $goos/$goarch"
                    ;;
            esac
        else
            echo "    本地平台构建"
        fi
        
        if [ "$DEBUG_MODE" -eq 1 ]; then \
            go build -gcflags="all=-N -l" -buildmode=plugin -o "./$output_name" ./plugin.go; \
        else \
            go build -buildmode=plugin -o "./$output_name" ./plugin.go; \
        fi
    )
    
    local build_result=$?
    if [ $build_result -ne 0 ]; then
        echo "  ✗ 构建失败: $plugin_name for $platform"
        return 1
    else
        echo "  ✓ 构建成功: $plugin_name for $platform"
        return 0
    fi
}

# 主构建循环
for plugin_path in "$plugins_dir"/*; do
    if [ -d "$plugin_path" ] && [ -f "$plugin_path/plugin.go" ]; then
        plugin_name=$(basename "$plugin_path")
        
        # 如果指定了特定插件，检查是否匹配
        if [ -n "$SPECIFIC_PLUGIN" ] && [ "$plugin_name" != "$SPECIFIC_PLUGIN" ]; then
            continue
        fi
        
        echo "==> 构建插件: $plugin_name"
        
        # 为每个平台构建
        for platform in $build_platforms_str; do
            platform_info=$(get_platform_info "$platform")
            read -r goos goarch <<< "$platform_info"
            
            if ! build_plugin "$plugin_path" "$plugin_name" "$platform" "$goos" "$goarch"; then
                fail=1
            fi
        done
    fi
done

# 如果指定了特定插件但未找到
if [ -n "$SPECIFIC_PLUGIN" ]; then
    plugin_found=0
    for plugin_path in "$plugins_dir"/*; do
        if [ -d "$plugin_path" ] && [ -f "$plugin_path/plugin.go" ]; then
            plugin_name=$(basename "$plugin_path")
            if [ "$plugin_name" = "$SPECIFIC_PLUGIN" ]; then
                plugin_found=1
                break
            fi
        fi
    done
    
    if [ $plugin_found -eq 0 ]; then
        echo "错误: 未找到插件 '$SPECIFIC_PLUGIN'"
        echo "可用插件:"
        for plugin_path in "$plugins_dir"/*; do
            if [ -d "$plugin_path" ] && [ -f "$plugin_path/plugin.go" ]; then
                echo "  - $(basename "$plugin_path")"
            fi
        done
        exit 1
    fi
fi

if [ "$fail" -ne 0 ]; then
    echo "部分插件构建失败"
    exit 1
fi

echo "构建完成！"
