
export LOG_LEVEL="debug"
export LOG_FORMAT="text"
export CGO_ENABLED=1
export PORT="8088"
export CONFIG_FILE=".config/config.yaml"

plugins_dir="./plugins"
fail=0

for plugin_path in "$plugins_dir"/*; do
    if [ -d "$plugin_path" ] && [ -f "$plugin_path/plugin.go" ]; then
        plugin_name=$(basename "$plugin_path")
        echo "==> 构建插件: $plugin_name"
        (
            cd "$plugin_path" && \
            rm -f ./*.so && \
            go build -gcflags="all=-N -l" -buildmode=plugin -o ./plugin.so ./plugin.go
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
