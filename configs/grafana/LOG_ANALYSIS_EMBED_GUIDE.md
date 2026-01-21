# 日志分析页面 Grafana 嵌入指南

## 概述

您看到的日志分析页面（包含日志统计概览、过滤选项、错误日志列表等）可以通过 Grafana 嵌入到其他页面中。

## 前提条件

1. ✅ **Loki 数据源已配置**：已在 `configs/grafana/provisioning/datasources/loki.yml` 中配置
2. ✅ **Grafana 支持 iframe 嵌入**：已在 `compose.yaml` 中配置 `GF_SECURITY_ALLOW_EMBEDDING=true`
3. ✅ **Loki 服务运行中**：确保 `compose.log.yaml` 中的 Loki 服务已启动

## 步骤 1：在 Grafana 中创建日志分析 Dashboard

### 1.1 创建新的 Dashboard

1. 登录 Grafana: http://localhost:3000
2. 进入 **Dashboards** → **New Dashboard**
3. 点击 **Add visualization**

### 1.2 创建日志统计面板

#### 面板 1：错误日志统计（Stat 面板）

1. 选择 **Loki** 数据源
2. 查询语句：
   ```logql
   sum(count_over_time({job="injective-node",level=~"ERR|ERROR|FATAL|PANIC"}[24h]))
   ```
3. 可视化类型：选择 **Stat**
4. 标题：`错误日志`
5. 颜色：红色

#### 面板 2：警告日志统计（Stat 面板）

1. 查询语句：
   ```logql
   sum(count_over_time({job="injective-node",level=~"WARN|WARNING"}[24h]))
   ```
2. 可视化类型：**Stat**
3. 标题：`警告日志`
4. 颜色：黄色

#### 面板 3：信息日志统计（Stat 面板）

1. 查询语句：
   ```logql
   sum(count_over_time({job="injective-node",level=~"INF|INFO"}[24h]))
   ```
2. 可视化类型：**Stat**
3. 标题：`信息日志`
4. 颜色：蓝色

#### 面板 4：总日志数统计（Stat 面板）

1. 查询语句：
   ```logql
   sum(count_over_time({job="injective-node"}[24h]))
   ```
2. 可视化类型：**Stat**
3. 标题：`总日志数`
4. 颜色：绿色

#### 面板 5：错误日志列表（Logs 面板）

1. 查询语句：
   ```logql
   {job="injective-node",level=~"ERR|ERROR|FATAL|PANIC"}
   ```
2. 可视化类型：**Logs**
3. 标题：`关键错误日志`
4. 配置：
   - 显示时间戳：是
   - 显示标签：是
   - 每页显示：50 条

#### 面板 6：按节点过滤的日志（Logs 面板，带变量）

1. 创建 Dashboard 变量：
   - 变量名：`node_id`
   - 类型：**Query**
   - 数据源：**Loki**
   - 查询：`label_values({job="injective-node"}, node_id)`
   - 标签：`节点ID`
   - 多选：是
   - 包含全部选项：是

2. 创建日志查询面板：
   - 查询语句：
     ```logql
     {job="injective-node",node_id=~"$node_id"}
     ```
   - 可视化类型：**Logs**
   - 标题：`节点日志 - $node_id`

### 1.3 配置 Dashboard 变量（可选）

为了支持动态过滤，可以添加以下变量：

- **时间范围变量**：`time_range`
  - 类型：**Interval**
  - 选项：`1h`, `6h`, `24h`, `7d`
  - 默认值：`24h`

- **日志级别变量**：`log_level`
  - 类型：**Custom**
  - 选项：
    - `全部级别` → `.*`
    - `错误` → `ERR|ERROR|FATAL|PANIC`
    - `警告` → `WARN|WARNING`
    - `信息` → `INF|INFO`
  - 默认值：`全部级别`

## 步骤 2：获取嵌入 URL

### 2.1 完整 Dashboard 嵌入

1. 在 Grafana 中打开创建的 Dashboard
2. 点击右上角 **Share** 按钮
3. 选择 **Embed** 标签
4. 复制 iframe URL
5. 修改 URL：
   - 将 `localhost:3000` 替换为实际服务器地址（如 `prv.grafana.biya.io`）
   - 添加参数：`&kiosk=tv`（隐藏导航栏）
   - 添加时间范围：`&from=now-24h&to=now`
   - 添加刷新间隔：`&refresh=10s`

**示例 URL：**
```
https://prv.grafana.biya.io/d/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&kiosk=tv
```

### 2.2 单个面板嵌入（推荐）

对于日志分析页面，可以分别嵌入不同的面板：

#### 日志统计概览（4个 Stat 面板并排）

```html
<div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px;">
  <!-- 错误日志 -->
  <iframe 
    src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=1&__feature.dashboardSceneSolo" 
    width="100%" 
    height="200" 
    frameborder="0">
  </iframe>
  
  <!-- 警告日志 -->
  <iframe 
    src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=2&__feature.dashboardSceneSolo" 
    width="100%" 
    height="200" 
    frameborder="0">
  </iframe>
  
  <!-- 信息日志 -->
  <iframe 
    src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=3&__feature.dashboardSceneSolo" 
    width="100%" 
    height="200" 
    frameborder="0">
  </iframe>
  
  <!-- 总日志数 -->
  <iframe 
    src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=4&__feature.dashboardSceneSolo" 
    width="100%" 
    height="200" 
    frameborder="0">
  </iframe>
</div>
```

