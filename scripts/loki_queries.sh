#!/bin/bash

# Loki 常用查询语句脚本
# 用途：提供常用的 LogQL 查询语句作为参考
# 使用方法：可以直接复制查询语句到 Grafana Explore 或通过 API 调用

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Loki API 地址（可通过环境变量覆盖）
LOKI_URL="${LOKI_URL:-http://localhost:3100}"
# export LOKI_URL=http://localhost:3100

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Loki 常用查询语句参考脚本${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# ============================================
# 1. 基础查询
# ============================================
echo -e "${GREEN}1. 基础查询${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 查询所有日志
{job="injective-node"}

# 查询特定节点日志
{job="injective-node", node_id="node-1"}

# 多标签查询
{job="injective-node", level="error", node_id="node-1"}

# 使用正则表达式匹配标签值
{job="injective-node", node_id=~"node-.*"}

# 排除特定标签值
{job="injective-node", node_id!="node-1"}

# 多个标签值匹配（OR）
{job="injective-node", level=~"error|warn"}
EOF
echo ""

# ============================================
# 2. 文本过滤查询
# ============================================
echo -e "${GREEN}2. 文本过滤查询${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 包含特定字符串
{job="injective-node"} |= "error"

# 不包含特定字符串
{job="injective-node"} != "debug"

# 正则表达式匹配
{job="injective-node"} |~ "error|exception|fatal"

# 不匹配正则表达式
{job="injective-node"} !~ "debug|trace"

# 大小写敏感匹配
{job="injective-node"} |= "Error"

# 大小写不敏感匹配（使用正则）
{job="injective-node"} |~ "(?i)error"

# 多条件组合（AND）
{job="injective-node"} |= "error" |= "timeout"

# 多条件组合（OR）
{job="injective-node"} |~ "error|timeout|failed"
EOF
echo ""

# ============================================
# 3. 时间范围查询
# ============================================
echo -e "${GREEN}3. 时间范围查询${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 最近5分钟
{job="injective-node"} [5m]

# 最近1小时
{job="injective-node"} [1h]

# 最近24小时
{job="injective-node"} [24h]

# 最近7天
{job="injective-node"} [7d]

# 指定时间范围（需要配合 API 使用）
# start 和 end 参数使用纳秒时间戳
EOF
echo ""

# ============================================
# 4. 聚合查询 - 计数
# ============================================
echo -e "${GREEN}4. 聚合查询 - 计数${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 统计日志条数（最近1小时）
count_over_time({job="injective-node"}[1h])

# 按级别统计日志数量
sum by (level) (count_over_time({job="injective-node"}[1h]))

# 按节点统计日志数量
sum by (node_id) (count_over_time({job="injective-node"}[1h]))

# 多维度统计
sum by (level, node_id) (count_over_time({job="injective-node"}[1h]))

# 统计错误日志数量
sum(count_over_time({job="injective-node", level="error"}[1h]))

# 统计包含特定关键词的日志数量
sum(count_over_time({job="injective-node"} |= "error" [1h]))
EOF
echo ""

# ============================================
# 5. 聚合查询 - 速率
# ============================================
echo -e "${GREEN}5. 聚合查询 - 速率${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 每秒日志速率
rate({job="injective-node"}[5m])

# 按级别统计速率
sum by (level) (rate({job="injective-node"}[5m]))

# 按节点统计速率
sum by (node_id) (rate({job="injective-node"}[5m]))

# 错误日志速率
sum(rate({job="injective-node", level="error"}[5m]))

# 错误率百分比
sum(rate({job="injective-node", level="error"}[5m])) / sum(rate({job="injective-node"}[5m])) * 100
EOF
echo ""

# ============================================
# 6. 聚合查询 - 其他函数
# ============================================
echo -e "${GREEN}6. 聚合查询 - 其他函数${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 计算平均值（需要数值字段）
avg_over_time({job="injective-node"} | json | __error__="" [5m])

# 计算最大值
max_over_time({job="injective-node"} | json | latency [5m])

# 计算最小值
min_over_time({job="injective-node"} | json | latency [5m])

# 计算总和
sum_over_time({job="injective-node"} | json | count [5m])

# 计算分位数（P95）
quantile_over_time(0.95, {job="injective-node"} | json | latency [5m])
EOF
echo ""

# ============================================
# 7. JSON 解析查询
# ============================================
echo -e "${GREEN}7. JSON 解析查询${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 解析 JSON 日志
{job="injective-node"} | json

# 解析并过滤 JSON 字段
{job="injective-node"} | json | level="error"

# 解析并提取特定字段
{job="injective-node"} | json | latency > 1000

# 解析并重命名字段
{job="injective-node"} | json | line_format "{{.timestamp}} {{.message}}"

# 解析并过滤多个字段
{job="injective-node"} | json | level="error" | latency > 1000
EOF
echo ""

# ============================================
# 8. 日志格式化和提取
# ============================================
echo -e "${GREEN}8. 日志格式化和提取${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 使用正则表达式提取字段
{job="injective-node"} | regexp "(?P<ip>\\d+\\.\\d+\\.\\d+\\.\\d+)"

# 使用正则表达式提取并过滤
{job="injective-node"} | regexp "(?P<status>\\d{3})" | status="500"

# 格式化输出
{job="injective-node"} | line_format "{{.timestamp}} [{{.level}}] {{.message}}"

# 标签格式化
{job="injective-node"} | label_format node="{{.node_id}}"
EOF
echo ""

# ============================================
# 9. 高级查询 - 子查询
# ============================================
echo -e "${GREEN}9. 高级查询 - 子查询${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 计算错误率（使用子查询）
sum(rate({job="injective-node", level="error"}[5m])) / sum(rate({job="injective-node"}[5m]))

# 计算增长率
(sum(rate({job="injective-node"}[5m])) - sum(rate({job="injective-node"}[5m] offset 1h))) / sum(rate({job="injective-node"}[5m] offset 1h)) * 100

# 计算同比（与昨天同时段对比）
sum(rate({job="injective-node"}[5m])) / sum(rate({job="injective-node"}[5m] offset 24h))
EOF
echo ""

# ============================================
# 10. 实际应用场景查询
# ============================================
echo -e "${GREEN}10. 实际应用场景查询${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 查询最近的错误日志（最近1小时，限制100条）
{job="injective-node", level="error"} [1h] limit 100

# 查询特定节点的警告和错误
{job="injective-node", node_id="node-1", level=~"warn|error"}

# 查询包含超时的日志
{job="injective-node"} |~ "timeout|deadline exceeded"

# 查询交易相关错误
{job="injective-node"} |= "transaction" |= "error"

# 查询区块同步相关日志
{job="injective-node"} |~ "block|sync|height"

# 查询性能问题（延迟超过阈值）
{job="injective-node"} | json | latency > 5000

# 查询特定时间段的错误统计
sum by (node_id) (count_over_time({job="injective-node", level="error"}[1h]))

# 查询错误趋势（按小时）
sum by (hour) (count_over_time({job="injective-node", level="error"}[1h]))
EOF
echo ""

# ============================================
# 11. 直接查询 Loki API (curl 命令)
# ============================================
echo -e "${GREEN}11. 直接查询 Loki API (curl 命令)${NC}"
echo "----------------------------------------"
echo -e "${YELLOW}说明: LogQL 是 Loki 的查询语言，可以直接通过 HTTP API 使用${NC}"
echo ""

# 计算时间戳的辅助函数
calc_timestamp() {
    local hours="${1:-1}"
    date -d "${hours} hour ago" +%s 2>/dev/null || date -v-${hours}H +%s 2>/dev/null || echo "$(($(date +%s) - ${hours} * 3600))"
}

# 生成纳秒时间戳
START_1H=$(($(calc_timestamp 1) * 1000000000))
END_NOW=$(($(date +%s) * 1000000000))
START_24H=$(($(calc_timestamp 24) * 1000000000))

cat << EOF
# ===== 基础查询 =====

# 1. 查询最近1小时的所有日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"}' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# 2. 查询最近1小时的错误日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node", level="error"}' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# 3. 查询特定节点的日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node", node_id="node-1"}' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# 4. 查询包含特定关键词的日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"} |= "error"' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# 5. 查询正则表达式匹配的日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"} |~ "timeout|error|failed"' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# ===== 即时查询 (当前时刻的快照) =====

# 6. 即时查询当前日志
curl -s -G "${LOKI_URL}/loki/api/v1/query" \\
  --data-urlencode 'query={job="injective-node"}' \\
  --data-urlencode "limit=100" | jq '.'

# 7. 即时查询错误日志
curl -s -G "${LOKI_URL}/loki/api/v1/query" \\
  --data-urlencode 'query={job="injective-node", level="error"}' \\
  --data-urlencode "limit=100" | jq '.'

# ===== 聚合查询 =====

# 8. 统计最近1小时的日志数量
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query=count_over_time({job="injective-node"}[1h])' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# 9. 按级别统计日志数量
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query=sum by (level) (count_over_time({job="injective-node"}[1h]))' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# 10. 统计错误日志数量
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query=sum(count_over_time({job="injective-node", level="error"}[1h]))' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# 11. 计算日志速率（每秒）
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query=sum(rate({job="injective-node"}[5m]))' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# 12. 按节点统计错误数量
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query=sum by (node_id) (count_over_time({job="injective-node", level="error"}[1h]))' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# ===== 标签查询 =====

# 13. 获取所有标签（可指定时间范围）
curl -s -G "${LOKI_URL}/loki/api/v1/labels" \\
  --data-urlencode "start=${START_24H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# 14. 获取特定标签的所有值（需要指定时间范围）
curl -s -G "${LOKI_URL}/loki/api/v1/label/node_id/values" \\
  --data-urlencode "start=${START_24H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# 15. 获取 level 标签的所有值
curl -s -G "${LOKI_URL}/loki/api/v1/label/level/values" \\
  --data-urlencode "start=${START_24H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# 16. 获取 job 标签的所有值
curl -s -G "${LOKI_URL}/loki/api/v1/label/job/values" \\
  --data-urlencode "start=${START_24H}" \\
  --data-urlencode "end=${END_NOW}" | jq '.'

# ===== 时间范围查询 =====

# 16. 查询最近24小时的日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"}' \\
  --data-urlencode "start=${START_24H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=1000" | jq '.'

# 17. 查询最近5分钟的日志（使用相对时间）
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"}[5m]' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# ===== 实用查询组合 =====

# 18. 查询交易相关错误
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"} |= "transaction" |= "error"' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# 19. 查询区块同步相关日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"} |~ "block|sync|height"' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# 20. 查询警告和错误级别日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node", level=~"warn|error"}' \\
  --data-urlencode "start=${START_1H}" \\
  --data-urlencode "end=${END_NOW}" \\
  --data-urlencode "limit=100" | jq '.'

