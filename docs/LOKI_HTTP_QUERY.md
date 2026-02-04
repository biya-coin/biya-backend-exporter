# 通过 HTTP 查询 Loki 日志

本文档说明通过 HTTP 使用的四个接口：

1. **获取 job 列表**：查询当前有日志的 job 标签值，用于下拉或筛选
2. **获取 node_id 列表**：按 job（可选）和时间段查询 node_id 列表，用于选择节点
3. **按 node_id 和时间段统计各级别数量**：返回 error、warn、info 的条数
4. **按 node_id、时间段和可选关键词搜索日志**：关键词可为空，为空时返回该范围内全部日志

---

## 一、前置说明

### 1.1 Loki 基础地址与 API

- **基础地址**：如 `https://prv.loki.biya.io` 或 `http://loki:3100`（容器内）
- **范围查询**：`GET /loki/api/v1/query_range`（带时间范围）
- **时间戳**：Loki 要求 **纳秒** 级 Unix 时间戳（秒级时间戳 × 10⁹）

### 1.2 常用标签（与本项目 Promtail 一致）

| 标签       | 说明                     | 示例值              |
|------------|--------------------------|---------------------|
| `job`      | 任务名（同一类服务/集群的统称） | `validator-node`（所有 Injective 节点共用） |
| `node_id`  | 节点唯一标识（区分同一集群里的不同实例） | `validator-0`, `local-node` |
| `level`    | 日志级别（Promtail 解析）| `INF`, `ERR`, `WRN`, `DBG`, `FATAL`, `PANIC` |

### 1.3 时间范围与纳秒时间戳

| 范围   | 起点相对当前 | 纳秒 start 示例（Bash） |
|--------|----------------|---------------------------|
| 1 小时 | 1h ago         | `$(($(date -d '1 hour ago' +%s) * 1000000000))` |
| 6 小时 | 6h ago         | `$(($(date -d '6 hours ago' +%s) * 1000000000))` |
| 24 小时| 24h ago        | `$(($(date -d '24 hours ago' +%s) * 1000000000))` |
| 7 天   | 7d ago         | `$(($(date -d '7 days ago' +%s) * 1000000000))` |

**当前时间纳秒**（end）：

```bash
END_NS=$(($(date +%s) * 1000000000))
```

---

## 二、接口一：获取 job 列表

用于获取 Loki 中存在的 **job** 标签取值列表（例如 `validator-node`），便于前端下拉或筛选。

### 2.1 直接调用 Loki HTTP API

**接口**：`GET /loki/api/v1/label/job/values`

可选参数 `start`、`end`（纳秒时间戳）限制在某一时间范围内出现过的 job；不传则返回所有曾出现过的 job。

**示例**：

```bash
LOKI_URL="https://prv.loki.biya.io"   # 按需修改

# 全部 job
curl -s -G "${LOKI_URL}/loki/api/v1/label/job/values"

```

**响应示例**：

```json
{
  "status": "success",
  "data": [ "validator-node" ]
}
```

---

## 三、接口二：获取 node_id 列表

根据 **job**（可选）和 **时间范围**，获取在该时间段内有日志的 **node_id** 列表，用于选择节点。

### 3.1 直接调用 Loki HTTP API

**接口**：`GET /loki/api/v1/label/node_id/values`

**示例**：

```bash
LOKI_URL="https://prv.loki.biya.io"

curl -s -G "${LOKI_URL}/loki/api/v1/label/node_id/values" 
```

**响应示例**：

```json
{
  "status":"success",
  "data":["sentry-0","sentry-1","validator-0","validator-1","validator-2","validator-3"]
}
```


---

## 四、接口三：按 node_id 和时间段统计 error / warn / info 数量

根据 **node_id** 和 **时间范围**，统计各级别日志条数（仅 error、warn、info）。

### 4.1 直接调用 Loki HTTP API

使用 LogQL：`sum by (level) (count_over_time({...}[duration]))`，再在结果中只取 `ERR`、`WRN`、`INF` 对应条数。

**为何每个级别会有多条？**  
`query_range` 会按 **step** 多次求值：不传 `step` 时 Loki 用默认步长，时间范围越长，每个 level 的 `values` 里就会出现多个 `[时间戳, 条数]`。若希望**每个级别在该时间范围内只得到一条（即该时段内的总数）**，需要显式传 `step` = 整个查询范围的**秒数**（Loki 的 step 为秒或时长字符串，纳秒会报错 overflow）。

**示例：指定节点、最近 1 小时（每个级别一条）**：

