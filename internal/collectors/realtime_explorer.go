package collectors

import (
	"context"
	"log/slog"
	"strings"

	"github.com/biya-coin/biya-dex-backend-exporter/internal/adapters/explorer"
	"github.com/biya-coin/biya-dex-backend-exporter/internal/adapters/tendermint"
	"github.com/biya-coin/biya-dex-backend-exporter/internal/config"
	"github.com/biya-coin/biya-dex-backend-exporter/internal/metrics"
)

// RealtimeExplorerCollector 负责用 biya-explorer API 填充 METRICS.md 中的 explorer 指标。
// 拿不到的数据（接口缺失/字段不明确/APIKey 未配置）会按 mock 配置兜底为固定值，确保 exporter 可用。
type RealtimeExplorerCollector struct {
	log  *slog.Logger
	m    *metrics.Metrics
	api  *explorer.Client
	tm   *tendermint.Client
	mock config.MockConfig
}

func NewRealtimeExplorerCollector(log *slog.Logger, m *metrics.Metrics, api *explorer.Client, tm *tendermint.Client, mock config.MockConfig) *RealtimeExplorerCollector {
	return &RealtimeExplorerCollector{log: log, m: m, api: api, tm: tm, mock: mock}
}

func (c *RealtimeExplorerCollector) Run(ctx context.Context) error {
	// provide.md 指标口径：
	// - block height:   GET /api/v1/block/latest                -> .data.data[0].height
	// - tx stats:       GET /api/v1/transaction/stats           -> .data.count_24h / .data.tps / .data.avg_block_time / .data.active_addresses_24h
	// - gas price gwei: GET /api/v1/block/gas-utilization       -> .data.gas_price（你已澄清：该字段即“平均 gas 费”）

	if v, ok := c.readLatestBlockHeight(ctx); ok {
		c.m.SetGauge("biya_block_height", nil, v)
	}

	if v, ok := c.readBlockInterval(ctx); ok {
		c.m.SetGauge("biya_block_interval", nil, v)
	}

	if stats, ok := c.readTransactionStats(ctx); ok {
		if stats.Count24H >= 0 {
			c.m.SetGauge("biya_tx_24h_total", nil, stats.Count24H)
		}
		if stats.TPS >= 0 {
			c.m.SetGauge("biya_tps_current", nil, stats.TPS)
		}
		if stats.AvgBlockTimeSeconds >= 0 {
			c.m.SetGauge("biya_avg_block_time", nil, stats.AvgBlockTimeSeconds)
		}
		if stats.ActiveAddresses24H >= 0 {
			c.m.SetGauge("biya_active_addresses_24h", nil, stats.ActiveAddresses24H)
		}
	}

	if v, ok := c.readGasPriceGwei(ctx); ok {
		c.m.SetGauge("biya_gas_price", nil, v)
	}

	if v, ok := c.readGasUtilization(ctx); ok {
		c.m.SetGauge("biya_gas_utilization", nil, v)
	}

	if v, ok := c.readTransactionSuccessRate(ctx); ok {
		c.m.SetGauge("biya_tx_success_rate", nil, v)
	}

	// 检测链分叉：对比 explorer 和 tendermint 节点的区块哈希
	if forkStatus, ok := c.detectFork(ctx); ok {
		c.m.SetGauge("biya_chain_fork", nil, forkStatus)
	} else {
		// 检测失败时，默认设置为 0（无分叉），避免误报
		c.m.SetGauge("biya_chain_fork", nil, 0)
	}

	return nil
}

func (c *RealtimeExplorerCollector) readLatestBlockHeight(ctx context.Context) (float64, bool) {
	raw, err := c.api.GetLatestBlocks(ctx, explorer.CursorPage{Page: 1, PageSize: 1})
	if err != nil {
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_latest_block"}, 0)
		return 0, false
	}
	c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_latest_block"}, 1)

	// apiclient 已剥离 envelope.data，因此这里的结构一般为：
	// {"data":[{"height":"123"}], ...}
	var resp struct {
		Data []struct {
			Height any `json:"height"`
		} `json:"data"`
	}
	if err := jsonUnmarshal(raw, &resp); err != nil {
		c.log.Warn("explorer latest block parse failed", "collector", "realtime_explorer", "method", "readLatestBlockHeight", "err", err)
		return 0, false
	}
	if len(resp.Data) == 0 {
		return 0, false
	}
	v, ok := toFloat64(resp.Data[0].Height)
	return v, ok
}