# ===== 时间戳转换说明 =====
# 
# 时间戳需要转换为纳秒（乘以 1000000000）
# 
# 获取1小时前的时间戳（秒）:
#   date -d '1 hour ago' +%s
# 
# 转换为纳秒:
#   echo \$(date -d '1 hour ago' +%s)000000000
# 
# 获取当前时间戳（纳秒）:
#   echo \$(date +%s)000000000
# 
# 示例：查询最近1小时的完整命令
#   START=\$(date -d '1 hour ago' +%s)000000000
#   END=\$(date +%s)000000000
#   curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
#     --data-urlencode 'query={job="injective-node"}' \\
#     --data-urlencode "start=\${START}" \\
#     --data-urlencode "end=\${END}" \\
#     --data-urlencode "limit=100" | jq '.'
EOF
echo ""

# ============================================
# 12. API 调用示例函数
# ============================================
echo -e "${GREEN}12. API 调用示例函数（可直接使用的函数）${NC}"
echo "----------------------------------------"

# 函数：查询范围日志
query_range() {
    local query="$1"
    local start="${2:-$(date -d '1 hour ago' +%s)000000000}"
    local end="${3:-$(date +%s)000000000}"
    local limit="${4:-100}"
    
    echo "查询: $query"
    echo "时间范围: $start 到 $end"
    echo ""
    
    curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
        --data-urlencode "query=${query}" \
        --data-urlencode "start=${start}" \
        --data-urlencode "end=${end}" \
        --data-urlencode "limit=${limit}" | jq '.'
}

