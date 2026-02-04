# 通过 HTTP 查询 Loki 日志

本文档说明通过 HTTP 使用的四个接口：

1. **获取 job 列表**：查询当前有日志的 job 标签值，用于下拉或筛选
2. **获取 node_id 列表**：按 job（可选）和时间段查询 node_id 列表，用于选择节点
3. **按 node_id 和时间段统计各级别数量**：返回 error、warn、info 的条数
4. **按 node_id、时间段和可选关键词搜索日志**：关键词可为空，为空时返回该范围内全部日志

---

## 一、前置说明

### 1.1 Loki 基础地址与 API

- **基础地址**：如 `https://prv.backend.biya.io` 或 `http://loki:3100`（容器内）
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
LOKI_URL="https://prv.backend.biya.io"   # 按需修改

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
LOKI_URL="https://prv.backend.biya.io"

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
LOKI_URL="https://prv.backend.biya.io"   # 按需修改
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
LOKI_URL="https://prv.backend.biya.io"
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

## 六、指定时间范围内指定节点所有日志下载

在接口四的基础上，通过**指定时间范围**和**指定 node_id** 拉取该节点在该时段内的**全部日志**并保存到本地文件。

### 6.1 单次请求拉取（数据量不大时）

- 使用 `query_range`，不加重度过滤，仅用 `job` + `node_id` 限定流。
- `limit` 控制单次返回条数：不传时使用 Loki 默认（通常较小）；要尽量多取可显式传较大值，Loki 单次上限一般为 **5000** 条（视服务端配置而定）。

**示例：最近 1 小时、指定节点、直接下载为原始日志文件（.log，每行一条日志，与常见日志文件一致）**：

```bash
LOKI_URL="https://prv.backend.biya.io"
NODE_ID="validator-0"
START_NS=$(($(date -d '1 hour ago' +%s) * 1000000000))
END_NS=$(($(date +%s) * 1000000000))

# 请求后经 jq 只取日志行，按时间顺序输出为原始日志文件
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
  --data-urlencode 'query={job="validator-node",node_id="'"${NODE_ID}"'"}' \
  --data-urlencode "start=${START_NS}" \
  --data-urlencode "end=${END_NS}" \
  --data-urlencode "limit=5000" \
  | jq -r '.data.result[].values[]? | .[1]' > "${NODE_ID}-logs.log"
```

得到的 `validator-0-logs.log` 即为原始 log 风格：每行一条日志内容，无 JSON 结构，可直接用 `less`、`grep` 等查看。

**若需保留完整 API 响应（JSON）**：去掉管道与 jq，改为 `-o "${NODE_ID}-logs.json"`；再从 JSON 提取日志行可用：  
`jq -r '.data.result[].values[]? | .[1]' "${NODE_ID}-logs.json" > "${NODE_ID}-logs.log"`

### 6.2 数据量超过单次 limit 时分页拉取

当该时间范围内日志条数超过 5000（或服务端配置的上限）时，需要**分页**：

- 使用 `direction=forward`（默认）：按时间正序，每次请求返回的最后一笔时间戳可作为下一页的 `start`，直到返回条数 &lt; limit。
- 或使用 `direction=backward`：按时间倒序，最后一笔时间戳作为下一页的 `end`。

**示例：分页拉取并合并为一份 JSON（forward，仅示意结构）**：

```bash
LOKI_URL="https://prv.backend.biya.io"
NODE_ID="validator-0"
START_NS=$(($(date -d '6 hours ago' +%s) * 1000000000))
END_NS=$(($(date +%s) * 1000000000))
LIMIT=5000
PAGE=1

while true; do
  curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
    --data-urlencode 'query={job="validator-node",node_id="'"${NODE_ID}"'"}' \
    --data-urlencode "start=${START_NS}" \
    --data-urlencode "end=${END_NS}" \
    --data-urlencode "limit=${LIMIT}" \
    --data-urlencode "direction=forward" \
    -o "page-${PAGE}.json"
  count=$(jq '[.data.result[].values[]?] | length' "page-${PAGE}.json")
  [ "$count" -lt "$LIMIT" ] && break
  # 取本页最后一条的时间戳作为下一页 start（纳秒）
  START_NS=$(jq -r '[.data.result[].values[]?[0]] | max' "page-${PAGE}.json")
  PAGE=$((PAGE + 1))
done
# 此处可将 page-*.json 合并或仅提取 values 到单一文件
```

### 6.3 前端下载按钮实现

前端需要“下载日志”按钮时，有两种常见做法：**前端直接请求并触发下载**，或**后端提供下载接口由浏览器直接下载**。

