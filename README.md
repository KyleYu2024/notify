##  仅对原项目jianxcao/notify的主题部分进行修改，感谢大佬的开源项目

### 使用 Docker Compose（推荐）

1. **创建 docker-compose.yml 文件**：
```yaml
services:
  notify:
    image: kyleyu2024/notify:latest
    container_name: notify
    ports:
      - "7879:7879"
    volumes:
      - ./config:/config
    environment:
      - TZ=Asia/Shanghai
      - NOTIFY_USERNAME=admin
      - NOTIFY_PASSWORD=password
      - CONFIG_FILE=/config/config.yaml
      - LOG_LEVEL=info
      - LOG_FORMAT=text
    restart: unless-stopped
```

2.其他使用说明请移步原项目