# 函数：即时查询
query_instant() {
    local query="$1"
    local limit="${2:-100}"
    
    echo "查询: $query"
    echo ""
    
    curl -s -G "${LOKI_URL}/loki/api/v1/query" \
        --data-urlencode "query=${query}" \
        --data-urlencode "limit=${limit}" | jq '.'
}

# 函数：获取标签列表
get_labels() {
    local hours="${1:-24}"
    local start=$(($(date -d "${hours} hour ago" +%s 2>/dev/null || date -v-${hours}H +%s 2>/dev/null || echo "$(($(date +%s) - ${hours} * 3600))") * 1000000000))
    local end=$(($(date +%s) * 1000000000))
    echo "获取最近 ${hours} 小时的所有标签..."
    curl -s -G "${LOKI_URL}/loki/api/v1/labels" \
        --data-urlencode "start=${start}" \
        --data-urlencode "end=${end}" | jq '.'
}

# 函数：获取标签值
get_label_values() {
    local label_name="$1"
    local hours="${2:-24}"
    local start=$(($(date -d "${hours} hour ago" +%s 2>/dev/null || date -v-${hours}H +%s 2>/dev/null || echo "$(($(date +%s) - ${hours} * 3600))") * 1000000000))
    local end=$(($(date +%s) * 1000000000))
    echo "获取标签 '${label_name}' 最近 ${hours} 小时的所有值..."
    curl -s -G "${LOKI_URL}/loki/api/v1/label/${label_name}/values" \
        --data-urlencode "start=${start}" \
        --data-urlencode "end=${end}" | jq '.'
}