func (c *RealtimeExplorerCollector) readBlockInterval(ctx context.Context) (float64, bool) {
	// 需要获取至少2个区块来计算时间差
	raw, err := c.api.GetLatestBlocks(ctx, explorer.CursorPage{Page: 1, PageSize: 2})
	if err != nil {
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_interval"}, 0)
		return 0, false
	}
	c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_interval"}, 1)

	// apiclient 已剥离 envelope.data，因此这里的结构一般为：
	// {"data":[{"block_unix_timestamp":1768443554890}, {"block_unix_timestamp":1768443553890}], ...}
	var resp struct {
		Data []struct {
			BlockUnixTimestamp any `json:"block_unix_timestamp"`
		} `json:"data"`
	}
	if err := jsonUnmarshal(raw, &resp); err != nil {
		c.log.Warn("explorer block interval parse failed", "collector", "realtime_explorer", "method", "readBlockInterval", "err", err)
		return 0, false
	}
	if len(resp.Data) < 2 {
		c.log.Warn("explorer block interval: need at least 2 blocks", "collector", "realtime_explorer", "method", "readBlockInterval", "blocks_count", len(resp.Data))
		return 0, false
	}

	// 获取两个区块的时间戳（毫秒）
	timestamp0, ok0 := toFloat64(resp.Data[0].BlockUnixTimestamp)
	timestamp1, ok1 := toFloat64(resp.Data[1].BlockUnixTimestamp)
	if !ok0 || !ok1 {
		c.log.Warn("explorer block interval: failed to parse timestamps", "collector", "realtime_explorer", "method", "readBlockInterval")
		return 0, false
	}

	// 计算时间差（毫秒）并转换为秒（小数）
	intervalMs := timestamp0 - timestamp1
	if intervalMs < 0 {
		// 时间戳顺序异常，返回失败
		c.log.Warn("explorer block interval: negative interval", "collector", "realtime_explorer", "method", "readBlockInterval", "interval_ms", intervalMs)
		return 0, false
	}

	// 毫秒转秒（保留小数）
	intervalSeconds := intervalMs / 1000.0
	return intervalSeconds, true
}

type txStats struct {
	Count24H            float64
	TPS                 float64
	AvgBlockTimeSeconds float64
	ActiveAddresses24H  float64
}

func (c *RealtimeExplorerCollector) readTransactionStats(ctx context.Context) (txStats, bool) {
	raw, err := c.api.GetTransactionStats(ctx)
	if err != nil {
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_transaction_stats"}, 0)
		return txStats{}, false
	}
	c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_transaction_stats"}, 1)

	// apiclient 已剥离 envelope.data，因此这里期望结构为：
	// {"count_24h":..., "tps":..., "avg_block_time":..., "active_addresses_24h":...}
	var resp struct {
		Count24H           any `json:"count_24h"`
		TPS                any `json:"tps"`
		AvgBlockTime       any `json:"avg_block_time"`
		ActiveAddresses24H any `json:"active_addresses_24h"`
	}
	if err := jsonUnmarshal(raw, &resp); err != nil {
		c.log.Warn("explorer tx stats parse failed", "collector", "realtime_explorer", "method", "readTransactionStats", "err", err)
		return txStats{}, false
	}

	out := txStats{
		Count24H:            -1,
		TPS:                 -1,
		AvgBlockTimeSeconds: -1,
		ActiveAddresses24H:  -1,
	}
	if v, ok := toFloat64(resp.Count24H); ok {
		out.Count24H = v
	}
	if v, ok := toFloat64(resp.TPS); ok {
		out.TPS = v
	}
	if v, ok := toFloat64(resp.AvgBlockTime); ok {
		out.AvgBlockTimeSeconds = v
	}
	if v, ok := toFloat64(resp.ActiveAddresses24H); ok {
		out.ActiveAddresses24H = v
	}
	return out, true
}

func (c *RealtimeExplorerCollector) readGasPriceGwei(ctx context.Context) (float64, bool) {
	raw, err := c.api.GetBlockGasUtilization(ctx)
	if err != nil {
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_price"}, 0)
		return 0, false
	}

	// apiclient 已剥离 envelope.data，因此这里期望结构为：
	// {"gas_price": ...}
	var resp struct {
		GasPrice any `json:"gas_price"`
	}
	if err := jsonUnmarshal(raw, &resp); err != nil {
		c.log.Warn("explorer gas price parse failed", "collector", "realtime_explorer", "method", "readGasPriceGwei", "err", err)
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_price"}, 0)
		return 0, false
	}
	v, ok := toFloat64(resp.GasPrice)
	if !ok {
		// 上游在部分环境可能不返回 gas_price 字段；此时视为该 source 不可用，避免“source_up=1 但指标为 0”的误导。
		c.log.Warn("explorer gas price field missing", "collector", "realtime_explorer", "method", "readGasPriceGwei")
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_price"}, 0)
		return 0, false
	}
	c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_price"}, 1)
	return v, true
}

func (c *RealtimeExplorerCollector) readGasUtilization(ctx context.Context) (float64, bool) {
	raw, err := c.api.GetBlockGasUtilization(ctx)
	if err != nil {
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_utilization"}, 0)
		return 0, false
	}

	// apiclient 已剥离 envelope.data，因此这里期望结构为：
	// {"gas_utilization": ...}
	var resp struct {
		GasUtilization any `json:"gas_utilization"`
	}
	if err := jsonUnmarshal(raw, &resp); err != nil {
		c.log.Warn("explorer gas utilization parse failed", "collector", "realtime_explorer", "method", "readGasUtilization", "err", err)
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_utilization"}, 0)
		return 0, false
	}
	v, ok := toFloat64(resp.GasUtilization)
	if !ok {
		// 上游在部分环境可能不返回 gas_utilization 字段；此时视为该 source 不可用，避免"source_up=1 但指标为 0"的误导。
		c.log.Warn("explorer gas utilization field missing", "collector", "realtime_explorer", "method", "readGasUtilization")
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_utilization"}, 0)
		return 0, false
	}
	c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_block_gas_utilization"}, 1)
	return v, true
}

