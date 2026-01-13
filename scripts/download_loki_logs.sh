#!/bin/bash
# Loki 日志下载脚本（简单版本）
# 使用 curl 直接调用 Loki API 下载日志

set -e

# 默认配置
LOKI_URL="${LOKI_URL:-http://localhost:3100}"
QUERY="${1:-'{job="injective-node"}'}"
HOURS="${2:-24}"
OUTPUT_FILE="${3:-loki_logs_$(date +%Y%m%d_%H%M%S).json}"

# 计算时间范围
END_TIME=$(date +%s)
START_TIME=$((END_TIME - HOURS * 3600))

# 转换为纳秒
START_NS=$((START_TIME * 1000000000))
END_NS=$((END_TIME * 1000000000))

echo "开始下载日志..."
echo "Loki URL: $LOKI_URL"
echo "查询: $QUERY"
echo "时间范围: $(date -d @$START_TIME) 到 $(date -d @$END_TIME)"
echo "输出文件: $OUTPUT_FILE"
echo "----------------------------------------"

# 查询日志（注意：这个简单版本只获取第一页，最多5000条）
LIMIT=5000
URL="${LOKI_URL}/loki/api/v1/query_range?query=$(echo "$QUERY" | sed 's/ /%20/g')&start=${START_NS}&end=${END_NS}&limit=${LIMIT}"

echo "请求 URL: $URL"
echo ""

# 下载日志
curl -s -G "$URL" \
  --data-urlencode "query=$QUERY" \
  -o "$OUTPUT_FILE"

# 检查响应
if [ $? -eq 0 ]; then
    # 检查是否是有效的 JSON
    if jq empty "$OUTPUT_FILE" 2>/dev/null; then
        # 统计日志条数
        COUNT=$(jq '[.data.result[]?.values[]?] | length' "$OUTPUT_FILE" 2>/dev/null || echo "0")
        echo "下载完成！"
        echo "文件: $OUTPUT_FILE"
        echo "日志条数: $COUNT"
        echo ""
        echo "注意: 此脚本只下载第一页（最多 $LIMIT 条）"
        echo "如需下载完整日志，请使用 download_loki_logs.py 脚本"
    else
        echo "错误: 响应不是有效的 JSON"
        echo "响应内容:"
        cat "$OUTPUT_FILE"
        exit 1
    fi
else
    echo "错误: 下载失败"
    exit 1
fi
