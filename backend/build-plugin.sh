
export LOG_LEVEL="debug"
export LOG_FORMAT="text"
export CGO_ENABLED=1
export PORT="7879"
export CONFIG_FILE=".config/config.yaml"

plugins_dir="./plugins"
fail=0

# 解析参数，支持 debug 模式
DEBUG_MODE=0
for arg in "$@"; do
    if [ "$arg" = "debug" ] || [ "$arg" = "--debug" ] || [ "$arg" = "-d" ]; then
        DEBUG_MODE=1
    fi
done

for plugin_path in "$plugins_dir"/*; do
    if [ -d "$plugin_path" ] && [ -f "$plugin_path/plugin.go" ]; then
        plugin_name=$(basename "$plugin_path")
        echo "==> 构建插件: $plugin_name"
        (
            cd "$plugin_path" && \
            rm -f ./*.so && \
            if [ "$DEBUG_MODE" -eq 1 ]; then \
                echo "使用 debug 模式构建" 
                go build -gcflags="all=-N -l" -buildmode=plugin -o ./plugin.so ./plugin.go; \
            else \
                echo "使用 release 模式构建"
                go build -buildmode=plugin -o ./plugin.so ./plugin.go; \
            fi
        )
        if [ $? -ne 0 ]; then
            echo "构建失败: $plugin_name"
            fail=1
        else
            echo "构建成功: $plugin_name"
        fi
    fi
done

if [ "$fail" -ne 0 ]; then
    echo "部分插件构建失败"
    exit 1
fi

echo "所有可用插件构建完成"