func (c *RealtimeExplorerCollector) readTransactionSuccessRate(ctx context.Context) (float64, bool) {
	raw, err := c.api.GetTransactionFailed1000(ctx)
	if err != nil {
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_transaction_failed_1000"}, 0)
		return 0, false
	}

	// apiclient 已剥离 envelope.data，因此这里期望结构为：
	// {"success_rate_1000": ...}
	var resp struct {
		SuccessRate1000 any `json:"success_rate_1000"`
	}
	if err := jsonUnmarshal(raw, &resp); err != nil {
		c.log.Warn("explorer transaction success rate parse failed", "collector", "realtime_explorer", "method", "readTransactionSuccessRate", "err", err)
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_transaction_failed_1000"}, 0)
		return 0, false
	}
	v, ok := toFloat64(resp.SuccessRate1000)
	if !ok {
		// 上游在部分环境可能不返回 success_rate_1000 字段；此时视为该 source 不可用，避免"source_up=1 但指标为 0"的误导。
		c.log.Warn("explorer transaction success rate field missing", "collector", "realtime_explorer", "method", "readTransactionSuccessRate")
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_transaction_failed_1000"}, 0)
		return 0, false
	}
	c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_transaction_failed_1000"}, 1)
	return v, true
}

// detectFork 检测链分叉：对比 explorer 和 tendermint 节点的区块哈希
// 返回值：(forkStatus, ok)
// - forkStatus: 0=无分叉, 1=检测到分叉
// - ok: 是否成功完成检测
func (c *RealtimeExplorerCollector) detectFork(ctx context.Context) (float64, bool) {
	// 1. 从 explorer 获取最新区块的高度和哈希
	raw, err := c.api.GetLatestBlocks(ctx, explorer.CursorPage{Page: 1, PageSize: 1})
	if err != nil {
		c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_fork_detection"}, 0)
		return 0, false
	}
	c.m.SetGauge("biya_exporter_source_up", map[string]string{"source": "explorer_fork_detection"}, 1)

	var explorerResp struct {
		Data []struct {
			Height    any    `json:"height"`
			BlockHash string `json:"block_hash"`
			Hash      string `json:"hash"` // 兼容不同的字段名
		} `json:"data"`
	}
	if err := jsonUnmarshal(raw, &explorerResp); err != nil {
		c.log.Warn("explorer fork detection: parse failed", "collector", "realtime_explorer", "method", "detectFork", "err", err)
		return 0, false
	}
	if len(explorerResp.Data) == 0 {
		c.log.Warn("explorer fork detection: no block data", "collector", "realtime_explorer", "method", "detectFork")
		return 0, false
	}

	// 获取区块高度和哈希（优先使用 block_hash，其次 hash）
	explorerHeight, ok := toFloat64(explorerResp.Data[0].Height)
	if !ok {
		c.log.Warn("explorer fork detection: failed to parse height", "collector", "realtime_explorer", "method", "detectFork")
		return 0, false
	}

	explorerHash := explorerResp.Data[0].BlockHash
	if explorerHash == "" {
		explorerHash = explorerResp.Data[0].Hash
	}
	if explorerHash == "" {
		c.log.Warn("explorer fork detection: block hash not found", "collector", "realtime_explorer", "method", "detectFork")
		return 0, false
	}

	// 2. 从 tendermint 节点获取对应高度的区块哈希
	if c.tm == nil {
		c.log.Warn("explorer fork detection: tendermint client not available", "collector", "realtime_explorer", "method", "detectFork")
		return 0, false
	}

	heightInt64 := int64(explorerHeight)
	tmBlock, err := c.tm.Block(ctx, heightInt64)
	if err != nil {
		c.log.Warn("explorer fork detection: failed to get block from tendermint", "collector", "realtime_explorer", "method", "detectFork", "height", heightInt64, "err", err)
		return 0, false
	}

	tmHash := strings.TrimSpace(tmBlock.Result.BlockID.Hash)
	explorerHash = strings.TrimSpace(explorerHash)

	// 3. 对比哈希值
	if tmHash == "" {
		c.log.Warn("explorer fork detection: tendermint block hash is empty", "collector", "realtime_explorer", "method", "detectFork", "height", heightInt64)
		return 0, false
	}

	if tmHash != explorerHash {
		c.log.Warn("explorer fork detection: fork detected", "collector", "realtime_explorer", "method", "detectFork",
			"height", heightInt64,
			"explorer_hash", explorerHash,
			"tendermint_hash", tmHash)
		return 1, true // 检测到分叉
	}

	// 哈希相同，无分叉
	return 0, true
}
