# Loki + Promtail 日志收集快速开始指南

本指南介绍如何使用 Loki + Promtail 收集多个 Injective 节点的日志。

## 架构概述

```
多个 Injective 节点
  ├── Node 1 (Promtail) → 收集日志 → Loki
  ├── Node 2 (Promtail) → 收集日志 → Loki
  └── Node 3 (Promtail) → 收集日志 → Loki

Loki (日志聚合存储)
  ├── REST API (查询日志)
  └── WebSocket API (实时日志流)

前端应用
  ├── HTTP API (查询日志)
  └── WebSocket (实时日志流)
```

## 快速开始

### 1. 启动 Loki 服务

```bash
# 启动 Loki
docker-compose -f compose.log.yaml up -d loki

# 验证服务
curl http://localhost:3100/ready

# 查看日志
docker logs loki
```

### 2. 配置 Promtail（每个节点）

#### 方式1：使用 Docker Compose（同一网络）

如果所有节点在同一个 Docker 网络中：

```bash
# 编辑配置文件
vi configs/promtail/promtail-config.yaml

# 修改以下内容：
# - node_id: node-1 (改为对应节点的唯一标识)
# - __path__: /var/log/injective/*.log (改为实际日志路径)

# 启动 Promtail（示例）
docker-compose -f compose.log.yaml up -d promtail-example
```

#### 方式2：独立部署（不同服务器）

如果节点在不同的服务器上，在每个节点上：

1. 复制配置文件到节点：
   ```bash
   scp configs/promtail/promtail-config.yaml user@node-ip:/path/to/promtail/
   ```

2. 修改配置文件（`node_id` 和日志路径）

3. 使用 Docker 或 systemd 运行 Promtail（参考 `configs/promtail/README.md`）

### 3. 验证日志收集

```bash
# 查询所有日志
curl "http://localhost:3100/loki/api/v1/query_range?query={job=\"injective-node\"}&start=$(date -d '1 hour ago' +%s)000000000&end=$(date +%s)000000000&limit=100"

# 查询特定节点日志
curl "http://localhost:3100/loki/api/v1/query_range?query={job=\"injective-node\",node_id=\"node-1\"}&start=$(date -d '1 hour ago' +%s)000000000&end=$(date +%s)000000000&limit=100"

# 查询错误日志
curl "http://localhost:3100/loki/api/v1/query_range?query={job=\"injective-node\",level=\"error\"}&start=$(date -d '1 hour ago' +%s)000000000&end=$(date +%s)000000000&limit=100"
```

### 4. 使用 Grafana 查看日志

1. 打开 Grafana: http://localhost:3000
2. 进入 **Explore**
3. 选择 **Loki** 数据源
4. 输入 LogQL 查询：
   ```
   {job="injective-node"}
   ```
5. 点击 **Run query** 查看日志
6. 点击 **Live** 按钮启用实时日志流

### 5. 前端集成（WebSocket 实时日志）

#### JavaScript 示例

```javascript
// 连接 Loki WebSocket API
function connectLokiLiveTail(query, limit = 100) {
  const ws = new WebSocket(
    `ws://localhost:3100/loki/api/v1/tail?query=${encodeURIComponent(query)}&limit=${limit}`
  );

  ws.onopen = () => {
    console.log('WebSocket connected to Loki');
  };

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    
    // 处理日志流
    data.streams?.forEach(stream => {
      const labels = stream.stream;
      const entries = stream.values;
      
      entries.forEach(([timestamp, logLine]) => {
        // 处理单条日志
        handleLogEntry({
          timestamp: parseInt(timestamp) / 1000000, // 转换为毫秒
          labels: labels,
          content: logLine
        });
      });
    });
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  ws.onclose = () => {
    console.log('WebSocket closed');
  };

  return ws;
}

// 使用示例
// 查询所有节点日志
const ws = connectLokiLiveTail('{job="injective-node"}');

// 查询特定节点日志
const wsNode1 = connectLokiLiveTail('{job="injective-node", node_id="node-1"}');

