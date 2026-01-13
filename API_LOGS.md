# 日志统计和查询 API 文档

本 API 提供日志级别统计和时间范围查询功能，直接使用 LogQL（Loki Query Language）进行查询。

## API 端点

### 1. 日志统计查询

**端点**: `GET /api/v1/node-logs/stats`

**功能**: 获取指定时间范围内按日志级别统计的数据

**查询参数**:
- `node_id` (可选): 节点 ID，例如 `node-1`。如果不提供，查询所有节点
- `start` (可选): 开始时间，Unix 时间戳（秒）。默认：24小时前
- `end` (可选): 结束时间，Unix 时间戳（秒）。默认：当前时间

**响应格式**:
```json
{
  "success": true,
  "data": {
    "total": 1879,      // 总日志数
    "error": 5,         // 错误日志数
    "warning": 42,      // 警告日志数
    "info": 1832,       // 信息日志数
    "debug": 0,         // 调试日志数
    "time_range": {
      "start": 1705123200,  // 开始时间戳（秒）
      "end": 1705209600     // 结束时间戳（秒）
    }
  }
}
```

**示例请求**:
```bash
# 查询最近24小时的统计
curl "http://localhost:18080/api/v1/node-logs/stats"

# 查询特定节点的统计
curl "http://localhost:18080/api/v1/node-logs/stats?node_id=local-node"

# 查询指定时间范围的统计
curl "http://localhost:18080/api/v1/node-logs/stats?start=1705123200&end=1705209600"

# 查询特定节点在指定时间范围的统计
curl "http://localhost:18080/api/v1/node-logs/stats?node_id=local-node&start=1705123200&end=1705209600"
```

### 2. 日志列表查询

**端点**: `GET /api/v1/node-logs/list`

**功能**: 获取日志列表，支持按级别、关键词、时间范围过滤和分页

**查询参数**:
- `node_id` (可选): 节点 ID
- `level` (可选): 日志级别，可选值：`all`, `error`, `warning`, `info`, `debug`。默认：`all`
- `keyword` (可选): 搜索关键词，在日志内容中搜索
- `start` (可选): 开始时间，Unix 时间戳（秒）。默认：24小时前
- `end` (可选): 结束时间，Unix 时间戳（秒）。默认：当前时间
- `page` (可选): 页码，从 1 开始。默认：1
- `page_size` (可选): 每页大小，最大 200。默认：50

**响应格式**:
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "timestamp": 1705209600000,  // Unix 时间戳（毫秒）
        "level": "error",
        "content": "ERROR consensus: Failed to reach consensus",
        "labels": {
          "job": "injective-node",
          "node_id": "local-node",
          "module": "consensus"
        }
      }
    ],
    "total": 1879,      // 总条数
    "page": 1,          // 当前页码
    "page_size": 50,    // 每页大小
    "has_more": true    // 是否有更多数据
  }
}
```

**示例请求**:
```bash
# 查询最近的日志（默认50条）
curl "http://localhost:18080/api/v1/node-logs/list"

# 查询错误日志
curl "http://localhost:18080/api/v1/node-logs/list?level=error"

# 搜索包含关键词的日志
curl "http://localhost:18080/api/v1/node-logs/list?keyword=consensus"

# 分页查询
curl "http://localhost:18080/api/v1/node-logs/list?page=1&page_size=100"

