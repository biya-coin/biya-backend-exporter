#!/bin/bash

# Promtail 部署脚本
# 用法: ./deploy-promtail.sh <日志文件路径> [节点ID]
# 示例: ./deploy-promtail.sh /home/ubuntu/.injectived/logs/inj.log node-1

set -e

# 检查参数
if [ $# -lt 1 ]; then
    echo "用法: $0 <日志文件路径> [节点ID]"
    echo "示例: $0 /home/ubuntu/.injectived/logs/inj.log node-1"
    exit 1
fi

LOG_PATH="$1"
NODE_ID="${2:-$(hostname)}"
LOKI_URL="http://10.8.190.46:3100/loki/api/v1/push"
PROMTAIL_IMAGE="grafana/promtail:3.0.0"
CONFIG_DIR="/opt/promtail"
CONFIG_FILE="${CONFIG_DIR}/promtail-config.yaml"
POSITIONS_DIR="${CONFIG_DIR}/positions"

# 检查日志文件是否存在
if [ ! -f "$LOG_PATH" ]; then
    echo "错误: 日志文件不存在: $LOG_PATH"
    exit 1
fi

# 获取日志文件的目录和文件名
LOG_DIR=$(dirname "$LOG_PATH")
LOG_FILE=$(basename "$LOG_PATH")

echo "=========================================="
echo "Promtail 部署配置"
echo "=========================================="
echo "日志文件路径: $LOG_PATH"
echo "节点ID: $NODE_ID"
echo "Loki 地址: $LOKI_URL"
echo "=========================================="

# 创建配置目录
mkdir -p "$CONFIG_DIR" "$POSITIONS_DIR"

# 生成 Promtail 配置文件
cat > "$CONFIG_FILE" <<EOF
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: ${POSITIONS_DIR}/positions.yaml

clients:
  - url: ${LOKI_URL}

scrape_configs:
  - job_name: validator-node
    static_configs:
      - targets:
          - localhost
        labels:
          job: validator-node
          node_id: ${NODE_ID}
          __path__: ${LOG_PATH}
    
    pipeline_stages:
      # 移除 ANSI 颜色码
      - regex:
          expression: '\x1b\[[0-9;]*m'
          replacement: ''
      
      # 解析日志级别（INF, ERR, WRN, DBG 等）
      - regex:
          expression: '(?P<timestamp>\d{1,2}:\d{2}(?:AM|PM)|\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})\s+(?P<level>INF|ERR|WRN|DBG|FATAL|PANIC)'
          source: log
      
      # 提取级别作为标签
      - match:
          selector: '{job="validator-node"}'
          stages:
            - regex:
                expression: '\s+(?P<level>INF|ERR|WRN|DBG|FATAL|PANIC)\s+'
            - labels:
                level:
      
      # 提取 module 信息（module=xxx）
      - match:
          selector: '{job="validator-node"}'
          stages:
            - regex:
                expression: 'module=(?P<module>\w+)'
            - labels:
                module:
      
      # 标记错误日志
      - match:
          selector: '{job="validator-node"}'
          stages:
            - regex:
                expression: '\s+(ERR|FATAL|PANIC)\s+'
            - labels:
                has_error: "true"
EOF

echo "配置文件已生成: $CONFIG_FILE"

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "错误: 未安装 Docker，请先安装 Docker"
    exit 1
fi

# 停止并删除旧容器（如果存在）
if docker ps -a --format '{{.Names}}' | grep -q "^promtail$"; then
    echo "停止并删除旧的 Promtail 容器..."
    docker stop promtail 2>/dev/null || true
    docker rm promtail 2>/dev/null || true
fi

# 启动 Promtail 容器
echo "启动 Promtail 容器..."
docker run -d \
  --name promtail \
  --restart unless-stopped \
  -p 9080:9080 \
  -v "${CONFIG_FILE}:/etc/promtail/config.yaml:ro" \
  -v "${LOG_DIR}:${LOG_DIR}:ro" \
  -v "${POSITIONS_DIR}:/tmp/positions" \
  "${PROMTAIL_IMAGE}" \
  -config.file=/etc/promtail/config.yaml

# 等待服务启动
echo "等待 Promtail 启动..."
sleep 3

# 检查服务状态
if docker ps --format '{{.Names}}' | grep -q "^promtail$"; then
    echo "✅ Promtail 部署成功！"
    echo ""
    echo "服务信息:"
    echo "  容器名称: promtail"
    echo "  状态检查: docker ps | grep promtail"
    echo "  查看日志: docker logs -f promtail"
    echo "  健康检查: curl http://localhost:9080/ready"
    echo ""
    echo "配置文件位置: $CONFIG_FILE"
    echo "如需修改配置，编辑配置文件后重启容器: docker restart promtail"
else
    echo "❌ Promtail 启动失败，请检查日志:"
    docker logs promtail
    exit 1
fi