# 函数：健康检查
health_check() {
    echo "检查 Loki 服务状态..."
    curl -s "${LOKI_URL}/ready" && echo -e "\n${GREEN}✓ Loki 服务正常${NC}" || echo -e "\n${RED}✗ Loki 服务异常${NC}"
}

# 函数：查询并格式化输出日志内容（只显示日志文本，不显示 JSON）
query_logs() {
    local query="$1"
    local start="${2:-$(date -d '1 hour ago' +%s)000000000}"
    local end="${3:-$(date +%s)000000000}"
    local limit="${4:-100}"
    
    echo -e "${BLUE}查询: ${query}${NC}"
    echo -e "${YELLOW}时间范围: $(date -d @$(($start / 1000000000))) 到 $(date -d @$(($end / 1000000000)))${NC}"
    echo ""
    
    curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
        --data-urlencode "query=${query}" \
        --data-urlencode "start=${start}" \
        --data-urlencode "end=${end}" \
        --data-urlencode "limit=${limit}" | \
    jq -r '.data.result[] | .values[] | "\(.[0] | tonumber / 1000000 | strftime("%Y-%m-%d %H:%M:%S")) | \(.[1])"'
}

# 函数：查询错误日志（快捷函数）
query_errors() {
    local hours="${1:-1}"
    local limit="${2:-100}"
    local start=$(($(date -d "${hours} hour ago" +%s 2>/dev/null || date -v-${hours}H +%s 2>/dev/null || echo "$(($(date +%s) - ${hours} * 3600))") * 1000000000))
    local end=$(($(date +%s) * 1000000000))
    
    query_logs '{job="injective-node", level="error"}' "${start}" "${end}" "${limit}"
}

# 函数：查询特定节点的日志（快捷函数）
query_node() {
    local node_id="$1"
    local hours="${2:-1}"
    local limit="${3:-100}"
    local start=$(($(date -d "${hours} hour ago" +%s 2>/dev/null || date -v-${hours}H +%s 2>/dev/null || echo "$(($(date +%s) - ${hours} * 3600))") * 1000000000))
    local end=$(($(date +%s) * 1000000000))
    
    query_logs "{job=\"injective-node\", node_id=\"${node_id}\"}" "${start}" "${end}" "${limit}"
}

# 函数：统计错误数量（快捷函数）
count_errors() {
    local hours="${1:-1}"
    local start=$(($(date -d "${hours} hour ago" +%s 2>/dev/null || date -v-${hours}H +%s 2>/dev/null || echo "$(($(date +%s) - ${hours} * 3600))") * 1000000000))
    local end=$(($(date +%s) * 1000000000))
    
    local result=$(curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \
        --data-urlencode 'query=sum(count_over_time({job="injective-node", level="error"}['${hours}'h]))' \
        --data-urlencode "start=${start}" \
        --data-urlencode "end=${end}" | jq -r '.data.result[0].value[1] // "0"')
    
    echo -e "${GREEN}最近 ${hours} 小时的错误日志数量: ${result}${NC}"
}

