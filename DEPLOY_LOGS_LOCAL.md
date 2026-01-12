# 本地部署 Loki + Promtail 日志收集

本指南说明如何在当前机器上部署 Loki + Promtail 来收集 `~/.injectived/logs/inj.log` 日志。

## 配置说明

已配置的路径：
- **日志文件**: `~/.injectived/logs/inj.log` (宿主机路径: `/home/ubuntu/.injectived/logs/inj.log`)
- **容器内路径**: `/home/injectived/logs/inj.log`
- **节点标识**: `local-node`

## 快速启动

### 方式1：使用部署脚本（推荐）

```bash
# 运行部署脚本
./deploy-logs.sh
```

脚本会自动：
1. 检查网络和日志文件
2. 启动 Loki 服务
3. 启动 Promtail 服务
4. 验证服务状态
5. 测试日志收集

### 方式2：手动启动

```bash
# 1. 确保网络存在
docker network ls | grep biya-backend
# 如果不存在，创建网络：
docker network create biya-backend

# 2. 启动 Loki
docker-compose -f compose.log.yaml up -d loki

# 3. 等待 Loki 就绪（约10秒）
curl http://localhost:3100/ready

# 4. 启动 Promtail
docker-compose -f compose.log.yaml up -d promtail

# 5. 验证 Promtail 就绪
curl http://localhost:9080/ready
```

## 验证日志收集

### 1. 查看服务状态

```bash
docker-compose -f compose.log.yaml ps
```

### 2. 查看 Promtail 日志

```bash
# 查看 Promtail 日志
docker-compose -f compose.log.yaml logs -f promtail

# 查看 Loki 日志
docker-compose -f compose.log.yaml logs -f loki
```

### 3. 查询日志（API）

**推荐方式：使用 `-G` 和 `--data-urlencode`（避免引号转义问题）**

```bash
# 查询最近的日志（前10条）
curl -s -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10" | jq -r '.data.result[0].values[] | .[1]'

# 查询特定节点的日志
curl -s -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node",node_id="local-node"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10"

# 查询错误日志
curl -s -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node",level="ERR"}' \
  --data-urlencode "start=$(date -d '1 hour ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10"

#使用 wscat 通过 WebSocket 实时查看日志（需先安装 wscat：npm install -g wscat）
wscat -c "ws://localhost:3100/loki/api/v1/tail?query={job=\"injective-node\"}&limit=100"
```


**或者使用 URL 编码方式：**

```bash
# 查询最近的日志（使用 URL 编码）
curl -s "http://localhost:3100/loki/api/v1/query_range?query=%7Bjob%3D%22injective-node%22%7D&start=$(date -d '1 hour ago' +%s)000000000&end=$(date +%s)000000000&limit=10"
```

### 4. 在 Grafana 中查看

1. 打开 Grafana: http://localhost:3000
2. 进入 **Explore**
3. 选择 **Loki** 数据源
4. 输入 LogQL 查询：
   ```
   {job="injective-node"}
   ```
5. 点击 **Run query** 查看日志
6. 点击 **Live** 按钮启用实时日志流（WebSocket）

### 5. WebSocket 实时日志（前端）

```javascript
// 连接 Loki WebSocket API
const ws = new WebSocket(
  'ws://localhost:3100/loki/api/v1/tail?query={job="injective-node"}&limit=100'

// 使用 wscat 也可以测试 WebSocket 日志流（需要先安装 wscat: npm install -g wscat）
// 示例命令：
// wscat -c "ws://localhost:3100/loki/api/v1/tail?query={job=\"injective-node\"}&limit=100"
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
      console.log({
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
```

## LogQL 查询示例

### 基础查询

```logql
# 查询所有日志
{job="injective-node"}

# 查询特定节点日志
{job="injective-node", node_id="local-node"}

# 查询错误日志
{job="injective-node", level="ERR"}

# 查询包含关键词的日志
{job="injective-node"} |= "error"

# 查询特定模块的日志
{job="injective-node", module="consensus"}
```

### 统计查询

```logql
# 统计错误数量（1小时内）
sum(count_over_time({job="injective-node", level="ERR"}[1h]))

# 按级别统计日志数量
sum by (level) (count_over_time({job="injective-node"}[1h]))

# 按模块统计日志数量
sum by (module) (count_over_time({job="injective-node"}[1h]))

# 错误率
sum(rate({job="injective-node", level="ERR"}[1h])) / sum(rate({job="injective-node"}[1h])) * 100
```

## 服务管理

### 停止服务

```bash
docker-compose -f compose.log.yaml down
```

### 重启服务

```bash
docker-compose -f compose.log.yaml restart promtail
docker-compose -f compose.log.yaml restart loki
```

### 查看服务状态

```bash
docker-compose -f compose.log.yaml ps
```

### 清理数据（注意：会删除所有日志数据）

```bash
docker-compose -f compose.log.yaml down -v
```

## 故障排查

### 1. Promtail 无法连接 Loki

**检查**:
```bash
# 检查 Loki 是否运行
curl http://localhost:3100/ready

# 检查网络
docker network inspect biya-backend | grep -A 5 "loki\|promtail"

# 查看 Promtail 日志
docker logs promtail
```

**解决方案**:
- 确保 Loki 先启动
- 检查网络配置
- 检查防火墙设置

### 2. 日志没有收集

**检查**:
```bash
# 检查日志文件是否存在
ls -la ~/.injectived/logs/inj.log

# 检查文件权限
stat ~/.injectived/logs/inj.log

# 查看 Promtail 日志
docker logs promtail | tail -50

# 检查 Promtail 配置
docker exec promtail cat /etc/promtail/config.yaml
```

**解决方案**:
- 确保日志文件存在且可读
- 检查 Promtail 配置中的路径是否正确
- 检查 Docker volume 挂载是否正确

### 3. 查询不到日志

**检查**:
```bash
# 查询所有标签
curl http://localhost:3100/loki/api/v1/labels

# 查询 job 标签的值
curl http://localhost:3100/loki/api/v1/label/job/values

# 查询 node_id 标签的值
curl http://localhost:3100/loki/api/v1/label/node_id/values
```

**解决方案**:
- 等待几分钟让 Promtail 开始收集日志
- 检查 LogQL 查询语法
- 扩大查询时间范围

### 4. WebSocket 连接失败

**检查**:
```bash
# 测试 WebSocket 连接（需要 wscat）
wscat -c 'ws://localhost:3100/loki/api/v1/tail?query={job="injective-node"}&limit=10'
```

**解决方案**:
- 确保 Loki 正常运行
- 检查 WebSocket URL 格式
- 如果通过代理，确保代理支持 WebSocket

## 数据保留

默认配置：
- 数据保留时间: 7 天（168小时）
- 可以通过修改 `configs/loki/loki-config.yaml` 中的 `limits_config.reject_old_samples_max_age` 调整

## 性能优化

如果日志量很大，可以：
1. 调整 `limits_config` 中的限制
2. 使用更精确的 LogQL 查询过滤日志
3. 增加 Loki 的资源配置

## 下一步

- 集成到 Grafana Dashboard 实现日志可视化
- 配置告警规则（使用 Loki Ruler）
- 添加更多节点的日志收集
- 配置日志轮转策略
