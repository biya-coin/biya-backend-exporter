# WebSocket 实时日志测试指南

本指南说明如何测试 Loki 的 WebSocket 实时日志流功能。

## 问题说明

使用 `curl` 测试 WebSocket 时，需要注意：
1. **curl 不支持 `ws://` 协议**，必须使用 `http://`
2. **必须添加 WebSocket 握手头**：`Sec-WebSocket-Key` 和 `Sec-WebSocket-Version`
3. **查询参数需要 URL 编码**

## 正确的 curl 命令

### 本地测试

```bash
curl --no-buffer \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: $(echo -n 'test' | base64)" \
  -H "Host: localhost:3100" \
  -H "Origin: http://localhost:3100" \
  "http://localhost:3100/loki/api/v1/tail?query=%7Bjob%3D%22injective-node%22%7D&limit=100"
```

### 远程服务器测试

```bash
curl --no-buffer \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: $(echo -n 'test' | base64)" \
  -H "Host: 45.249.245.183:3100" \
  -H "Origin: http://45.249.245.183:3100" \
  "http://45.249.245.183:3100/loki/api/v1/tail?query=%7Bjob%3D%22injective-node%22%7D&limit=100"
```

### URL 编码说明

查询参数 `{job="injective-node"}` 需要编码为 `%7Bjob%3D%22injective-node%22%7D`

可以使用以下命令生成编码：
```bash
echo -n '{job="injective-node"}' | python3 -c "import sys, urllib.parse; print(urllib.parse.quote(sys.stdin.read()))"
# 输出: %7Bjob%3D%22injective-node%22%7D
```

## 使用 wscat（推荐）

`wscat` 是专门用于测试 WebSocket 的工具，使用更简单：

### 安装 wscat

```bash
npm install -g wscat
```

### 使用 wscat 测试

```bash
# 本地测试
wscat -c "ws://localhost:3100/loki/api/v1/tail?query=%7Bjob%3D%22injective-node%22%7D&limit=100"

# 远程服务器测试
wscat -c "ws://45.249.245.183:3100/loki/api/v1/tail?query=%7Bjob%3D%22injective-node%22%7D&limit=100"
```

## 常见错误

### 错误 1: 400 Bad Request

**原因**：
- 使用了 `ws://` 协议（curl 不支持）
- 缺少 `Sec-WebSocket-Key` 或 `Sec-WebSocket-Version` 头
- 查询参数未正确编码

**解决方案**：
- 使用 `http://` 协议
- 添加所有必需的 WebSocket 头
- URL 编码查询参数

### 错误 2: 连接被拒绝

**原因**：
- Loki 服务未运行
- 端口未开放
- 防火墙阻止

**解决方案**：
```bash
# 检查 Loki 是否运行
curl http://localhost:3100/ready

# 检查端口
netstat -tlnp | grep 3100
```

### 错误 3: 查询语法错误

**原因**：
- LogQL 查询语法错误
- 查询参数未正确编码

**解决方案**：
- 验证 LogQL 语法
- 使用 URL 编码工具编码查询参数

## 前端 JavaScript 示例

```javascript
// 正确的 WebSocket 连接（浏览器会自动处理握手）
const ws = new WebSocket(
  'ws://45.249.245.183:3100/loki/api/v1/tail?query=' + 
  encodeURIComponent('{job="injective-node"}') + 
  '&limit=100'
);

ws.onopen = () => {
  console.log('WebSocket connected');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  data.streams?.forEach(stream => {
    stream.values.forEach(([timestamp, logLine]) => {
      console.log(logLine);
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

## 验证 WebSocket 连接

### 方法 1: 使用 wscat（最简单）

```bash
wscat -c "ws://localhost:3100/loki/api/v1/tail?query=%7Bjob%3D%22injective-node%22%7D&limit=10"
```

如果连接成功，会看到实时日志流输出。

### 方法 2: 使用 curl（需要正确配置）

```bash
curl --no-buffer \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: $(openssl rand -base64 16)" \
  -H "Host: localhost:3100" \
  -H "Origin: http://localhost:3100" \
  "http://localhost:3100/loki/api/v1/tail?query=%7Bjob%3D%22injective-node%22%7D&limit=10"
```

如果看到 `HTTP/1.1 101 Switching Protocols`，说明连接成功。

## 总结

- **推荐使用 wscat**：简单、可靠
- **curl 需要特殊配置**：使用 `http://` 协议，添加所有必需的 WebSocket 头
- **查询参数必须 URL 编码**
- **前端直接使用 WebSocket API**：浏览器会自动处理握手
