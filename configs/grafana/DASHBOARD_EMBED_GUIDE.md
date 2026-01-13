# Grafana Dashboard 完整嵌入指南

## 概述

Grafana 支持两种嵌入方式：
1. **整个 Dashboard 嵌入** - 显示完整的 Dashboard，包含所有面板
2. **单个面板嵌入** - 只显示单个面板（使用 `d-solo` 路径）

## 整个 Dashboard 嵌入

### URL 格式

```
http://45.249.245.183:3000/d/{dashboard-uid}/{dashboard-title}?orgId=1&from=now-24h&to=now&refresh=10s&kiosk=tv
```

### 日志分析 Dashboard 嵌入 URL

```
http://45.249.245.183:3000/d/log-analysis/log-analysis?orgId=1&from=now-24h&to=now&refresh=10s&kiosk=tv&theme=light
```

### HTML 代码示例

```html
<!DOCTYPE html>
<html>
<head>
    <title>日志分析 Dashboard</title>
    <style>
        body {
            margin: 0;
            padding: 0;
        }
        .dashboard-container {
            width: 100%;
            height: 100vh;
        }
    </style>
</head>
<body>
    <div class="dashboard-container">
        <iframe 
            src="http://45.249.245.183:3000/d/log-analysis/log-analysis?orgId=1&from=now-24h&to=now&refresh=10s&kiosk=tv" 
            width="100%" 
            height="100%" 
            frameborder="0"
            allowfullscreen>
        </iframe>
    </div>
</body>
</html>
```

## URL 参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| `orgId` | 组织 ID（默认 1） | `orgId=1` |
| `from` | 开始时间 | `from=now-24h`（最近24小时）<br>`from=now-7d`（最近7天）<br>`from=2024-01-01T00:00:00Z`（绝对时间） |
| `to` | 结束时间 | `to=now`（当前时间）<br>`to=2024-01-02T00:00:00Z`（绝对时间） |
| `refresh` | 自动刷新间隔 | `refresh=10s`（10秒）<br>`refresh=1m`（1分钟）<br>`refresh=5m`（5分钟） |
| `kiosk` | 显示模式 | `kiosk=tv`（隐藏导航栏，推荐用于嵌入）<br>`kiosk`（隐藏导航栏和工具栏） |
| `theme` | 主题模式 | `theme=light`（浅色模式）<br>`theme=dark`（深色模式，默认） |

## 时间范围选项

### 相对时间（推荐）

- `from=now-1h&to=now` - 最近1小时
- `from=now-6h&to=now` - 最近6小时
- `from=now-24h&to=now` - 最近24小时
- `from=now-7d&to=now` - 最近7天
- `from=now-30d&to=now` - 最近30天

### 绝对时间

- `from=2024-01-01T00:00:00Z&to=2024-01-02T00:00:00Z` - 指定日期范围

## 刷新间隔选项

- `refresh=5s` - 5秒
- `refresh=10s` - 10秒
- `refresh=30s` - 30秒
- `refresh=1m` - 1分钟
- `refresh=5m` - 5分钟
- `refresh=15m` - 15分钟
- `refresh=30m` - 30分钟

## 完整示例

### 示例 1：全屏嵌入（推荐）

```html
<iframe 
    src="http://45.249.245.183:3000/d/log-analysis/log-analysis?orgId=1&from=now-24h&to=now&refresh=10s&kiosk=tv&theme=light" 
    width="100%" 
    height="100vh" 
    frameborder="0"
    style="border: none;"
    allowfullscreen>
</iframe>
```

### 示例 2：固定高度嵌入

```html
<div style="width: 100%; height: 800px;">
    <iframe 
        src="http://45.249.245.183:3000/d/log-analysis/log-analysis?orgId=1&from=now-24h&to=now&refresh=10s&kiosk=tv&theme=light" 
        width="100%" 
        height="100%" 
        frameborder="0"
        allowfullscreen>
    </iframe>
</div>
```

### 示例 3：响应式嵌入

```html
<div style="position: relative; width: 100%; padding-bottom: 56.25%; height: 0; overflow: hidden;">
    <iframe 
        src="http://45.249.245.183:3000/d/log-analysis/log-analysis?orgId=1&from=now-24h&to=now&refresh=10s&kiosk=tv&theme=light" 
        style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; border: none;"
        allowfullscreen>
    </iframe>
</div>
```

## 单个面板嵌入（可选）

如果需要只嵌入单个面板，可以使用 `d-solo` 路径：

```
http://45.249.245.183:3000/d-solo/log-analysis/log-analysis?orgId=1&from=now-24h&to=now&refresh=10s&panelId=1&__feature.dashboardSceneSolo
```

**注意：** 需要先获取面板的实际 ID（在 Grafana 中点击面板右上角菜单 → "Share" → "Embed"）

## 配置要求

确保 Grafana 已配置允许 iframe 嵌入：

1. **环境变量**（在 `compose.yaml` 中）：
   ```yaml
   environment:
     - GF_SECURITY_ALLOW_EMBEDDING=true
   ```

2. **配置文件**（在 `grafana.ini` 中）：
   ```ini
   [security]
   allow_embedding = true
   ```

3. **重启 Grafana 服务**：
   ```bash
   docker compose restart grafana
   ```

## 测试

1. 打开测试页面：`configs/grafana/test-embed.html`
2. 或直接在浏览器中访问 Dashboard URL 验证是否正常显示
3. 检查浏览器控制台是否有错误信息

## 常见问题

### Q: iframe 显示空白或无法加载？

A: 检查以下几点：
1. Grafana 服务是否正常运行
2. `allow_embedding` 配置是否生效（需要重启 Grafana）
3. 网络连接是否正常
4. 浏览器控制台是否有错误信息

### Q: 如何隐藏 Grafana 的导航栏？

A: 在 URL 中添加 `kiosk=tv` 参数

### Q: 如何设置自动刷新？

A: 在 URL 中添加 `refresh=10s` 参数（10秒刷新一次）

### Q: 如何设置时间范围？

A: 使用 `from` 和 `to` 参数，例如：
- `from=now-24h&to=now` - 最近24小时
- `from=now-7d&to=now` - 最近7天

## 相关文档

- [Grafana iframe 嵌入指南](./IFRAME_EMBED_GUIDE.md)
- [日志分析嵌入指南](./LOG_ANALYSIS_EMBED_GUIDE.md)
- [Grafana 配置说明](./README.md)