#### 方式一：前端请求 API，解析后触发下载（推荐）

前端调用 Loki 的 `query_range`（或你们自己的后端代理该接口），拿到 JSON 后抽成原始日志文本，用 Blob + 临时 `<a>` 触发浏览器下载。

**要点**：

1. **请求**：对 `/loki/api/v1/query_range` 发 GET，参数 `query`、`start`、`end`、`limit`（如 5000）。若 Loki 与前端不同域，需后端代理或配置 CORS。
2. **解析**：响应 `data.result[]` 中每个元素的 `values` 为 `[[时间戳纳秒, 日志行], ...]`，只取每项的 `[1]` 即日志行。
3. **合并与下载**：把所有日志行用换行符拼成一段文本，转成 Blob，再通过创建临时 `<a>` 的 `download` 触发下载。

**示例（JavaScript / TypeScript）**：

```javascript
async function downloadNodeLogs(nodeId, startNs, endNs, job = 'validator-node') {
  const params = new URLSearchParams({
    query: `{job="${job}",node_id="${nodeId}"}`,
    start: String(startNs),
    end: String(endNs),
    limit: '5000',
  });
  const url = `${LOKI_BASE_URL}/loki/api/v1/query_range?${params}`;
  const res = await fetch(url);
  const json = await res.json();

  const lines = (json?.data?.result ?? []).flatMap((stream) =>
    (stream.values ?? []).map((v) => v[1])
  );
  const text = lines.join('\n');

  const blob = new Blob([text], { type: 'text/plain;charset=utf-8' });
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = `${nodeId}-logs.log`;
  a.click();
  URL.revokeObjectURL(a.href);
}
```

- 时间范围由前端传入纳秒 `startNs`、`endNs`（可与现有时间选择器统一，选好后转纳秒）。
- 若数据量可能超过 5000 条，可在前端做分页循环：用上一页最后一条的 `values[][0]` 作为下一页的 `start`，直到某页条数 &lt; limit，再把各页 `lines` 合并后一起下载。

#### 方式二：后端提供“下载日志”接口

后端封装对 Loki 的 `query_range` 调用，将结果整理成原始日志文本后，以 **流式或整段** 返回，并设置响应头让浏览器直接下载文件：

- `Content-Type: text/plain; charset=utf-8`
- `Content-Disposition: attachment; filename="<node_id>-logs.log"`

前端下载按钮只需跳转到该接口或使用 `<a href="...">` / `window.location.href`，无需在浏览器里解析 JSON。适合不希望暴露 Loki 地址、或需要鉴权、限流的场景。

**后端响应示例（Go 伪代码）**：

```go
w.Header().Set("Content-Type", "text/plain; charset=utf-8")
w.Header().Set("Content-Disposition", `attachment; filename="`+nodeID+`-logs.log"`)
// 从 Loki 拉取 query_range，解析 data.result[].values，逐行写入 w
```

**前端**：按钮指向该 URL（可带 `node_id`、`start`、`end` 等 query），例如：

```html
<a :href="`/api/logs/download?node_id=${nodeId}&start=${startNs}&end=${endNs}`" download>下载日志</a>
```

或用 `window.open` / `location.href` 同理；若需 POST 或 Header 鉴权，可改为 `fetch` 拿到 `blob()` 后再用上面的 Blob + `<a>.download` 方式触发下载。

---

### 6.4 小结

| 目的           | 做法 |
|----------------|------|
| 指定时间 + 节点 | `query_range`，`query={job="validator-node",node_id="<节点>"}`，`start`/`end` 纳秒 |
| 尽量多取单次   | 显式传 `limit=5000`（或服务端允许的最大值） |
| 超过单次上限   | 用 `direction=forward` + 上页最后一条时间戳作为新 `start` 分页 |
| 下载为原始日志 | `curl ... | jq -r '.data.result[].values[]? | .[1]' > 节点.log`；要完整 JSON 则 `curl ... -o 文件.json` |

---

## 七、快速参考

| 范围   | start（纳秒，Loki 直连） | start（秒，Exporter） |
|--------|---------------------------|------------------------|
| 1 小时 | `$(($(date -d '1 hour ago' +%s)*1000000000))` | `$(($(date +%s) - 3600))` |
| 6 小时 | `$(($(date -d '6 hours ago' +%s)*1000000000))` | `$(($(date +%s) - 21600))` |
| 24 小时| `$(($(date -d '24 hours ago' +%s)*1000000000))` | `$(($(date +%s) - 86400))` |
| 7 天   | `$(($(date -d '7 days ago' +%s)*1000000000))` | `$(($(date +%s) - 604800))` |

---
