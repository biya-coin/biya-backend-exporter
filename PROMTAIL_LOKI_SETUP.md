# Promtail 节点日志收集与 Loki 集成文档

本文档详细说明如何使用 Promtail 从各个节点收集日志并发送到 Loki 进行集中存储和查询。

## 目录

- [架构概述](#架构概述)
- [组件说明](#组件说明)
- [配置文件](#配置文件)
- [部署步骤](#部署步骤)
- [验证与测试](#验证与测试)
- [日志查询](#日志查询)
- [故障排查](#故障排查)
- [性能优化](#性能优化)

## 架构概述

```
┌─────────────────────────────────────────────────────────────┐
│                    日志收集架构图                              │
└─────────────────────────────────────────────────────────────┘

多个 Injective 节点
  ├── Node 1
  │   ├── 日志文件: ~/.injectived/logs/inj.log
  │   └── Promtail (收集代理)
  │       └── 推送日志 → Loki
  │
  ├── Node 2
  │   ├── 日志文件: ~/.injectived/logs/inj.log
  │   └── Promtail (收集代理)
  │       └── 推送日志 → Loki
  │
  └── Node 3
      ├── 日志文件: ~/.injectived/logs/inj.log
      └── Promtail (收集代理)
          └── 推送日志 → Loki

                    ↓
            ┌───────────────┐
            │  Loki 服务器   │
            │  (日志聚合)    │
            └───────────────┘
                    ↓
        ┌───────────────────────┐
        │  查询与可视化          │
        ├── Grafana Dashboard  │
        ├── HTTP API            │
        └── WebSocket (实时流)  │
```

### 数据流向

1. **日志生成**: Injective 节点将日志写入本地文件（如 `~/.injectived/logs/inj.log`）
2. **日志收集**: Promtail 监控日志文件，实时读取新日志
3. **日志处理**: Promtail 解析日志格式，提取标签（node_id、level、module 等）
4. **日志推送**: Promtail 通过 HTTP API 将日志推送到 Loki
5. **日志存储**: Loki 接收并存储日志，建立索引
6. **日志查询**: 通过 Grafana 或 API 查询日志

## 组件说明

### Promtail

**作用**: 日志收集代理，部署在各个节点上

**功能**:
- 监控指定目录下的日志文件
- 实时读取新日志内容
- 解析日志格式，提取结构化信息
- 添加标签（node_id、level、module 等）
- 通过 HTTP API 推送日志到 Loki
- 支持断点续传（positions 文件）

**端口**: 9080 (HTTP API)

### Loki

**作用**: 日志聚合存储服务

**功能**:
- 接收来自多个 Promtail 的日志
- 存储日志数据（chunks 和索引）
- 提供查询 API（HTTP 和 WebSocket）
- 支持 LogQL 查询语言
- 数据保留策略管理

**端口**: 
- 3100 (HTTP API)
- 9096 (gRPC，内部通信)

## 配置文件

### Promtail 配置

配置文件位置: `configs/promtail/promtail-config.yaml`

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

# 记录读取位置，用于断点续传
positions:
  filename: /tmp/positions.yaml

# Loki 客户端配置
clients:
  - url: http://loki:3100/loki/api/v1/push
    # 如果 Loki 需要认证，可以配置：
    # bearer_token_file: /path/to/token
    # basic_auth:
    #   username: user
    #   password_file: /path/to/password

# 日志采集配置
scrape_configs:
  - job_name: injective-node
    static_configs:
      - targets:
          - localhost
        labels:
          job: injective-node
          # 重要：每个节点需要修改 node_id 为唯一标识
          node_id: local-node
          # 日志路径（当前机器上的日志文件）
          __path__: /home/injectived/logs/inj.log
    
    # 管道处理（解析日志格式）
    pipeline_stages:
      # 移除 ANSI 颜色码
      - regex:
          expression: '\x1b\[[0-9;]*m'
          replacement: ''
      
      # 解析日志级别（INF, ERR, WRN, DBG 等）
      - regex:
          expression: '(?P<timestamp>\d{1,2}:\d{2}(?:AM|PM)|\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})\s+(?P<level>INF|ERR|WRN|DBG|FATAL|PANIC)'
          source: log
      
      # 提取级别作为标签
      - match:
          selector: '{job="injective-node"}'
          stages:
            - regex:
                expression: '\s+(?P<level>INF|ERR|WRN|DBG|FATAL|PANIC)\s+'
            - labels:
                level:
      
      # 提取 module 信息（module=xxx）
      - match:
          selector: '{job="injective-node"}'
          stages:
            - regex:
                expression: 'module=(?P<module>\w+)'
            - labels:
                module:
      
      # 标记错误日志
      - match:
          selector: '{job="injective-node"}'
          stages:
            - regex:
                expression: '\s+(ERR|FATAL|PANIC)\s+'
            - labels:
                has_error: "true"
```

**关键配置项说明**:

1. **node_id**: 每个节点必须设置唯一标识，用于区分不同节点的日志
2. **__path__**: 日志文件的路径，支持通配符（如 `/var/log/injective/*.log`）
3. **clients.url**: Loki 服务器的地址
   - 同一 Docker 网络: `http://loki:3100/loki/api/v1/push`
   - 不同服务器: `http://loki-server-ip:3100/loki/api/v1/push`
4. **pipeline_stages**: 日志处理管道，用于解析和提取标签

### Loki 配置

配置文件位置: `configs/loki/loki-config.yaml`

```yaml
auth_enabled: false

server:
  http_listen_port: 3100
  grpc_listen_port: 9096

common:
  instance_addr: 127.0.0.1
  path_prefix: /loki
  storage:
    filesystem:
      chunks_directory: /loki/chunks
      rules_directory: /loki/rules
  replication_factor: 1
  ring:
    kvstore:
      store: inmemory

# Schema 配置（存储格式）
schema_config:
  configs:
    - from: 2020-10-24
      store: tsdb
      object_store: filesystem
      schema: v13
      index:
        prefix: index_
        period: 24h

# 限制配置
limits_config:
  # 查询限制
  max_query_length: 720h
  max_query_parallelism: 16
  max_query_series: 500
  
  # 写入限制
  ingestion_rate_mb: 16
  ingestion_burst_size_mb: 32
  max_streams_per_user: 0
  max_line_size: 256KB
  
  # 保留策略
  reject_old_samples: true
  reject_old_samples_max_age: 168h
  retention_period: 168h  # 7天保留期
```

**关键配置项说明**:

1. **retention_period**: 日志保留时间，默认 7 天
2. **ingestion_rate_mb**: 日志写入速率限制
3. **max_query_length**: 最大查询时间范围

### Docker Compose 配置

配置文件位置: `compose.log.yaml`

```yaml
services:
  # Loki - 日志聚合存储服务
  loki:
    image: grafana/loki:3.0.0
    container_name: loki
    command: -config.file=/etc/loki/local-config.yaml
    ports:
      - "3100:3100"
      - "9096:9096"
    volumes:
      - ./configs/loki/loki-config.yaml:/etc/loki/local-config.yaml:ro
      - loki-data:/loki
    networks:
      - biya-network
    restart: unless-stopped

  # Promtail - 日志收集代理
  promtail:
    image: grafana/promtail:3.0.0
    container_name: promtail
    command: -config.file=/etc/promtail/config.yaml
    volumes:
      - ./configs/promtail/promtail-config.yaml:/etc/promtail/config.yaml:ro
      - /home/ubuntu/.injectived/logs:/home/injectived/logs:ro
      - promtail-positions:/tmp
    networks:
      - biya-network
    restart: unless-stopped
    depends_on:
      - loki

networks:
  biya-network:
    external: true
    name: biya-backend

volumes:
  loki-data:
  promtail-positions:
```

## 部署步骤

### 方式1: 单节点部署（本地测试）

适用于在同一台机器上部署 Loki 和 Promtail，收集本地节点日志。

#### 1. 准备工作

```bash
# 检查日志文件是否存在
ls -la ~/.injectived/logs/inj.log

# 确保 Docker 网络存在
docker network ls | grep biya-backend
# 如果不存在，创建网络：
docker network create biya-backend
```

#### 2. 修改 Promtail 配置

编辑 `configs/promtail/promtail-config.yaml`:

```yaml
labels:
  node_id: local-node  # 修改为唯一标识
  __path__: /home/injectived/logs/inj.log  # 确认日志路径
```

#### 3. 启动服务

```bash
# 启动 Loki
docker-compose -f compose.log.yaml up -d loki

# 等待 Loki 就绪（约10秒）
curl http://localhost:3100/ready

# 启动 Promtail
docker-compose -f compose.log.yaml up -d promtail

# 验证 Promtail 就绪
curl http://localhost:9080/ready
```

#### 4. 验证部署

```bash
# 查看服务状态
docker-compose -f compose.log.yaml ps

# 查看服务日志
docker-compose -f compose.log.yaml logs -f promtail
docker-compose -f compose.log.yaml logs -f loki
```

### 方式2: 多节点部署

适用于在多个服务器上部署 Promtail，所有日志发送到中央 Loki 服务器。

#### 节点1: 部署 Loki（中央服务器）

```bash
# 1. 在中央服务器上启动 Loki
docker-compose -f compose.log.yaml up -d loki

# 2. 验证 Loki 可访问（确保防火墙开放 3100 端口）
curl http://loki-server-ip:3100/ready
```

#### 节点2-N: 部署 Promtail（各个节点）

在每个节点上：

```bash
# 1. 创建 Promtail 配置目录
mkdir -p /opt/promtail/config

# 2. 复制配置文件
scp configs/promtail/promtail-config.yaml user@node-ip:/opt/promtail/config/

# 3. 修改配置文件
vi /opt/promtail/config/promtail-config.yaml
```

修改以下内容：

```yaml
clients:
  - url: http://loki-server-ip:3100/loki/api/v1/push  # 修改为 Loki 服务器地址

labels:
  node_id: node-1  # 修改为节点唯一标识（node-1, node-2, ...）
  __path__: /home/injectived/logs/inj.log  # 确认日志路径
```

#### 使用 Docker 运行 Promtail

```bash
# 创建 docker-compose.yml
cat > /opt/promtail/docker-compose.yml <<EOF
version: '3'
services:
  promtail:
    image: grafana/promtail:3.0.0
    container_name: promtail
    command: -config.file=/etc/promtail/config.yaml
    volumes:
      - ./config/promtail-config.yaml:/etc/promtail/config.yaml:ro
      - /home/ubuntu/.injectived/logs:/home/injectived/logs:ro
      - promtail-positions:/tmp
    restart: unless-stopped
volumes:
  promtail-positions:
EOF

# 启动 Promtail
cd /opt/promtail
docker-compose up -d
```

#### 使用 systemd 运行 Promtail

```bash
# 1. 下载 Promtail 二进制文件
wget https://github.com/grafana/loki/releases/download/v3.0.0/promtail-linux-amd64.zip
unzip promtail-linux-amd64.zip
sudo mv promtail-linux-amd64 /usr/local/bin/promtail
sudo chmod +x /usr/local/bin/promtail

# 2. 创建 systemd 服务文件
sudo cat > /etc/systemd/system/promtail.service <<EOF
[Unit]
Description=Promtail
After=network.target

[Service]
Type=simple
User=ubuntu
ExecStart=/usr/local/bin/promtail -config.file=/opt/promtail/config/promtail-config.yaml
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 3. 启动服务
sudo systemctl daemon-reload
sudo systemctl enable promtail
sudo systemctl start promtail
sudo systemctl status promtail
```

## 验证与测试

### 1. 检查服务状态

```bash
# 检查 Loki
curl http://localhost:3100/ready
# 应该返回: ready

# 检查 Promtail
curl http://localhost:9080/ready
# 应该返回: ready

# 查看容器状态
docker-compose -f compose.log.yaml ps
```

### 2. 查看服务日志

```bash
# Promtail 日志
docker logs promtail -f

# Loki 日志
docker logs loki -f
```

### 3. 测试日志收集

```bash
# 查询所有日志（最近1小时）
curl -s -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10" | jq

# 查询特定节点日志
curl -s -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node",node_id="local-node"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10" | jq

# 查询错误日志
curl -s -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node",level="ERR"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10" | jq
```

### 4. 检查标签

```bash
# 查询所有标签
curl http://localhost:3100/loki/api/v1/labels

# 查询 job 标签的值
curl http://localhost:3100/loki/api/v1/label/job/values

# 查询 node_id 标签的值
curl http://localhost:3100/loki/api/v1/label/node_id/values
```

### 5. WebSocket 实时日志测试

```bash
# 使用 wscat 测试（需要先安装: npm install -g wscat）
wscat -c "ws://localhost:3100/loki/api/v1/tail?query=%7Bjob%3D%22injective-node%22%7D&limit=100"
```

## 日志查询

### LogQL 查询语言

LogQL 是 Loki 的查询语言，类似于 PromQL。

#### 基础查询

```logql
# 查询所有节点日志
{job="injective-node"}

# 查询特定节点日志
{job="injective-node", node_id="node-1"}

# 查询错误日志
{job="injective-node", level="ERR"}

# 查询包含关键词的日志
{job="injective-node"} |= "error"

# 排除特定关键词
{job="injective-node"} != "debug"

# 正则匹配
{job="injective-node"} |~ "error|fatal|panic"
```

#### 统计查询

```logql
# 统计错误数量（1小时内）
sum(count_over_time({job="injective-node", level="ERR"}[1h]))

# 按级别统计日志数量
sum by (level) (count_over_time({job="injective-node"}[1h]))

# 按节点统计日志数量
sum by (node_id) (count_over_time({job="injective-node"}[1h]))

# 错误率（1小时内）
sum(rate({job="injective-node", level="ERR"}[1h])) / 
sum(rate({job="injective-node"}[1h])) * 100

# 日志速率（每秒日志条数）
rate({job="injective-node"}[5m])
```

#### 时间范围查询

```logql
# 查询最近1小时的日志
{job="injective-node"} [1h]

# 查询最近24小时的日志
{job="injective-node"} [24h]

# 使用时间函数
{job="injective-node"} | json | timestamp > now() - 1h
```

### HTTP API 查询

#### 查询范围（Query Range）

```bash
curl -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=100"
```

参数说明:
- `query`: LogQL 查询表达式
- `start`: 开始时间（纳秒 Unix 时间戳）
- `end`: 结束时间（纳秒 Unix 时间戳）
- `limit`: 返回的最大条目数

#### 即时查询（Instant Query）

```bash
curl -G "http://localhost:3100/loki/api/v1/query" \
  --data-urlencode 'query={job="injective-node"}' \
  --data-urlencode "limit=100"
```

#### 标签查询

```bash
# 查询所有标签
curl http://localhost:3100/loki/api/v1/labels

# 查询标签值
curl http://localhost:3100/loki/api/v1/label/node_id/values
```

### WebSocket 实时日志流

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
const ws = connectLokiLiveTail('{job="injective-node"}');
```

### Grafana 查询

1. 打开 Grafana: http://localhost:3000
2. 进入 **Explore**
3. 选择 **Loki** 数据源
4. 输入 LogQL 查询：
   ```
   {job="injective-node"}
   ```
5. 点击 **Run query** 查看日志
6. 点击 **Live** 按钮启用实时日志流

## 故障排查

### 1. Promtail 无法连接 Loki

**症状**: Promtail 日志显示连接错误

**检查步骤**:

```bash
# 检查 Loki 是否运行
curl http://localhost:3100/ready

# 检查网络连接（同一 Docker 网络）
docker exec promtail wget -qO- http://loki:3100/ready

# 检查网络配置
docker network inspect biya-backend | grep -A 5 "loki\|promtail"

# 查看 Promtail 日志
docker logs promtail | tail -50
```

**解决方案**:
- 确保 Loki 先启动
- 检查网络配置（同一 Docker 网络或公网地址）
- 如果节点不在同一网络，修改 Promtail 配置中的 `clients.url` 为 Loki 的公网地址
- 检查防火墙设置（确保 3100 端口开放）

### 2. 日志没有收集

**症状**: 查询不到日志

**检查步骤**:

```bash
# 检查日志文件是否存在
ls -la ~/.injectived/logs/inj.log

# 检查文件权限
stat ~/.injectived/logs/inj.log

# 查看 Promtail 日志
docker logs promtail | tail -50

# 检查 Promtail 配置
docker exec promtail cat /etc/promtail/config.yaml

# 检查 positions 文件
docker exec promtail cat /tmp/positions.yaml
```

**解决方案**:
- 确保日志文件存在且可读
- 检查 Promtail 配置中的路径是否正确
- 检查 Docker volume 挂载是否正确
- 等待几分钟让 Promtail 开始收集日志
- 检查 positions 文件是否正常更新

### 3. 查询不到日志

**症状**: API 查询返回空结果

**检查步骤**:

```bash
# 查询所有标签
curl http://localhost:3100/loki/api/v1/labels

# 查询 job 标签的值
curl http://localhost:3100/loki/api/v1/label/job/values

# 查询 node_id 标签的值
curl http://localhost:3100/loki/api/v1/label/node_id/values

# 扩大时间范围查询
curl -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node"}' \
  --data-urlencode "start=$(date -d '24 hours ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=100"
```

**解决方案**:
- 等待几分钟让 Promtail 开始收集日志
- 检查 LogQL 查询语法
- 扩大查询时间范围
- 检查标签值是否正确

### 4. WebSocket 连接失败

**症状**: 前端 WebSocket 连接失败

**检查步骤**:

```bash
# 测试 WebSocket 连接（需要 wscat）
wscat -c 'ws://localhost:3100/loki/api/v1/tail?query={job="injective-node"}&limit=10'

# 检查 Loki 是否正常运行
curl http://localhost:3100/ready
```

**解决方案**:
- 确保 Loki 正常运行
- 检查 WebSocket URL 格式
- 如果通过代理，确保代理支持 WebSocket（nginx 需要配置 `proxy_set_header Upgrade $http_upgrade;`）
- 检查 CORS 设置（如果需要跨域）

### 5. 性能问题

**症状**: 查询速度慢或内存占用高

**检查步骤**:

```bash
# 查看容器资源使用
docker stats loki promtail

# 查看 Loki 日志
docker logs loki | grep -i error
```

**解决方案**:
- 限制查询时间范围
- 使用标签过滤减少数据量
- 调整 `limits_config` 中的限制
- 使用统计查询代替原始日志查询
- 增加 Loki 的资源配置

### 6. 日志标签未提取

**症状**: 查询时无法使用 level、module 等标签

**检查步骤**:

```bash
# 查看 Promtail 配置中的 pipeline_stages
docker exec promtail cat /etc/promtail/config.yaml | grep -A 20 pipeline_stages

# 查看实际日志格式
tail -20 ~/.injectived/logs/inj.log
```

**解决方案**:
- 检查日志格式是否匹配 pipeline_stages 中的正则表达式
- 调整正则表达式以匹配实际日志格式
- 查看 Promtail 日志确认是否有解析错误

## 性能优化

### 1. 日志保留策略

修改 `configs/loki/loki-config.yaml`:

```yaml
limits_config:
  retention_period: 168h  # 7天，可根据需要调整
  reject_old_samples_max_age: 168h
```

### 2. 写入速率限制

```yaml
limits_config:
  ingestion_rate_mb: 16  # 每秒写入速率（MB）
  ingestion_burst_size_mb: 32  # 突发写入大小（MB）
```

### 3. 查询优化

- 使用标签过滤减少数据量
- 限制查询时间范围
- 使用统计查询代替原始日志查询
- 避免全量扫描

### 4. 日志轮转

配置日志轮转避免单个文件过大:

```bash
# 使用 logrotate
cat > /etc/logrotate.d/injectived <<EOF
/home/ubuntu/.injectived/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0644 ubuntu ubuntu
}
EOF
```

### 5. 多文件收集

如果节点有多个日志文件，可以配置多个 scrape_configs:

```yaml
scrape_configs:
  - job_name: injective-node-app
    static_configs:
      - targets: [localhost]
        labels:
          job: injective-node
          node_id: node-1
          log_type: application
          __path__: /var/log/injective/application.log
  
  - job_name: injective-node-error
    static_configs:
      - targets: [localhost]
        labels:
          job: injective-node
          node_id: node-1
          log_type: error
          __path__: /var/log/injective/error.log
```

## 最佳实践

1. **节点标识**: 每个节点使用唯一的 `node_id`，便于区分和查询
2. **标签设计**: 合理使用标签（job、node_id、level、module），避免标签基数过高
3. **日志格式**: 保持日志格式一致，便于 Promtail 解析
4. **监控告警**: 配置 Loki Ruler 对错误日志进行告警
5. **数据保留**: 根据存储容量和需求设置合理的保留期
6. **性能监控**: 监控 Loki 和 Promtail 的资源使用情况

## 相关文档

- [LOKI_QUICKSTART.md](./LOKI_QUICKSTART.md) - Loki 快速开始指南
- [DEPLOY_LOGS_LOCAL.md](./DEPLOY_LOGS_LOCAL.md) - 本地部署日志收集
- [Loki 官方文档](https://grafana.com/docs/loki/latest/)
- [LogQL 查询语言](https://grafana.com/docs/loki/latest/logql/)
- [Promtail 配置文档](https://grafana.com/docs/loki/latest/clients/promtail/)

## 总结

本文档详细介绍了如何使用 Promtail 从各个节点收集日志并发送到 Loki。通过合理的配置和部署，可以实现：

- ✅ 集中式日志存储和查询
- ✅ 多节点日志统一管理
- ✅ 实时日志流监控
- ✅ 强大的日志查询和分析能力
- ✅ 与 Grafana 集成实现可视化

如有问题，请参考故障排查章节或查看相关文档。
