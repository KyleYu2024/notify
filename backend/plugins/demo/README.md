# Demo Plugin

这是一个演示插件，展示了如何实现 Notify 系统的插件。

## 功能

- 从输入数据中提取标题、内容、图片、URL等信息
- 根据配置设置格式化消息
- 支持添加时间戳和调试信息
- 返回标准化的通知格式

## 配置选项

- **prefix**: 消息前缀，会在标题前添加
- **add_timestamp**: 是否添加时间戳
- **default_image**: 默认图片 URL
- **debug**: 是否启用调试模式

## 构建

```bash
cd backend
go build -buildmode=plugin -o ../plugins/demo/plugin.so ../plugins/demo/plugin.go
```

## 测试数据

```json
{
  "title": "测试通知",
  "content": "这是一条测试消息",
  "image": "https://example.com/image.jpg",
  "url": "https://example.com",
  "targets": ["user1", "user2"]
}
```