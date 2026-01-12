# Promtail 日志收集配置

Promtail 是 Loki 的日志收集代理，部署在各个 Injective 节点上，用于收集日志并发送到 Loki。

## 部署说明

### 单节点部署（同一 Docker 网络）

如果所有节点都在同一个 Docker 网络中，可以使用 `compose.log.yaml` 中的 `promtail-example` 服务作为模板。

### 多节点部署（独立节点）

如果节点在不同的服务器上，需要在每个节点上独立部署 Promtail：

#### 方式1：Docker Compose（推荐）

在每个节点上创建 `docker-compose.yml`：

```yaml
version: '3'

services:
  promtail:
    image: grafana/promtail:3.0.0
    container_name: promtail
    command: -config.file=/etc/promtail/config.yaml
    volumes:
      - ./promtail-config.yaml:/etc/promtail/config.yaml:ro
      - /var/log/injective:/var/log/injective:ro  # 根据实际路径修改
      - /tmp:/tmp  # positions 文件
    restart: unless-stopped
    networks:
      - default
    # 如果 Loki 不在同一网络，需要配置外部网络或使用公网地址

networks:
  default:
    external: true
    name: biya-network  # 或使用其他网络
```

#### 方式2：systemd 服务

创建 `/etc/systemd/system/promtail.service`：

```ini
[Unit]
Description=Promtail
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/promtail -config.file=/etc/promtail/config.yaml
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## 配置说明

### promtail-config.yaml

每个节点的配置需要修改以下内容：

1. **node_id**: 节点的唯一标识
   ```yaml
   labels:
     node_id: node-1  # 修改为 node-2, node-3 等
   ```

2. **日志路径**: 实际的日志文件路径
   ```yaml
   __path__: /var/log/injective/*.log  # 根据实际路径修改
   ```

3. **Loki 地址**: Loki 服务器的地址
   ```yaml
   clients:
     - url: http://loki:3100/loki/api/v1/push  # Docker 网络
       # 或者使用公网地址：
       # - url: http://your-loki-server:3100/loki/api/v1/push
   ```

## 日志格式解析

### JSON 格式日志

如果日志是 JSON 格式，使用 `json` 解析器：

```yaml
pipeline_stages:
  - json:
      expressions:
        level: level
        message: message
        timestamp: timestamp
  - labels:
      level:
  - timestamp:
      source: timestamp
      format: RFC3339
```

### 文本格式日志

如果日志是普通文本，可以使用正则表达式解析：

```yaml
pipeline_stages:
  - regex:
      expression: '^(?P<timestamp>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+\[(?P<level>\w+)\]\s+(?P<message>.*)$'
  - labels:
      level:
  - timestamp:
      source: timestamp
      format: '2006-01-02 15:04:05'
```

## 多日志文件配置

如果节点有多个日志文件，可以添加多个 `static_configs`：

```yaml
scrape_configs:
  - job_name: injective-node
    static_configs:
      # 应用日志
      - targets: [localhost]
        labels:
          job: injective-node
          node_id: node-1
          log_type: application
          __path__: /var/log/injective/application.log
      # 错误日志
      - targets: [localhost]
        labels:
          job: injective-node
          node_id: node-1
          log_type: error
          __path__: /var/log/injective/error.log
      # 访问日志
      - targets: [localhost]
        labels:
          job: injective-node
          node_id: node-1
          log_type: access
          __path__: /var/log/injective/access.log
```

## 验证和调试

### 检查 Promtail 状态

```bash
# 健康检查
curl http://localhost:9080/ready

# 查看指标
curl http://localhost:9080/metrics
```

### 查看位置文件

Promtail 使用 `positions.yaml` 记录读取位置，用于断点续传：

```bash
cat /tmp/positions.yaml
```

### 查看日志

```bash
# Docker
docker logs promtail

# systemd
journalctl -u promtail -f
```

## 注意事项

1. **日志路径权限**: 确保 Promtail 有权限读取日志文件
2. **网络连接**: 确保 Promtail 可以访问 Loki 服务器
3. **标签唯一性**: 每个节点的 `node_id` 必须唯一
4. **日志轮转**: Promtail 支持日志轮转，但需要确保路径匹配（使用通配符）
5. **资源消耗**: Promtail 资源消耗较低，但大量日志文件可能增加 CPU 和内存使用

## 故障排查

### Promtail 无法连接 Loki

- 检查网络连接：`curl http://loki:3100/ready`
- 检查 Loki 地址配置是否正确
- 如果是跨网络，确保防火墙允许访问

### 日志没有收集

- 检查日志路径是否正确
- 检查文件权限
- 查看 Promtail 日志：`docker logs promtail`
- 检查 `positions.yaml` 文件

### 日志格式解析失败

- 检查 `pipeline_stages` 配置
- 使用 LogQL 测试查询：`{job="injective-node"}`
- 查看原始日志内容，调整解析规则
