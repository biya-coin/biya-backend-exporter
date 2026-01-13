# Loki 日志下载指南

本指南介绍如何从 Loki 下载完整的日志数据。

## 方法一：使用 Python 脚本（推荐）

### 安装依赖

```bash
# 确保已安装 Python 3 和 requests 库
pip3 install requests
```

### 基本用法

```bash
# 下载最近24小时的所有日志（JSON格式）
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24

# 下载指定时间范围的日志
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --start "2024-01-01T00:00:00Z" \
  --end "2024-01-02T00:00:00Z"

# 导出为文本格式（类似日志文件）
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --format txt

# 导出为 CSV 格式
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --format csv

# 指定输出文件
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --output my_logs.json

# 连接远程 Loki 服务器
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --loki-url "http://45.249.245.183:3100"
```

### 高级用法

```bash
# 下载错误日志
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"} |~ "(ERR|FATAL|PANIC)"' \
  --hours 24

# 下载特定节点的日志
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node", node_id="node-1"}' \
  --hours 24

# 限制最大下载条目数（避免下载过多数据）
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --max-entries 10000

# 调整每次请求的条目数（默认5000，最大10000）
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --limit-per-request 10000

# 增加超时时间（适用于大量日志）
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 168 \
  --timeout 600
```

### 脚本特性

- ✅ **自动分页**: 自动处理分页，下载所有日志
- ✅ **时间窗口**: 智能分割时间范围，避免单次查询过大
- ✅ **进度显示**: 实时显示下载进度
- ✅ **多种格式**: 支持 JSON、文本、CSV 格式导出
- ✅ **错误处理**: 自动重试和错误恢复
- ✅ **断点续传**: 支持限制最大条目数，可分批下载

## 方法二：使用 Bash 脚本（简单快速）

```bash
# 下载最近24小时的日志
./scripts/download_loki_logs.sh '{job="injective-node"}' 24

# 下载最近7天的日志
./scripts/download_loki_logs.sh '{job="injective-node"}' 168

# 指定输出文件
LOKI_URL="http://localhost:3100" \
./scripts/download_loki_logs.sh '{job="injective-node"}' 24 logs.json
```

**注意**: Bash 脚本只下载第一页（最多5000条），如需完整下载请使用 Python 脚本。

## 方法三：直接使用 curl（手动）

### 基本查询

```bash
# 设置变量
LOKI_URL="http://localhost:3100"
QUERY='{job="injective-node"}'
START_TIME=$(date -d '24 hours ago' +%s)000000000
END_TIME=$(date +%s)000000000
LIMIT=5000

# 查询日志
curl -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode "query=${QUERY}" \
  --data-urlencode "start=${START_TIME}" \
  --data-urlencode "end=${END_TIME}" \
  --data-urlencode "limit=${LIMIT}" \
  -o logs.json
```

### 查看结果

```bash
# 查看 JSON 格式的日志
cat logs.json | jq '.'

# 提取日志内容
cat logs.json | jq -r '.data.result[]?.values[]?[1]'

# 统计日志条数
cat logs.json | jq '[.data.result[]?.values[]?] | length'

# 转换为文本格式
cat logs.json | jq -r '.data.result[] | .values[] | "\(.[0] | tonumber / 1000000 | strftime("%Y-%m-%d %H:%M:%S")) \(.[1])"'
```

## 方法四：使用 LogCLI（Loki 官方工具）

### 安装 LogCLI

```bash
# 下载 LogCLI
wget https://github.com/grafana/loki/releases/download/v3.0.0/logcli-linux-amd64.zip
unzip logcli-linux-amd64.zip
chmod +x logcli-linux-amd64
sudo mv logcli-linux-amd64 /usr/local/bin/logcli
```

### 使用 LogCLI

```bash
# 设置 Loki 地址
export LOKI_ADDR=http://localhost:3100

# 查询最近24小时的日志
logcli query '{job="injective-node"}' --since=24h --output=jsonl > logs.jsonl

# 查询指定时间范围
logcli query '{job="injective-node"}' \
  --from="2024-01-01T00:00:00Z" \
  --to="2024-01-02T00:00:00Z" \
  --output=jsonl > logs.jsonl

# 导出为文本格式
logcli query '{job="injective-node"}' --since=24h --output=raw > logs.txt

# 限制条数
logcli query '{job="injective-node"}' --since=24h --limit=10000 --output=jsonl > logs.jsonl
```

## 导出格式说明

### JSON 格式

```json
[
  {
    "timestamp": 1704067200000,
    "timestamp_ns": 1704067200000000000,
    "labels": {
      "job": "injective-node",
      "node_id": "node-1",
      "level": "info"
    },
    "content": "日志内容..."
  }
]
```

### 文本格式

```
[2024-01-01T00:00:00] [job=injective-node node_id=node-1 level=info] 日志内容...
[2024-01-01T00:00:01] [job=injective-node node_id=node-1 level=error] 错误日志...
```

### CSV 格式

```csv
timestamp,timestamp_iso,labels,content
1704067200000,2024-01-01T00:00:00,"{""job"":""injective-node""}","日志内容..."
```

## 常见查询示例

### 查询所有日志

```bash
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24
```

### 查询错误日志

```bash
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"} |~ "(ERR|FATAL|PANIC|ERROR)"' \
  --hours 24
```

### 查询特定节点

```bash
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node", node_id="node-1"}' \
  --hours 24
```

### 查询包含关键词的日志

```bash
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"} |= "transaction"' \
  --hours 24
```

### 查询特定级别的日志

```bash
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node", level="error"}' \
  --hours 24
```

## 性能优化建议

1. **限制时间范围**: 不要查询过长时间范围的日志，建议每次查询不超过7天
2. **使用标签过滤**: 尽量使用标签（如 `node_id`、`level`）来减少数据量
3. **分批下载**: 如果日志量很大，使用 `--max-entries` 参数分批下载
4. **调整 limit**: 根据网络情况调整 `--limit-per-request`（最大10000）

## 故障排查

### 问题：连接超时

```bash
# 增加超时时间
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --timeout 600
```

### 问题：内存不足

```bash
# 限制最大条目数，分批下载
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --max-entries 10000
```

### 问题：Loki 返回 429（请求过多）

```bash
# 减少每次请求的条目数
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --limit-per-request 1000
```

### 问题：无法连接到 Loki

```bash
# 检查 Loki 是否运行
curl http://localhost:3100/ready

# 检查网络连接
ping loki-server-ip

# 使用正确的 URL
python3 scripts/download_loki_logs.py \
  --query '{job="injective-node"}' \
  --hours 24 \
  --loki-url "http://loki-server-ip:3100"
```

## 相关文档

- [Loki 快速开始指南](./LOKI_QUICKSTART.md)
- [Promtail + Loki 设置文档](./PROMTAIL_LOKI_SETUP.md)
- [Loki API 文档](https://grafana.com/docs/loki/latest/api/)
- [LogQL 查询语言](https://grafana.com/docs/loki/latest/logql/)
