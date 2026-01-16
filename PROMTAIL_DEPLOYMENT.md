# Promtail 分布式部署文档

## 概述

本文档说明如何在每个验证者节点上部署 Promtail，将日志收集到中央 Loki 服务器。

**Loki 服务器地址**: `10.8.190.46:3100`（已写死）

## 快速部署

### 1. 使用部署脚本（推荐）

```bash
# 基本用法：只需提供日志文件路径
./scripts/deploy-promtail.sh /home/ubuntu/.injectived/logs/inj.log

# 指定节点ID（可选，默认使用主机名）
./scripts/deploy-promtail.sh /home/ubuntu/.injectived/logs/inj.log node-1
```

### 2. 脚本功能

- ✅ 自动生成 Promtail 配置文件
- ✅ 自动部署 Docker 容器
- ✅ 自动挂载日志文件目录
- ✅ 自动配置日志解析管道
- ✅ 自动设置节点标识

## 部署步骤

### 步骤 1: 准备日志文件路径

确认验证者节点的日志文件路径，例如：
- `/home/ubuntu/.injectived/logs/inj.log`
- `/var/log/injective/inj.log`
- 或其他自定义路径

### 步骤 2: 运行部署脚本

```bash
# 进入项目目录
cd /path/to/biya-backend-exporter

# 运行部署脚本
./scripts/deploy-promtail.sh /home/ubuntu/.injectived/logs/inj.log node-1
```

### 步骤 3: 验证部署

```bash
# 检查容器状态
docker ps | grep promtail

# 查看容器日志
docker logs -f promtail

# 健康检查
curl http://localhost:9080/ready
```

## 配置说明

### 自动生成的配置

脚本会在 `/opt/promtail/promtail-config.yaml` 生成配置文件，包含：

- **Loki 地址**: `http://10.8.190.46:3100/loki/api/v1/push`
- **节点标识**: 使用传入的节点ID或主机名
- **日志路径**: 使用传入的日志文件路径
- **日志解析**: 自动提取日志级别、模块等信息

### 修改配置

如需修改配置：

```bash
# 编辑配置文件
vi /opt/promtail/promtail-config.yaml

# 重启容器使配置生效
docker restart promtail
```

## 服务管理

### 查看服务状态

```bash
docker ps | grep promtail
```

### 查看服务日志

```bash
docker logs -f promtail
```

### 重启服务

```bash
docker restart promtail
```

### 停止服务

```bash
docker stop promtail
```

### 删除服务

```bash
docker stop promtail
docker rm promtail
```

## 验证日志收集

### 在 Loki 中查询日志

```bash
# 查询所有节点日志
curl -G "http://10.8.190.46:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10"

# 查询特定节点日志（替换 node-1 为实际节点ID）
curl -G "http://10.8.190.46:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node",node_id="node-1"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10"
```

### 在 Grafana 中查看

1. 打开 Grafana: http://10.8.190.46:3000（或你的 Grafana 地址）
2. 进入 **Explore**
3. 选择 **Loki** 数据源
4. 输入查询: `{job="validator-node"}`
5. 点击 **Run query**

## 故障排查

### 问题 1: 容器无法启动

**检查**:
```bash
docker logs promtail
```

**可能原因**:
- 配置文件格式错误
- 日志文件路径不存在
- 端口 9080 被占用

### 问题 2: 无法连接到 Loki

**检查**:
```bash
# 测试网络连接
curl http://10.8.190.46:3100/ready

# 查看容器日志
docker logs promtail | grep -i error
```

**可能原因**:
- 网络不通
- 防火墙阻止
- Loki 服务未启动

### 问题 3: 日志没有收集

**检查**:
```bash
# 检查日志文件是否存在
ls -la /home/ubuntu/.injectived/logs/inj.log

# 检查文件权限
stat /home/ubuntu/.injectived/logs/inj.log

# 查看 Promtail 日志
docker logs promtail | tail -50
```

**可能原因**:
- 日志文件路径错误
- 文件权限不足
- 日志文件没有新内容

## 多节点部署

在每个验证者节点上重复执行部署脚本：

```bash
# 节点 1
./scripts/deploy-promtail.sh /home/ubuntu/.injectived/logs/inj.log node-1

# 节点 2
./scripts/deploy-promtail.sh /home/ubuntu/.injectived/logs/inj.log node-2

# 节点 3
./scripts/deploy-promtail.sh /home/ubuntu/.injectived/logs/inj.log node-3
```

## 注意事项

1. **节点ID唯一性**: 确保每个节点使用不同的节点ID，便于在 Loki 中区分日志
2. **日志文件路径**: 确保传入的日志文件路径正确且可读
3. **网络连通性**: 确保节点能够访问 Loki 服务器（10.8.190.46:3100）
4. **防火墙**: 如果节点间有防火墙，需要开放 3100 端口（Loki）和 9080 端口（Promtail 健康检查）

## 文件位置

- **配置文件**: `/opt/promtail/promtail-config.yaml`
- **位置文件**: `/opt/promtail/positions/positions.yaml`
- **容器名称**: `promtail`
- **部署脚本**: `./scripts/deploy-promtail.sh`
