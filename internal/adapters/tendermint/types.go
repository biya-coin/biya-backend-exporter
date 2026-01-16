package tendermint

import "time"

// 注意：Tendermint/CometBFT 返回结构会随版本变化，本类型只取我们用到的字段。

type StatusResponse struct {
	Result struct {
		NodeInfo struct {
			Network string `json:"network"`
		} `json:"node_info"`
		SyncInfo struct {
			LatestBlockHeight string    `json:"latest_block_height"`
			LatestBlockTime   time.Time `json:"latest_block_time"`
			CatchingUp        bool      `json:"catching_up"`
		} `json:"sync_info"`
	} `json:"result"`
}

// BlockResponse 对应 Tendermint/CometBFT JSON-RPC 2.0 响应格式
// 实际响应格式：
//
//	{
//	  "jsonrpc": "2.0",
//	  "id": -1,
//	  "result": {
//	    "block_id": { "hash": "...", "parts": {...} },
//	    "block": { "header": {...}, "data": {...} }
//	  }
//	}
type BlockResponse struct {
	JSONRPC string `json:"jsonrpc"` // JSON-RPC 版本，通常为 "2.0"
	ID      int    `json:"id"`      // 请求 ID，通常为 -1
	Result  struct {
		BlockID struct {
			Hash  string `json:"hash"`
			Parts struct {
				Total int    `json:"total"`
				Hash  string `json:"hash"`
			} `json:"parts"`
		} `json:"block_id"`
		Block struct {
			Header struct {
				Height string    `json:"height"`
				Time   time.Time `json:"time"`
			} `json:"header"`
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	} `json:"result"`
}

type NumUnconfirmedTxsResponse struct {
	Result struct {
		NTxs  string `json:"n_txs"`
		Total string `json:"total"`
	} `json:"result"`
}