# ============================================
# 13. 使用示例（函数调用）
# ============================================
echo -e "${GREEN}13. 使用示例（函数调用）${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 注意：以下函数需要先 source 脚本才能使用
# source scripts/loki_queries.sh

# ===== 基础查询函数 =====

# 示例1: 查询最近1小时的错误日志（返回 JSON）
query_range '{job="injective-node", level="error"}'

# 示例2: 指定时间范围查询
query_range '{job="injective-node"}' "$(date -d '2 hours ago' +%s)000000000" "$(date +%s)000000000" 200

# 示例3: 即时查询
query_instant '{job="injective-node"} |= "error"'

# 示例4: 查询并格式化输出日志（只显示日志内容，不显示 JSON）
query_logs '{job="injective-node", level="error"}'

# ===== 快捷查询函数 =====

# 示例5: 查询最近1小时的错误日志（格式化输出）
query_errors 1

# 示例6: 查询最近24小时的错误日志
query_errors 24 500

# 示例7: 查询特定节点的日志
query_node "node-1" 1

# 示例8: 查询特定节点最近2小时的日志
query_node "node-1" 2 200

# 示例9: 统计错误数量
count_errors 1    # 最近1小时
count_errors 24   # 最近24小时

# ===== 标签和状态函数 =====

# 示例10: 获取所有标签（默认最近24小时）
get_labels

# 示例10a: 获取最近7天的所有标签
get_labels 168

# 示例11: 获取 node_id 标签的所有值（默认最近24小时）
get_label_values "node_id"

# 示例11a: 获取 node_id 标签最近7天的所有值
get_label_values "node_id" 168

# 示例12: 获取 level 标签的所有值
get_label_values "level"

# 示例13: 健康检查
health_check
EOF
echo ""

# ============================================
# 14. 常用查询快捷命令（LogQL 语句）
# ============================================
echo -e "${GREEN}14. 常用查询快捷命令（LogQL 语句）${NC}"
echo "----------------------------------------"
cat << 'EOF'
# 查询最近错误（最近1小时）
{job="injective-node", level="error"} [1h]

# 查询特定节点最近日志
{job="injective-node", node_id="node-1"} [5m]

# 统计错误数量（最近1小时）
sum(count_over_time({job="injective-node", level="error"}[1h]))

# 错误速率（每秒）
sum(rate({job="injective-node", level="error"}[5m]))

# 按节点统计错误数量
sum by (node_id) (count_over_time({job="injective-node", level="error"}[1h]))

# 查询包含超时的日志
{job="injective-node"} |~ "timeout|deadline"

# 查询交易相关错误
{job="injective-node"} |= "transaction" |= "error"

# 错误率百分比
sum(rate({job="injective-node", level="error"}[5m])) / sum(rate({job="injective-node"}[5m])) * 100
EOF
echo ""

# ============================================
# 15. 快速参考：直接查询命令
# ============================================
echo -e "${GREEN}15. 快速参考：直接查询命令${NC}"
echo "----------------------------------------"
cat << EOF
# 快速查询最近1小时的错误日志（一行命令）
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node", level="error"}' \\
  --data-urlencode "start=\$(date -d '1 hour ago' +%s)000000000" \\
  --data-urlencode "end=\$(date +%s)000000000" \\
  --data-urlencode "limit=100" | jq '.data.result[] | .values[] | .[1]'

# 快速查询包含关键词的日志
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query={job="injective-node"} |= "error"' \\
  --data-urlencode "start=\$(date -d '1 hour ago' +%s)000000000" \\
  --data-urlencode "end=\$(date +%s)000000000" \\
  --data-urlencode "limit=50" | jq '.data.result[] | .values[] | .[1]'

# 快速统计错误数量
curl -s -G "${LOKI_URL}/loki/api/v1/query_range" \\
  --data-urlencode 'query=sum(count_over_time({job="injective-node", level="error"}[1h]))' \\
  --data-urlencode "start=\$(date -d '1 hour ago' +%s)000000000" \\
  --data-urlencode "end=\$(date +%s)000000000" | jq '.data.result[0].value[1]'

# 快速获取所有标签（最近24小时）
curl -s -G "${LOKI_URL}/loki/api/v1/labels" \\
  --data-urlencode "start=\$(date -d '24 hours ago' +%s)000000000" \\
  --data-urlencode "end=\$(date +%s)000000000" | jq '.data[]'

