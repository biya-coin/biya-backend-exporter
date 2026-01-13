package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client Loki API 客户端
type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient 创建 Loki 客户端
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

// QueryRangeResponse Loki 查询范围响应
// 注意：对于日志查询返回 stream 格式，对于统计查询（sum/count_over_time）返回 matrix 格式
type QueryRangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"` // "streams" 或 "matrix"
		Result     []struct {
			// 日志查询格式（resultType="streams"）
			Stream map[string]string `json:"stream,omitempty"`
			// 统计查询格式（resultType="matrix"）
			Metric map[string]string `json:"metric,omitempty"` // 标签（如 level）
			// Values: 对于 streams 是 [][]string，对于 matrix 是 [][]interface{}（时间戳和值都是数字）
			Values []interface{} `json:"values"`
		} `json:"result"`
		Stats interface{} `json:"stats"`
	} `json:"data"`
}

// QueryResponse Loki 即时查询响应
type QueryResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// QueryRange 查询时间范围内的日志
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, limit int) (*QueryRangeResponse, error) {
	u := c.baseURL + "/loki/api/v1/query_range"

	params := url.Values{}
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	params.Set("end", strconv.FormatInt(end.UnixNano(), 10))
	params.Set("limit", strconv.Itoa(limit))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("loki query_range failed: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result QueryRangeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Query 即时查询日志
func (c *Client) Query(ctx context.Context, query string, limit int) (*QueryResponse, error) {
	u := c.baseURL + "/loki/api/v1/query"

	params := url.Values{}
	params.Set("query", query)
	params.Set("limit", strconv.Itoa(limit))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("loki query failed: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result QueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// QueryStats 使用 LogQL 统计查询（返回数值结果）
func (c *Client) QueryStats(ctx context.Context, query string, start, end time.Time) (float64, error) {
	u := c.baseURL + "/loki/api/v1/query_range"

	params := url.Values{}
	params.Set("query", query)
	params.Set("start", strconv.FormatInt(start.UnixNano(), 10))
	params.Set("end", strconv.FormatInt(end.UnixNano(), 10))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+params.Encode(), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("loki query_stats failed: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var result struct {
		Status string `json:"status"`
		Data   struct {
			ResultType string `json:"resultType"`
			Result     []struct {
				Metric map[string]string `json:"metric"`
				Value  []interface{}     `json:"value"` // [timestamp, value]
			} `json:"result"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// 提取统计值
	if len(result.Data.Result) > 0 && len(result.Data.Result[0].Value) >= 2 {
		if val, ok := result.Data.Result[0].Value[1].(string); ok {
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				return f, nil
			}
		}
	}

	return 0, nil
}
