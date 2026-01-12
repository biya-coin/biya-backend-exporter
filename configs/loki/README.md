# Loki 日志聚合配置

Loki 是 Grafana 开源的日志聚合系统，用于收集、存储和查询多个 Injective 节点的日志。

## 架构说明

```
多个 Injective 节点
  ├── Node 1 (Promtail) → 收集日志 → Loki
  ├── Node 2 (Promtail) → 收集日志 → Loki
  └── Node 3 (Promtail) → 收集日志 → Loki

Loki (日志聚合存储)
  ├── 标签索引
  ├── LogQL 查询引擎
  ├── REST API (HTTP)
  └── WebSocket API (实时日志流)
```

## 配置文件说明

### loki-config.yaml

Loki 主配置文件，包含：

- **server**: HTTP 和 gRPC 监听端口
- **storage**: 存储配置（文件系统）
- **schema_config**: 存储格式配置（TSDB）
- **limits_config**: 查询和写入限制
- **compactor**: 压缩和保留策略

## 使用方法

### 1. 启动 Loki 服务

```bash
# 使用 Docker Compose 启动
docker-compose -f compose.log.yaml up -d loki

# 验证服务
curl http://localhost:3100/ready
```

### 2. API 端点

- **健康检查**: `GET /ready`
- **查询范围**: `GET /loki/api/v1/query_range?query={job="injective-node"}&start=...&end=...`
- **实时日志（WebSocket）**: `ws://localhost:3100/loki/api/v1/tail?query={job="injective-node"}&limit=100`

### 3. LogQL 查询示例

```logql
# 查询所有节点日志
{job="injective-node"}

# 查询特定节点日志
{job="injective-node", node_id="node-1"}

# 查询错误日志
{job="injective-node", level="error"}

# 查询包含关键词的日志
{job="injective-node"} |= "error"

# 统计错误数量
sum(count_over_time({job="injective-node", level="error"}[1h]))
```

## 与 Grafana 集成

Grafana 已经配置了 Loki 数据源，可以直接在 Grafana Explore 中查询日志：

1. 打开 Grafana: http://localhost:3000
2. 进入 Explore
3. 选择 Loki 数据源
4. 输入 LogQL 查询
5. 点击 "Live" 按钮启用实时日志流

## 存储和保留

- 数据存储在 Docker volume: `loki-data`
- 默认保留时间: 7 天（168h）
- 可以通过 `limits_config.reject_old_samples_max_age` 修改

## 性能优化

- 查询限制：`max_query_parallelism: 32`
- 写入限制：`ingestion_rate_mb: 16`
- 缓存：`embedded_cache: max_size_mb: 100`

## 注意事项

1. 生产环境建议启用认证（`auth_enabled: true`）
2. 可以根据实际需求调整 `limits_config` 中的限制
3. 如果需要多租户支持，需要配置 `auth_enabled` 和相关认证
4. 数据存储在 Docker volume 中，需要定期备份
