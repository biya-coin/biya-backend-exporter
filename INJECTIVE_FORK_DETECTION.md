# Injective 网络分叉检测方法

## 概述

Injective 网络基于 Tendermint/CometBFT 共识引擎。网络分叉（Fork）是指不同节点对同一高度的区块产生了不同的区块哈希，导致链分裂为多个分支。

## 分叉检测原理

### 1. 基本原理

在 Tendermint 共识中，正常情况下：
- 同一高度的区块应该有**唯一的区块哈希**
- 所有节点应该看到相同的区块哈希
- 如果同一高度出现**多个不同的区块哈希**，说明发生了分叉

### 2. 检测方法

#### 方法一：多源区块哈希对比（推荐）

通过对比不同数据源的同一高度区块哈希来判断分叉：

```
1. 查询 Tendermint RPC 节点的区块哈希
   - 端点：GET /block?height={height}
   - 获取：result.block_id.hash

2. 查询 Explorer API 的区块哈希
   - 端点：GET /api/v1/block/by-height?height={height}
   - 获取：data.block_hash 或 data.hash

3. 对比两个哈希值
   - 如果哈希值相同 → 无分叉
   - 如果哈希值不同 → 检测到分叉
```

#### 方法二：多节点对比

如果有多个 Tendermint RPC 节点：

```
1. 查询节点A的区块哈希（高度H）
2. 查询节点B的区块哈希（高度H）
3. 查询节点C的区块哈希（高度H）
4. 如果存在不同的哈希值 → 检测到分叉
```

#### 方法三：分叉深度计算

如果检测到分叉，需要计算分叉深度：

```
1. 从当前高度开始，向前回溯
2. 比较每个高度的区块哈希
3. 找到最后一个所有节点都一致的区块高度
4. 分叉深度 = 当前高度 - 最后一个共同区块高度
```

## 实现方案

### 当前项目状态

在 `configs/prometheus/alert_rules.yml` 中已有分叉告警规则：

```yaml
- alert: 网络分叉
  expr: |
    (biya_chain_fork_depth or vector(0)) > 3
  for: 1m
  labels:
    severity: critical
```

但 `biya_chain_fork_depth` 指标尚未实现。

### 推荐实现步骤

#### 1. 创建分叉检测 Collector

在 `internal/collectors/` 下创建 `fork_detector.go`：

```go
package collectors

import (
    "context"
    "log/slog"
    "strconv"
    
    "github.com/biya-coin/biya-dex-backend-exporter/internal/adapters/tendermint"
    "github.com/biya-coin/biya-dex-backend-exporter/internal/adapters/explorer"
    "github.com/biya-coin/biya-dex-backend-exporter/internal/metrics"
)

type ForkDetector struct {
    log     *slog.Logger
    m       *metrics.Metrics
    tm      *tendermint.Client
    explorer *explorer.Client
}

func (d *ForkDetector) DetectFork(ctx context.Context) (int64, error) {
    // 1. 获取当前最新高度
    st, err := d.tm.Status(ctx)
    if err != nil {
        return 0, err
    }
    
    height, _ := strconv.ParseInt(st.Result.SyncInfo.LatestBlockHeight, 10, 64)
    
    // 2. 查询 Tendermint 节点的区块哈希
    tmBlock, err := d.tm.Block(ctx, height)
    if err != nil {
        return 0, err
    }
    tmHash := tmBlock.Result.BlockID.Hash
    
    // 3. 查询 Explorer 的区块哈希
    explorerBlock, err := d.explorer.GetBlockByHeight(ctx, strconv.FormatInt(height, 10))
    if err != nil {
        // Explorer 可能暂时不可用，不视为分叉
        return 0, nil
    }
    
    // 4. 解析 Explorer 返回的区块哈希
    // 需要根据实际 API 响应格式解析
    explorerHash := parseBlockHash(explorerBlock)
    
    // 5. 比较哈希值
    if explorerHash != "" && tmHash != explorerHash {
        // 检测到分叉，计算分叉深度
        return d.calculateForkDepth(ctx, height, tmHash, explorerHash)
    }
    
    return 0, nil // 无分叉
}

func (d *ForkDetector) calculateForkDepth(ctx context.Context, startHeight int64, hash1, hash2 string) (int64, error) {
    // 向前回溯，找到最后一个共同区块
    maxDepth := int64(100) // 最多回溯100个区块
    for i := int64(1); i <= maxDepth && startHeight-i > 0; i++ {
        height := startHeight - i
        
        tmBlock, err := d.tm.Block(ctx, height)
        if err != nil {
            continue
        }
        
        explorerBlock, err := d.explorer.GetBlockByHeight(ctx, strconv.FormatInt(height, 10))
        if err != nil {
            continue
        }
        
        tmHash := tmBlock.Result.BlockID.Hash
        explorerHash := parseBlockHash(explorerBlock)
        
        if explorerHash != "" && tmHash == explorerHash {
            // 找到共同区块
            return i, nil
        }
    }
    
    return maxDepth, nil // 分叉深度超过最大回溯范围
}
```

#### 2. 注册指标

在 `internal/metrics/metrics.go` 中添加：

```go
reg.MustDeclare("biya_chain_fork_depth", TypeGauge, "Chain fork depth (number of blocks since last common block).", []string{"chain_id"})
```

#### 3. 定期执行检测

在 `internal/collectors/scheduler.go` 中调度分叉检测：

```go
// 每30秒检测一次分叉
scheduler.Every(30 * time.Second).Do(func() {
    forkDepth, err := forkDetector.DetectFork(ctx)
    if err != nil {
        log.Warn("fork detection failed", "err", err)
        return
    }
    metrics.SetGauge("biya_chain_fork_depth", map[string]string{"chain_id": chainID}, float64(forkDepth))
})
```

## 检测策略

### 检测频率

- **实时检测**：每 30 秒检测一次最新区块
- **深度回溯**：仅在检测到分叉时执行（避免频繁查询）

### 容错处理

1. **Explorer 不可用**：不视为分叉，仅记录警告
2. **节点暂时离线**：跳过该节点，使用其他节点对比
3. **网络延迟**：设置合理的超时时间（建议 5-10 秒）

### 告警阈值

根据需求文档：
- **分叉深度 > 3 个区块**：触发严重告警
- **分叉深度 1-3 个区块**：记录警告日志

## Tendermint/CometBFT 原生指标

CometBFT 本身也提供了一些共识相关的指标，可用于辅助判断：

```
cometbft_consensus_byzantine_validators          # 拜占庭验证者数量
cometbft_consensus_missing_validators            # 缺失签名的验证者数量
cometbft_consensus_round_duration_seconds        # 共识轮次持续时间
```

这些指标异常可能预示着潜在的分叉风险。

## 注意事项

1. **正常网络延迟**：不同节点可能在不同时间看到新区块，需要区分正常延迟和真实分叉
2. **计划内分叉**：网络升级时的硬分叉是正常的，需要结合官方公告判断
3. **临时分叉**：Tendermint 的最终确定性机制会快速解决临时分叉（通常 < 1 分钟）
4. **数据源可靠性**：确保对比的数据源是独立且可靠的

## 参考资源

- [Tendermint 共识文档](https://docs.tendermint.com/v0.37/core/consensus.html)
- [CometBFT 指标文档](https://docs.cometbft.com/v0.38/core/metrics)
- [Injective 网络文档](https://docs.injective.network/)