# 组合查询：特定节点的错误日志，包含关键词
curl "http://localhost:18080/api/v1/node-logs/list?node_id=local-node&level=error&keyword=timeout"
```

## LogQL 查询说明

本 API 直接使用 LogQL（Loki Query Language）进行查询，类似于 SQL 的聚合查询。

### 统计查询使用的 LogQL

1. **总日志数**:
   ```logql
   sum(count_over_time({job="injective-node"}[1h]))
   ```

2. **按级别统计**:
   ```logql
   sum by (level) (count_over_time({job="injective-node"}[1h]))
   ```

3. **错误日志数**:
   ```logql
   sum(count_over_time({job="injective-node",level=~"ERR|ERROR|FATAL|PANIC"}[1h]))
   ```

### 日志列表查询使用的 LogQL

1. **基础查询**:
   ```logql
   {job="injective-node"}
   ```

2. **按节点过滤**:
   ```logql
   {job="injective-node",node_id="local-node"}
   ```

3. **按级别过滤**:
   ```logql
   {job="injective-node",level=~"ERR|ERROR|FATAL|PANIC"}
   ```

4. **关键词搜索**:
   ```logql
   {job="injective-node"} |= "error"
   ```

5. **组合查询**:
   ```logql
   {job="injective-node",node_id="local-node",level=~"ERR|ERROR"} |= "timeout"
   ```

## 时间范围说明

时间范围参数使用 Unix 时间戳（秒）：

- **最近1小时**: `start=$(date -d '1 hour ago' +%s)&end=$(date +%s)`
- **最近24小时**: `start=$(date -d '24 hours ago' +%s)&end=$(date +%s)`
- **最近7天**: `start=$(date -d '7 days ago' +%s)&end=$(date +%s)`
- **自定义范围**: `start=1705123200&end=1705209600`

## 日志级别说明

支持的日志级别：
- **error**: ERR, ERROR, FATAL, PANIC
- **warning**: WARN, WARNING
- **info**: INF, INFO
- **debug**: DBG, DEBUG

## 前端集成示例

### JavaScript/TypeScript

```javascript
// 获取日志统计
async function getLogStats(nodeId, start, end) {
  const params = new URLSearchParams();
  if (nodeId) params.set('node_id', nodeId);
  if (start) params.set('start', start.toString());
  if (end) params.set('end', end.toString());

  const response = await fetch(`http://localhost:18080/api/v1/node-logs/stats?${params}`);
  const data = await response.json();
  
  if (data.success) {
    console.log('总日志数:', data.data.total);
    console.log('错误日志:', data.data.error);
    console.log('警告日志:', data.data.warning);
    console.log('信息日志:', data.data.info);
  }
  
  return data;
}

// 获取日志列表
async function getLogList(nodeId, level, keyword, page = 1, pageSize = 50) {
  const params = new URLSearchParams();
  if (nodeId) params.set('node_id', nodeId);
  if (level && level !== 'all') params.set('level', level);
  if (keyword) params.set('keyword', keyword);
  params.set('page', page.toString());
  params.set('page_size', pageSize.toString());

  const response = await fetch(`http://localhost:18080/api/v1/node-logs/list?${params}`);
  const data = await response.json();
  
  if (data.success) {
    console.log('日志列表:', data.data.logs);
    console.log('总数:', data.data.total);
    console.log('是否有更多:', data.data.has_more);
  }
  
  return data;
}

// 使用示例
const stats = await getLogStats('local-node', Date.now() / 1000 - 86400, Date.now() / 1000);
const logs = await getLogList('local-node', 'error', 'timeout', 1, 100);
```

## 性能优化

1. **时间范围**: 建议查询时间范围不超过 7 天，避免查询超时
2. **分页**: 使用合理的 `page_size`（建议 50-100），避免单次查询数据量过大
3. **缓存**: 前端可以实现缓存机制，减少重复查询
4. **索引**: Loki 会根据标签（job, node_id, level）自动建立索引，查询速度很快

## 错误处理

API 返回错误时，`success` 字段为 `false`，`message` 字段包含错误信息：

```json
{
  "success": false,
  "message": "查询日志失败: connection timeout"
}
```

HTTP 状态码：
- `200`: 成功
- `400`: 请求参数错误
- `500`: 服务器内部错误

## 与 Loki 直接查询的对比

| 方式 | 优点 | 缺点 |
|------|------|------|
| **通过 API** | 统一接口、参数验证、错误处理、响应格式统一 | 多一层转发 |
| **直接查询 Loki** | 性能最优、功能完整 | 需要了解 LogQL、需要处理错误格式 |

**推荐**: 前端应用使用本 API，统一接口和错误处理；高级用户可以直接查询 Loki。