// 查询错误日志
const wsError = connectLokiLiveTail('{job="injective-node", level="error"}');

// 查询包含关键词的日志
const wsKeyword = connectLokiLiveTail('{job="injective-node"} |= "error"');
```

## LogQL 查询示例

### 基础查询

```logql
# 查询所有节点日志
{job="injective-node"}

# 查询特定节点日志
{job="injective-node", node_id="node-1"}

# 查询错误日志
{job="injective-node", level="error"}

# 查询包含关键词的日志
{job="injective-node"} |= "error"

# 排除特定关键词
{job="injective-node"} != "debug"
```

### 统计查询

```logql
# 统计错误数量（1小时内）
sum(count_over_time({job="injective-node", level="error"}[1h]))

# 按级别统计日志数量
sum by (level) (count_over_time({job="injective-node"}[1h]))

# 按节点统计日志数量
sum by (node_id) (count_over_time({job="injective-node"}[1h]))

# 错误率（1小时内）
sum(rate({job="injective-node", level="error"}[1h])) / sum(rate({job="injective-node"}[1h])) * 100
```

### 时间范围查询

```logql
# 查询最近1小时的日志
{job="injective-node"} [1h]

# 查询最近24小时的日志
{job="injective-node"} [24h]

# 使用时间函数
{job="injective-node"} | json | timestamp > now() - 1h
```

## API 端点

### HTTP API

- **查询范围**: `GET /loki/api/v1/query_range?query={job="injective-node"}&start=...&end=...&limit=100`
- **即时查询**: `GET /loki/api/v1/query?query={job="injective-node"}&limit=100`
- **标签查询**: `GET /loki/api/v1/labels`
- **标签值查询**: `GET /loki/api/v1/label/<name>/values`
- **健康检查**: `GET /ready`

### WebSocket API

- **实时日志流**: `ws://localhost:3100/loki/api/v1/tail?query={job="injective-node"}&limit=100`

参数：
- `query`: LogQL 查询表达式（必需）
- `delay_for`: 延迟秒数（0-5，默认0）
- `limit`: 返回的最大条目数（默认100）
- `start`: 查询开始时间（纳秒 Unix 时间戳，默认1小时前）

## 常见问题

### 1. Promtail 无法连接 Loki

**问题**: Promtail 报错无法连接到 Loki

**解决方案**:
- 检查网络连接：`curl http://loki:3100/ready`
- 如果节点不在同一网络，修改 Promtail 配置中的 `clients.url` 为 Loki 的公网地址
- 检查防火墙设置

### 2. 日志没有收集

**问题**: 查询不到日志

**解决方案**:
- 检查日志路径是否正确
- 检查文件权限
- 查看 Promtail 日志：`docker logs promtail-example`
- 检查 `positions.yaml` 文件

### 3. WebSocket 连接失败

**问题**: 前端 WebSocket 连接失败

**解决方案**:
- 检查 Loki 是否正常运行：`curl http://localhost:3100/ready`
- 检查 WebSocket URL 是否正确
- 如果通过代理，确保代理支持 WebSocket（nginx 需要配置 `proxy_set_header Upgrade $http_upgrade;`）
- 检查 CORS 设置（如果需要跨域）

### 4. 性能问题

**问题**: 查询速度慢或内存占用高

**解决方案**:
- 限制查询时间范围
- 使用标签过滤减少数据量
- 调整 `limits_config` 中的限制
- 使用统计查询代替原始日志查询

## 下一步

- 查看 `configs/loki/README.md` 了解 Loki 详细配置
- 查看 `configs/promtail/README.md` 了解 Promtail 部署详情
- 集成到 Grafana Dashboard 实现日志可视化
- 配置告警规则（使用 Loki Ruler）

## 参考资源

- [Loki 官方文档](https://grafana.com/docs/loki/latest/)
- [LogQL 查询语言](https://grafana.com/docs/loki/latest/logql/)
- [Promtail 配置文档](https://grafana.com/docs/loki/latest/clients/promtail/)
- [Loki HTTP API](https://grafana.com/docs/loki/latest/api/)