# 快速获取节点列表（最近24小时）
curl -s -G "${LOKI_URL}/loki/api/v1/label/node_id/values" \\
  --data-urlencode "start=\$(date -d '24 hours ago' +%s)000000000" \\
  --data-urlencode "end=\$(date +%s)000000000" | jq '.data[]'
EOF
echo ""

# ============================================
# 16. 重要说明
# ============================================
echo -e "${GREEN}16. 重要说明${NC}"
echo "----------------------------------------"
cat << 'EOF'
1. LogQL 是 Loki 的查询语言，不是 Grafana 专用的
   - LogQL 可以直接通过 Loki HTTP API 使用
   - Grafana 只是提供了一个可视化界面来执行 LogQL 查询

2. 直接查询 Loki 的方式：
   - 使用 curl 命令调用 HTTP API（见第11章节）
   - 使用脚本中的函数（需要 source 脚本）
   - 使用任何支持 HTTP 的客户端工具

3. API 端点：
   - 范围查询: GET /loki/api/v1/query_range
   - 即时查询: GET /loki/api/v1/query
   - 标签列表: GET /loki/api/v1/labels
   - 标签值:   GET /loki/api/v1/label/<name>/values
   - 健康检查: GET /ready

4. 时间戳格式：
   - Loki API 需要纳秒级时间戳（Unix 时间戳 * 1000000000）
   - 示例: $(date +%s)000000000

5. 查询参数：
   - query: LogQL 查询语句（必需）
   - start: 开始时间（纳秒时间戳，范围查询必需）
   - end:   结束时间（纳秒时间戳，范围查询必需）
   - limit: 返回结果数量限制（默认 100）
EOF
echo ""

# ============================================
# 17. 导出函数供外部使用
# ============================================
echo -e "${GREEN}17. 导出函数供外部使用${NC}"
echo "----------------------------------------"
echo -e "${YELLOW}提示: 要使用 API 调用函数，请运行:${NC}"
echo -e "${YELLOW}  source $0${NC}"
echo -e "${YELLOW}然后就可以直接调用函数，例如:${NC}"
echo -e "${YELLOW}  health_check${NC}"
echo -e "${YELLOW}  query_range '{job=\"injective-node\", level=\"error\"}'${NC}"
echo -e "${YELLOW}  get_labels${NC}"
echo -e "${YELLOW}  get_label_values \"node_id\"${NC}"
echo ""

# 如果脚本被 source 执行，导出函数
if [[ "${BASH_SOURCE[0]}" != "${0}" ]]; then
    export -f query_range
    export -f query_instant
    export -f query_logs
    export -f query_errors
    export -f query_node
    export -f count_errors
    export -f get_labels
    export -f get_label_values
    export -f health_check
    echo -e "${GREEN}✓ 函数已导出，可以直接使用${NC}"
    echo ""
    echo -e "${BLUE}可用函数:${NC}"
    echo "  ${GREEN}基础查询函数:${NC}"
    echo "    - query_range <query> [start] [end] [limit]  # 范围查询（返回 JSON）"
    echo "    - query_instant <query> [limit]              # 即时查询（返回 JSON）"
    echo "    - query_logs <query> [start] [end] [limit]   # 范围查询（格式化输出日志）"
    echo ""
    echo "  ${GREEN}快捷查询函数:${NC}"
    echo "    - query_errors [hours] [limit]               # 查询错误日志（默认最近1小时）"
    echo "    - query_node <node_id> [hours] [limit]      # 查询特定节点日志"
    echo "    - count_errors [hours]                       # 统计错误数量"
    echo ""
    echo "  ${GREEN}标签和状态函数:${NC}"
    echo "    - get_labels [hours]                          # 获取所有标签（默认最近24小时）"
    echo "    - get_label_values <label_name> [hours]       # 获取标签值（默认最近24小时）"
    echo "    - health_check                               # 健康检查"
    echo ""
    echo -e "${YELLOW}使用示例:${NC}"
    echo "  query_errors 1        # 查询最近1小时的错误日志"
    echo "  query_node node-1 2   # 查询 node-1 最近2小时的日志"
    echo "  count_errors 24       # 统计最近24小时的错误数量"
    echo "  get_labels 168        # 获取最近7天的所有标签"
    echo "  get_label_values node_id 48  # 获取 node_id 最近48小时的所有值"
fi
