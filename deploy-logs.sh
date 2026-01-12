#!/bin/bash
# 部署 Loki 和 Promtail 日志收集服务
# 用于收集 ~/.injectived/logs/inj.log 日志

set -e

echo "=== 部署 Loki 和 Promtail 日志收集服务 ==="

# 检查网络是否存在
if ! docker network ls | grep -q "biya-backend"; then
    echo "错误: biya-backend 网络不存在，请先创建网络："
    echo "  docker network create biya-backend"
    exit 1
fi

# 检查日志文件是否存在
if [ ! -f ~/.injectived/logs/inj.log ]; then
    echo "警告: 日志文件 ~/.injectived/logs/inj.log 不存在"
    echo "请确保 Injective 节点正在运行并生成日志"
fi

# 启动 Loki
echo ""
echo "1. 启动 Loki 服务..."
docker compose -f compose.log.yaml up -d loki

# 等待 Loki 就绪
echo "   等待 Loki 就绪..."
for i in {1..30}; do
    if curl -s http://localhost:3100/ready > /dev/null 2>&1; then
        echo "   ✓ Loki 已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "   ✗ Loki 启动超时"
        docker compose -f compose.log.yaml logs loki
        exit 1
    fi
    sleep 1
done

# 启动 Promtail
echo ""
echo "2. 启动 Promtail 服务..."
docker compose -f compose.log.yaml up -d promtail

# 等待 Promtail 就绪
echo "   等待 Promtail 就绪..."
for i in {1..30}; do
    if docker exec promtail wget -qO- http://localhost:9080/ready > /dev/null 2>&1; then
        echo "   ✓ Promtail 已就绪"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "   ✗ Promtail 启动超时"
        docker compose -f compose.log.yaml logs promtail
        exit 1
    fi
    sleep 1
done

# 验证服务状态
echo ""
echo "=== 服务状态 ==="
docker compose -f compose.log.yaml ps

echo ""
echo "=== 验证日志收集 ==="
echo "等待 5 秒让 Promtail 开始收集日志..."
sleep 5

# 查询日志验证
echo ""
echo "查询最近的日志（前10条）..."
curl -s -G "http://localhost:3100/loki/api/v1/query_range" \
  --data-urlencode 'query={job="injective-node"}' \
  --data-urlencode "start=$(date -d '5 minutes ago' +%s)000000000" \
  --data-urlencode "end=$(date +%s)000000000" \
  --data-urlencode "limit=10" | jq -r '.data.result[0].values[-10:] | .[] | .[1]' 2>/dev/null || echo "提示: 如果查询结果为空，可能是日志还未被收集，请稍后重试"

echo ""
echo "=== 部署完成 ==="
echo ""
echo "服务端点："
echo "  - Loki API: http://localhost:3100"
echo "  - Loki WebSocket: ws://localhost:3100/loki/api/v1/tail"
echo "  - Promtail: http://localhost:9080"
echo ""
echo "查看日志："
echo "  docker compose -f compose.log.yaml logs -f promtail"
echo ""
echo "查询日志示例："
echo '  curl -s -G "http://localhost:3100/loki/api/v1/query_range" \'
echo '    --data-urlencode '\''query={job="injective-node"}'\'' \'
echo '    --data-urlencode "start=$(date -d '\''1 hour ago'\'' +%s)000000000" \'
echo '    --data-urlencode "end=$(date +%s)000000000" \'
echo '    --data-urlencode "limit=100"'
echo ""
echo "在 Grafana 中查看："
echo "  1. 打开 http://localhost:3000"
echo "  2. 进入 Explore"
echo "  3. 选择 Loki 数据源"
echo "  4. 输入查询: {job=\"injective-node\"}"
echo "  5. 点击 Live 按钮启用实时日志流"
