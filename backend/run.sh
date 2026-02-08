# export GOOS=linux
# export GOARCH=amd64
export NOTIFY_USERNAME="admin"
export NOTIFY_PASSWORD="admin"
export LOG_LEVEL="debug"
export LOG_FORMAT="text"
export CGO_ENABLED=1
export PORT="7879"
export CONFIG_FILE=".config/config.yaml"
export PLUGINS_DIR="./plugins"

go run cmd/notify/main.go