#### 错误日志列表

```html
<iframe 
  src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=5&__feature.dashboardSceneSolo" 
  width="100%" 
  height="600" 
  frameborder="0">
</iframe>
```

## 步骤 3：完整嵌入示例

### 示例 1：模态窗口样式嵌入

```html
<!DOCTYPE html>
<html>
<head>
    <title>节点日志分析</title>
    <style>
        body {
            margin: 0;
            padding: 20px;
            font-family: Arial, sans-serif;
        }
        .modal {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .modal-header {
            padding: 20px;
            border-bottom: 1px solid #eee;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .modal-title {
            font-size: 20px;
            font-weight: bold;
        }
        .modal-actions {
            display: flex;
            gap: 10px;
        }
        .btn {
            padding: 8px 16px;
            border: 1px solid #ddd;
            border-radius: 4px;
            background: white;
            cursor: pointer;
        }
        .btn:hover {
            background: #f5f5f5;
        }
        .modal-content {
            padding: 20px;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 15px;
            margin-bottom: 20px;
        }
        .stat-card {
            border: 1px solid #ddd;
            border-radius: 4px;
            overflow: hidden;
        }
        .stat-card iframe {
            border: none;
        }
        .logs-section {
            margin-top: 20px;
        }
        .section-title {
            font-size: 16px;
            font-weight: bold;
            margin-bottom: 10px;
        }
        .logs-container {
            border: 1px solid #ddd;
            border-radius: 4px;
            overflow: hidden;
        }
    </style>
</head>
<body>
    <div class="modal">
        <div class="modal-header">
            <div class="modal-title">节点日志分析 - validator-001</div>
            <div class="modal-actions">
                <button class="btn">刷新日志</button>
                <button class="btn">导出日志</button>
                <button class="btn">关闭</button>
            </div>
        </div>
        
        <div class="modal-content">
            <!-- 日志统计概览 -->
            <div class="stats-grid">
                <div class="stat-card">
                    <iframe 
                        src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=1&__feature.dashboardSceneSolo" 
                        width="100%" 
                        height="150" 
                        frameborder="0">
                    </iframe>
                </div>
                <div class="stat-card">
                    <iframe 
                        src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=2&__feature.dashboardSceneSolo" 
                        width="100%" 
                        height="150" 
                        frameborder="0">
                    </iframe>
                </div>
                <div class="stat-card">
                    <iframe 
                        src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=3&__feature.dashboardSceneSolo" 
                        width="100%" 
                        height="150" 
                        frameborder="0">
                    </iframe>
                </div>
                <div class="stat-card">
                    <iframe 
                        src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=4&__feature.dashboardSceneSolo" 
                        width="100%" 
                        height="150" 
                        frameborder="0">
                    </iframe>
                </div>
            </div>
            
            <!-- 关键错误日志 -->
            <div class="logs-section">
                <div class="section-title">● 关键错误日志</div>
                <div class="logs-container">
                    <iframe 
                        src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=5&__feature.dashboardSceneSolo" 
                        width="100%" 
                        height="500" 
                        frameborder="0">
                    </iframe>
                </div>
            </div>
        </div>
    </div>
</body>
</html>
```

### 示例 2：使用 Dashboard 变量进行过滤

```html
<!-- 按节点ID过滤的日志 -->
<iframe 
  src="https://prv.grafana.biya.io/d-solo/log-analysis/log-analysis-dashboard?orgId=1&from=now-24h&to=now&refresh=10s&panelId=6&var-node_id=validator-001&__feature.dashboardSceneSolo" 
  width="100%" 
  height="600" 
  frameborder="0">
</iframe>
```

## 常用 LogQL 查询示例

### 按节点ID查询
```logql
{job="injective-node",node_id="validator-001"}
```

### 按级别查询
```logql
{job="injective-node",level=~"ERR|ERROR|FATAL|PANIC"}
```

### 按关键词搜索
```logql
{job="injective-node"} |= "consensus"
```

### 组合查询
```logql
{job="injective-node",node_id="validator-001",level=~"ERR|ERROR"} |= "timeout"
```

### 时间范围统计
```logql
sum(count_over_time({job="injective-node",level=~"ERR|ERROR"}[24h]))
```

## 注意事项

1. **面板 ID 获取**：
   - 在 Grafana 中点击面板右上角菜单 → "Share" → "Embed"
   - 复制 URL 中的 `panelId` 参数值

2. **URL 参数说明**：
   - `d-solo`：Grafana 推荐的单个面板嵌入格式
   - `from=now-24h&to=now`：时间范围（最近24小时）
   - `refresh=10s`：自动刷新间隔
   - `panelId=X`：面板 ID
   - `var-node_id=xxx`：Dashboard 变量值

3. **性能优化**：
   - 使用 `d-solo` 格式比完整 Dashboard 嵌入性能更好
   - 合理设置刷新间隔，避免过于频繁的查询
   - 限制日志查询数量（在 Logs 面板中设置）

4. **安全性**：
   - 生产环境建议使用 API Key 认证
   - 配置 CORS 策略（如果需要跨域嵌入）

## 相关文档

- [Grafana iframe 嵌入指南](./IFRAME_EMBED_GUIDE.md)
- [Loki 配置文档](../loki/README.md)
- [LogQL 查询语言](https://grafana.com/docs/loki/latest/logql/)
