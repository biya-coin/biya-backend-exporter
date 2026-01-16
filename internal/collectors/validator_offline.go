package collectors

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/biya-coin/biya-dex-backend-exporter/internal/adapters/explorer"
	"github.com/biya-coin/biya-dex-backend-exporter/internal/adapters/tendermint"
	"github.com/biya-coin/biya-dex-backend-exporter/internal/metrics"
)

type ValidatorOfflineCollector struct {
	log         *slog.Logger
	m           *metrics.Metrics
	nodeIPs     []string
	httpTimeout time.Duration
	explorerCli *explorer.Client
}

func NewValidatorOfflineCollector(log *slog.Logger, m *metrics.Metrics, nodeIPs []string, httpTimeout time.Duration, explorerCli *explorer.Client) *ValidatorOfflineCollector {
	return &ValidatorOfflineCollector{
		log:         log,
		m:           m,
		nodeIPs:     nodeIPs,
		httpTimeout: httpTimeout,
		explorerCli: explorerCli,
	}
}

func (c *ValidatorOfflineCollector) Run(ctx context.Context) error {
	if len(c.nodeIPs) == 0 {
		c.log.Debug("no validator nodes configured, skipping offline check")
		return nil
	}

	// 获取当前区块高度（biya_block_height）
	currentBlockHeight, err := c.getCurrentBlockHeight(ctx)
	if err != nil {
		c.log.Warn("failed to get current block height, skipping behind_blocks calculation", "err", err)
		// 即使获取不到当前区块高度，仍然可以检查节点离线状态
		currentBlockHeight = 0
	}

	// 顺序检查所有节点
	for _, nodeIP := range c.nodeIPs {
		status, nodeHeight := c.checkNodeStatus(ctx, nodeIP)
		
		// 设置离线指标：1=离线, 0=在线
		offlineValue := 0.0
		if status == nodeStatusOffline {
			offlineValue = 1.0
		}
		c.m.SetGauge("biya_validator_offline", map[string]string{"node_ip": nodeIP}, offlineValue)

		// 计算并设置落后区块数
		if currentBlockHeight > 0 && nodeHeight > 0 {
			behindBlocks := currentBlockHeight - nodeHeight
			if behindBlocks < 0 {
				behindBlocks = 0 // 如果节点高度大于当前高度（不应该发生），设为0
			}
			c.m.SetGauge("biya_validator_behind_blocks", map[string]string{"node_ip": nodeIP}, float64(behindBlocks))
		} else {
			// 如果无法获取高度，设置为0（表示无法计算）
			c.m.SetGauge("biya_validator_behind_blocks", map[string]string{"node_ip": nodeIP}, 0)
		}
	}

	return nil
}

type nodeStatus int

const (
	nodeStatusOnline  nodeStatus = 0
	nodeStatusOffline nodeStatus = 1
)

// getCurrentBlockHeight 获取当前区块高度（biya_block_height）
// 根据 provide.md，从 explorer API 获取
func (c *ValidatorOfflineCollector) getCurrentBlockHeight(ctx context.Context) (int64, error) {
	if c.explorerCli == nil {
		return 0, fmt.Errorf("explorer client not configured")
	}

	raw, err := c.explorerCli.GetLatestBlocks(ctx, explorer.CursorPage{Page: 1, PageSize: 1})
	if err != nil {
		return 0, fmt.Errorf("failed to get latest blocks: %w", err)
	}

	// 解析响应，获取 height
	var resp struct {
		Data []struct {
			Height any `json:"height"`
		} `json:"data"`
	}
	if err := jsonUnmarshal(raw, &resp); err != nil {
		return 0, fmt.Errorf("failed to parse latest blocks response: %w", err)
	}
	if len(resp.Data) == 0 {
		return 0, fmt.Errorf("no block data in response")
	}

	height, ok := toFloat64(resp.Data[0].Height)
	if !ok {
		return 0, fmt.Errorf("invalid height value")
	}

	return int64(height), nil
}

// checkNodeStatus 检查节点状态并返回节点高度
// 返回：节点状态（在线/离线）和节点区块高度
func (c *ValidatorOfflineCollector) checkNodeStatus(ctx context.Context, nodeIP string) (nodeStatus, int64) {
	// 构建节点 URL，根据 provide.md，使用 https://${节点ip}:26757/status
	// 但配置中的 IP 可能不包含协议，需要判断
	baseURL := nodeIP
	if !strings.HasPrefix(nodeIP, "http://") && !strings.HasPrefix(nodeIP, "https://") {
		// 默认使用 http（内网环境通常使用 http）
		baseURL = fmt.Sprintf("http://%s:26757", nodeIP)
	}

	// 创建临时的 tendermint client
	tmClient := tendermint.NewClient(baseURL, c.httpTimeout)

	// 尝试调用 /status 接口
	statusResp, err := tmClient.Status(ctx)
	if err != nil {
		c.log.Debug("validator node offline", "node_ip", nodeIP, "err", err)
		return nodeStatusOffline, 0 // 请求失败，节点离线
	}

	// 解析节点区块高度
	nodeHeight, err := strconv.ParseInt(statusResp.Result.SyncInfo.LatestBlockHeight, 10, 64)
	if err != nil {
		c.log.Debug("failed to parse node block height", "node_ip", nodeIP, "err", err)
		return nodeStatusOnline, 0 // 节点在线但无法获取高度
	}

	c.log.Debug("validator node online", "node_ip", nodeIP, "height", nodeHeight)
	return nodeStatusOnline, nodeHeight // 请求成功，节点在线
}