```bash
LOKI_URL="https://prv.loki.biya.io"   # 按需修改
START_NS=$(($(date -d '1 hour ago' +%s) * 1000000000))
END_NS=$(($(date +%s) * 1000000000))
# step 单位为秒：整段范围的秒数，只求值一次，每个 level 只返回一条
STEP_S=$(((END_NS - START_NS) / 1000000000))

curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode 'query=sum by (level) (count_over_time({job="validator-node",node_id="validator-0"}[1h]))' \
  --data-urlencode "start=${START_NS}" \
  --data-urlencode "end=${END_NS}" \
  --data-urlencode "step=${STEP_S}"
```

**响应示例**：每个级别一条，从 `data.result[]` 中根据 `metric.level` 取 ERR、WRN、INF 的条数（`values[0][1]` 即该级别总数）。

```json
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      { "metric": { "level": "ERR" },  "values": [ [ 1738657200000000000, "5" ] ] },
      { "metric": { "level": "INF" },  "values": [ [ 1738657200000000000, "120" ] ] },
      { "metric": { "level": "WRN" },  "values": [ [ 1738657200000000000, "3" ] ] }
    ]
  }
}
```

---

## 五、接口四：按 node_id、时间段和可选关键词搜索日志

根据 **node_id**、**时间范围** 以及 **搜索关键词**（可为空）查询日志。关键词为空时返回该节点在该时间范围内的全部日志。

### 5.1 直接调用 Loki HTTP API

- 指定节点：`{job="validator-node",node_id="validator-0"}`
- 有关键词时在 query 后加 `|= "关键词"`；无关键词则不加重度过滤，只按流与时间查询。

**示例：指定节点、最近 1 小时、无关键词（全部日志）**：

```bash
LOKI_URL="https://prv.loki.biya.io"
START_NS=$(($(date -d '1 hour ago' +%s) * 1000000000))
END_NS=$(($(date +%s) * 1000000000))

curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node",node_id="validator-0"}' \
  --data-urlencode "start=${START_NS}" \
  --data-urlencode "end=${END_NS}" \
  --data-urlencode "limit=50"
```

**示例：带关键词搜索**（例如包含 "timeout"）：

```bash
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node",node_id="validator-0"} |= "Timed"' \
  --data-urlencode "start=${START_NS}" \
  --data-urlencode "end=${END_NS}" \
  --data-urlencode "limit=5"
```

**示例：仅 error 级别 + 关键词搜索**（与 Grafana 日志大盘一致，按**日志行内容**匹配，不依赖 stream 的 `level` 标签）：

```bash
# 用 |~ 正则匹配日志行中的 ERR/FATAL/PANIC，再用 |= 限制关键词；job 可选，有 node_id 即可
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node",node_id="validator-0"} |~ "(ERR|FATAL|PANIC)" |= "Failed"' \
  --data-urlencode "start=${START_NS}" \
  --data-urlencode "end=${END_NS}" \
  --data-urlencode "limit=50"
```

**示例：仅 warn / 仅 info 级别**（同样用行内正则匹配日志内容）：

```bash
# 仅 warn：行内包含 WRN
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node",node_id="validator-0"} |~ "WRN"' \
  --data-urlencode "start=${START_NS}" \
  --data-urlencode "end=${END_NS}" \
  --data-urlencode "limit=50"

# 仅 info：行内包含 INF
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node",node_id="validator-0"} |~ "INF"' \
  --data-urlencode "start=${START_NS}" \
  --data-urlencode "end=${END_NS}" \
  --data-urlencode "limit=50"
```

按级别过滤推荐用 **行内正则**：**error** 用 `|~ "(ERR|FATAL|PANIC)"`，**warn** 用 `|~ "WRN"`，**info** 用 `|~ "INF"`。这样依赖日志行里出现的级别关键字，不要求 Promtail 写入 `level` 标签。若流上已有 `level` 标签，也可用 `level=~"ERR|..."`，但不少环境下行内匹配更稳定。

响应为 `resultType: streams`，每条日志在 `data.result[].values` 中，格式为 `[ 时间戳纳秒, 日志行文本 ]`。



---

## 六、快速参考

| 范围   | start（纳秒，Loki 直连） | start（秒，Exporter） |
|--------|---------------------------|------------------------|
| 1 小时 | `$(($(date -d '1 hour ago' +%s)*1000000000))` | `$(($(date +%s) - 3600))` |
| 6 小时 | `$(($(date -d '6 hours ago' +%s)*1000000000))` | `$(($(date +%s) - 21600))` |
| 24 小时| `$(($(date -d '24 hours ago' +%s)*1000000000))` | `$(($(date +%s) - 86400))` |
| 7 天   | `$(($(date -d '7 days ago' +%s)*1000000000))` | `$(($(date +%s) - 604800))` |

---
