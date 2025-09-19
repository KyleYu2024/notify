# 前端构建阶段
FROM node:22-alpine AS frontend-builder

# 设置工作目录
WORKDIR /app

# 安装pnpm
RUN npm install -g pnpm

# 复制前端代码
COPY frontend/package.json frontend/pnpm-lock.yaml ./

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制前端源代码
COPY frontend/ .

# 构建前端应用
RUN pnpm build

# 后端构建阶段
FROM golang:1.25.1-alpine AS backend-builder

# 设置工作目录
WORKDIR /app

# 安装git、ca-certificates 以及构建所需工具链（gcc、musl-dev 等）
RUN apk add --no-cache git ca-certificates build-base pkgconfig

# 复制后端代码
COPY backend/ ./

# 下载依赖
RUN go mod download

# 编译应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o notify cmd/notify/main.go

# 构建插件（排除 demo 目录）
# RUN set -eux; \
#     for d in plugins/*; do \
#         if [ -d "$d" ] && [ -f "$d/plugin.go" ] && [ "$(basename "$d")" != "demo" ]; then \
#             echo "==> 构建插件: $(basename "$d")"; \
#             (cd "$d" && rm -f ./*.so && CGO_ENABLED=1 GOOS=linux go build -buildmode=plugin -o ./plugin.so ./plugin.go); \
#         fi; \
#     done

# 仅收集需要的文件到 plugins-dist（每个插件仅保留 setting.json 与 .so，排除 demo）
# RUN set -eux; \
#     rm -rf plugins-dist; mkdir -p plugins-dist; \
#     for d in plugins/*; do \
#         name=$(basename "$d"); \
#         if [ -d "$d" ] && [ "$name" != "demo" ]; then \
#             [ -f "$d/setting.json" ] || continue; \
#             mkdir -p "plugins-dist/$name"; \
#             cp "$d/setting.json" "plugins-dist/$name/"; \
#             if ls "$d"/*.so >/dev/null 2>&1; then cp "$d"/*.so "plugins-dist/$name/"; fi; \
#         fi; \
#     done

# 运行阶段
FROM alpine:latest

# 安装ca-certificates、timezone数据、wget (用于健康检查)、gosu和shadow-utils
RUN apk --no-cache add ca-certificates tzdata wget gosu shadow
RUN apk add --no-cache dumb-init
# 设置时区
ENV TZ=Asia/Shanghai

# 创建notify用户和组
RUN addgroup -g 1000 notify && \
    adduser -D -u 1000 -G notify notify

RUN mkdir -p /app
RUN mkdir -p /config
RUN mkdir -p /app/static
# 复制配置文件模板
RUN touch config/config.yaml
# 设置工作目录
WORKDIR /app
VOLUME /config

# 从后端构建阶段复制二进制文件
COPY --from=backend-builder /app/notify .

# 复制精简后的插件目录到镜像种子目录（仅 setting.json 和 .so）
# COPY --from=backend-builder /app/plugins-dist /opt/plugins-dist

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /app/dist /app/static
COPY ./entrypoint /app/entrypoint
RUN chmod +x /app/entrypoint
RUN chmod +x /app/notify

# 设置目录权限
RUN chown -R notify:notify /app /config
RUN chmod -R 755 /app

ENV PGID=1000
ENV PUID=1000
ENV UMASK=022
ENV CONFIG_FILE=/config/config.yaml
ENV PORT=8088
ENV LOG_LEVEL=info
ENV LOG_FORMAT=text
ENV STATIC_DIR=/app/static
ENV NOTIFY_USERNAME=
ENV NOTIFY_PASSWORD=
ENV PLUGINS_DIR=/config/plugins

# 启动应用
ENTRYPOINT ["/app/entrypoint" ]
EXPOSE 8088