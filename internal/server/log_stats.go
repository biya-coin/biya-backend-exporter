package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/biya-coin/biya-dex-backend-exporter/internal/adapters/loki"
)

// LogStatsService 处理日志统计查询
type LogStatsService struct {
	lokiClient *loki.Client
	logger     *slog.Logger
}

// NewLogStatsService 创建日志统计服务
func NewLogStatsService(lokiClient *loki.Client, logger *slog.Logger) *LogStatsService {
	return &LogStatsService{
		lokiClient: lokiClient,
		logger:     logger,
	}
}

// LogStatsResponse 日志统计响应
type LogStatsResponse struct {
	Success bool          `json:"success"`
	Data    *LogStatsData `json:"data"`
	Message string        `json:"message,omitempty"`
}

// LogStatsData 日志统计数据
type LogStatsData struct {
	Total     int64     `json:"total"`      // 总日志数
	Error     int64     `json:"error"`      // 错误日志数
	Warning   int64     `json:"warning"`    // 警告日志数
	Info      int64     `json:"info"`       // 信息日志数
	Debug     int64     `json:"debug"`      // 调试日志数
	TimeRange TimeRange `json:"time_range"` // 时间范围
}

// TimeRange 时间范围
type TimeRange struct {
	Start int64 `json:"start"` // Unix 时间戳（秒）
	End   int64 `json:"end"`   // Unix 时间戳（秒）
}

// LogEntry 日志条目
type LogEntry struct {
	Timestamp int64             `json:"timestamp"` // Unix 时间戳（毫秒）
	Level     string            `json:"level"`     // 日志级别
	Content   string            `json:"content"`   // 日志内容
	Labels    map[string]string `json:"labels"`    // 标签（node_id, module 等）
}

// LogListResponse 日志列表响应
type LogListResponse struct {
	Success bool         `json:"success"`
	Data    *LogListData `json:"data"`
	Message string       `json:"message,omitempty"`
}

// LogListData 日志列表数据
type LogListData struct {
	Logs     []LogEntry `json:"logs"`
	Total    int        `json:"total"`     // 总条数
	Page     int        `json:"page"`      // 当前页码
	PageSize int        `json:"page_size"` // 每页大小
	HasMore  bool       `json:"has_more"`  // 是否有更多
}

// GetLogStats 获取日志统计（按级别和时间范围）
// 直接使用 LogQL 的聚合函数一次性查询所有统计信息
func (s *LogStatsService) GetLogStats(ctx context.Context, nodeID string, start, end time.Time) (*LogStatsResponse, error) {
	// 构建基础查询
	baseQuery := `{job="injective-node"}`
	if nodeID != "" {
		baseQuery = fmt.Sprintf(`{job="injective-node",node_id="%s"}`, nodeID)
	}

	// 计算时间窗口（使用较小的窗口，比如 1h，然后在 query_range 中指定时间范围）
	// count_over_time 的窗口应该小于等于查询时间范围
	timeWindow := end.Sub(start)
	if timeWindow > 24*time.Hour {
		timeWindow = 24 * time.Hour // 最大24小时窗口
	}
	if timeWindow < time.Minute {
		timeWindow = time.Minute // 最小1分钟窗口
	}
	duration := s.durationToPromQL(timeWindow)

	// 使用 LogQL 直接查询统计：
	// 1. 总日志数: sum(count_over_time(...))
	// 2. 按级别统计: sum by (level) (count_over_time(...))

	var total, errorCount, warningCount, infoCount, debugCount int64

	// 查询总日志数（使用 LogQL 的 sum 聚合函数）
	// 注意：count_over_time 返回时间序列（matrix 格式），需要取最后一个值（最新的统计）
	totalQuery := fmt.Sprintf(`sum(count_over_time(%s[%s]))`, baseQuery, duration)
	totalResult, err := s.lokiClient.QueryRange(ctx, totalQuery, start, end, 1000)
	if err != nil {
		s.logger.Warn("failed to query total logs", "error", err, "query", totalQuery)
	} else if len(totalResult.Data.Result) > 0 {
		// 取最后一个值（最新的统计值）
		item := totalResult.Data.Result[0]
		if len(item.Values) > 0 {
			// 取最后一个值（matrix 格式：[[timestamp, value], ...]，都是数字）
			lastValue := item.Values[len(item.Values)-1]
			if valueArray, ok := lastValue.([]interface{}); ok && len(valueArray) >= 2 {
				// 第二个元素是统计值（可能是字符串或数字）
				var val float64
				switch v := valueArray[1].(type) {
				case string:
					if parsed, err := strconv.ParseFloat(v, 64); err == nil {
						val = parsed
					}
				case float64:
					val = v
				case int64:
					val = float64(v)
				}
				if val > 0 {
					total = int64(val)
					s.logger.Debug("total logs count", "total", total, "query", totalQuery)
				}
			}
		}
	} else {
		s.logger.Warn("no result from total query", "query", totalQuery)
	}

	// 使用 LogQL 的 sum by (level) 一次性查询所有级别的统计
	// 注意：返回的是 matrix 格式（Prometheus 格式），metric 字段包含 level 标签
	statsQuery := fmt.Sprintf(`sum by (level) (count_over_time(%s[%s]))`, baseQuery, duration)
	result, err := s.lokiClient.QueryRange(ctx, statsQuery, start, end, 1000)
	if err != nil {
		s.logger.Warn("failed to query stats by level", "error", err, "query", statsQuery)
		// 如果聚合查询失败，使用备用方案分别查询
		return s.getLogStatsFallback(ctx, baseQuery, duration, start, end)
	} else if len(result.Data.Result) > 0 {
		// 解析按级别统计的结果（取每个级别的最后一个值）
		for _, item := range result.Data.Result {
			level := ""
			// 统计查询返回的是 metric 格式，不是 stream 格式
			if l, ok := item.Metric["level"]; ok {
				level = strings.ToUpper(l)
			} else if l, ok := item.Stream["level"]; ok {
				// 兼容 stream 格式
				level = strings.ToUpper(l)
			}

			// 取最后一个值（最新的统计值）
			// matrix 格式：values 是 []interface{}，每个元素是 [timestamp, value]
			if len(item.Values) > 0 {
				lastValue := item.Values[len(item.Values)-1]
				if valueArray, ok := lastValue.([]interface{}); ok && len(valueArray) >= 2 {
					// 第二个元素是统计值（可能是字符串或数字）
					var count float64
					switch v := valueArray[1].(type) {
					case string:
						if parsed, err := strconv.ParseFloat(v, 64); err == nil {
							count = parsed
						} else {
							continue
						}
					case float64:
						count = v
					case int64:
						count = float64(v)
					default:
						continue
					}

					switch {
					case strings.Contains(level, "ERR") || strings.Contains(level, "ERROR") || strings.Contains(level, "FATAL") || strings.Contains(level, "PANIC"):
						errorCount += int64(count)
					case strings.Contains(level, "WARN"):
						warningCount += int64(count)
					case strings.Contains(level, "INF") || strings.Contains(level, "INFO"):
						infoCount += int64(count)
					case strings.Contains(level, "DBG") || strings.Contains(level, "DEBUG"):
						debugCount += int64(count)
					}
				}
			}
		}
		s.logger.Debug("stats by level", "error", errorCount, "warning", warningCount, "info", infoCount, "debug", debugCount)
	} else {
		s.logger.Warn("no result from stats query", "query", statsQuery)
		// 如果聚合查询失败，使用备用方案分别查询
		return s.getLogStatsFallback(ctx, baseQuery, duration, start, end)
	}

	// 如果总数未查询到，使用各级别之和
	if total == 0 {
		total = errorCount + warningCount + infoCount + debugCount
	}

	return &LogStatsResponse{
		Success: true,
		Data: &LogStatsData{
			Total:   total,
			Error:   errorCount,
			Warning: warningCount,
			Info:    infoCount,
			Debug:   debugCount,
			TimeRange: TimeRange{
				Start: start.Unix(),
				End:   end.Unix(),
			},
		},
	}, nil
}

// getLogStatsFallback 备用方案：分别查询各级别统计
func (s *LogStatsService) getLogStatsFallback(ctx context.Context, baseQuery, duration string, start, end time.Time) (*LogStatsResponse, error) {
	var total, errorCount, warningCount, infoCount, debugCount int64

	// 查询总日志数（取最后一个值，matrix 格式）
	totalQuery := fmt.Sprintf(`sum(count_over_time(%s[%s]))`, baseQuery, duration)
	if result, err := s.lokiClient.QueryRange(ctx, totalQuery, start, end, 1000); err == nil && len(result.Data.Result) > 0 {
		item := result.Data.Result[0]
		if len(item.Values) > 0 {
			lastValue := item.Values[len(item.Values)-1]
			if valueArray, ok := lastValue.([]interface{}); ok && len(valueArray) >= 2 {
				switch v := valueArray[1].(type) {
				case string:
					if val, err := strconv.ParseFloat(v, 64); err == nil {
						total = int64(val)
					}
				case float64:
					total = int64(v)
				case int64:
					total = v
				}
			}
		}
	}

	// 查询错误日志数（取最后一个值，matrix 格式）
	errorQuery := fmt.Sprintf(`sum(count_over_time(%s{level=~"ERR|ERROR|FATAL|PANIC"}[%s]))`, baseQuery, duration)
	if result, err := s.lokiClient.QueryRange(ctx, errorQuery, start, end, 1000); err == nil && len(result.Data.Result) > 0 {
		item := result.Data.Result[0]
		if len(item.Values) > 0 {
			lastValue := item.Values[len(item.Values)-1]
			if valueArray, ok := lastValue.([]interface{}); ok && len(valueArray) >= 2 {
				switch v := valueArray[1].(type) {
				case string:
					if val, err := strconv.ParseFloat(v, 64); err == nil {
						errorCount = int64(val)
					}
				case float64:
					errorCount = int64(v)
				case int64:
					errorCount = v
				}
			}
		}
	}

	// 查询警告日志数（取最后一个值，matrix 格式）
	warningQuery := fmt.Sprintf(`sum(count_over_time(%s{level=~"WARN|WARNING"}[%s]))`, baseQuery, duration)
	if result, err := s.lokiClient.QueryRange(ctx, warningQuery, start, end, 1000); err == nil && len(result.Data.Result) > 0 {
		item := result.Data.Result[0]
		if len(item.Values) > 0 {
			lastValue := item.Values[len(item.Values)-1]
			if valueArray, ok := lastValue.([]interface{}); ok && len(valueArray) >= 2 {
				switch v := valueArray[1].(type) {
				case string:
					if val, err := strconv.ParseFloat(v, 64); err == nil {
						warningCount = int64(val)
					}
				case float64:
					warningCount = int64(v)
				case int64:
					warningCount = v
				}
			}
		}
	}

	// 查询信息日志数（取最后一个值，matrix 格式）
	infoQuery := fmt.Sprintf(`sum(count_over_time(%s{level=~"INF|INFO"}[%s]))`, baseQuery, duration)
	if result, err := s.lokiClient.QueryRange(ctx, infoQuery, start, end, 1000); err == nil && len(result.Data.Result) > 0 {
		item := result.Data.Result[0]
		if len(item.Values) > 0 {
			lastValue := item.Values[len(item.Values)-1]
			if valueArray, ok := lastValue.([]interface{}); ok && len(valueArray) >= 2 {
				switch v := valueArray[1].(type) {
				case string:
					if val, err := strconv.ParseFloat(v, 64); err == nil {
						infoCount = int64(val)
					}
				case float64:
					infoCount = int64(v)
				case int64:
					infoCount = v
				}
			}
		}
	}

	// 查询调试日志数（取最后一个值，matrix 格式）
	debugQuery := fmt.Sprintf(`sum(count_over_time(%s{level=~"DBG|DEBUG"}[%s]))`, baseQuery, duration)
	if result, err := s.lokiClient.QueryRange(ctx, debugQuery, start, end, 1000); err == nil && len(result.Data.Result) > 0 {
		item := result.Data.Result[0]
		if len(item.Values) > 0 {
			lastValue := item.Values[len(item.Values)-1]
			if valueArray, ok := lastValue.([]interface{}); ok && len(valueArray) >= 2 {
				switch v := valueArray[1].(type) {
				case string:
					if val, err := strconv.ParseFloat(v, 64); err == nil {
						debugCount = int64(val)
					}
				case float64:
					debugCount = int64(v)
				case int64:
					debugCount = v
				}
			}
		}
	}

	if total == 0 {
		total = errorCount + warningCount + infoCount + debugCount
	}

	return &LogStatsResponse{
		Success: true,
		Data: &LogStatsData{
			Total:   total,
			Error:   errorCount,
			Warning: warningCount,
			Info:    infoCount,
			Debug:   debugCount,
			TimeRange: TimeRange{
				Start: start.Unix(),
				End:   end.Unix(),
			},
		},
	}, nil
}

// GetLogList 获取日志列表（支持级别、时间范围、关键词过滤）
func (s *LogStatsService) GetLogList(ctx context.Context, nodeID, level, keyword string, start, end time.Time, page, pageSize int) (*LogListResponse, error) {
	// 构建查询
	query := `{job="injective-node"}`

	if nodeID != "" {
		query = fmt.Sprintf(`{job="injective-node",node_id="%s"}`, nodeID)
	}

	// 添加级别过滤
	if level != "" && level != "all" {
		levelMap := map[string]string{
			"error":   `{level=~"ERR|ERROR|FATAL|PANIC"}`,
			"warning": `{level=~"WARN|WARNING"}`,
			"info":    `{level=~"INF|INFO"}`,
			"debug":   `{level=~"DBG|DEBUG"}`,
		}
		if levelFilter, ok := levelMap[strings.ToLower(level)]; ok {
			// 合并标签
			query = strings.TrimSuffix(query, "}") + "," + strings.TrimPrefix(levelFilter, "{")
		}
	}

	// 添加关键词过滤
	if keyword != "" {
		query = fmt.Sprintf(`%s |= "%s"`, query, keyword)
	}

	// 查询日志（从新到旧，取足够多的数据用于分页）
	limit := page * pageSize
	if limit < 100 {
		limit = 100 // 至少查询100条
	}
	if limit > 1000 {
		limit = 1000 // 最多查询1000条
	}

	result, err := s.lokiClient.QueryRange(ctx, query, start, end, limit)
	if err != nil {
		return &LogListResponse{
			Success: false,
			Message: fmt.Sprintf("查询日志失败: %v", err),
		}, nil
	}

	// 解析日志条目（streams 格式：values 是 [][]string，但 JSON 解析后可能是 []interface{}）
	var allLogs []LogEntry
	for _, stream := range result.Data.Result {
		level := "info"
		if l, ok := stream.Stream["level"]; ok {
			level = strings.ToLower(l)
		}

		for _, valueInterface := range stream.Values {
			// streams 格式：value 是 [timestamp, log_line]
			// JSON 解析后可能是 []interface{} 或 []string
			var timestampStr, content string

			if valueArray, ok := valueInterface.([]interface{}); ok {
				// 如果是 interface{} 数组
				if len(valueArray) < 2 {
					continue
				}
				// 第一个元素是时间戳（可能是字符串或数字）
				switch v := valueArray[0].(type) {
				case string:
					timestampStr = v
				case float64:
					timestampStr = strconv.FormatFloat(v, 'f', 0, 64)
				case int64:
					timestampStr = strconv.FormatInt(v, 10)
				default:
					timestampStr = fmt.Sprintf("%v", v)
				}
				// 第二个元素是日志内容（字符串）
				if str, ok := valueArray[1].(string); ok {
					content = str
				} else {
					content = fmt.Sprintf("%v", valueArray[1])
				}
			} else if valueStrArray, ok := valueInterface.([]string); ok {
				// 如果是 string 数组
				if len(valueStrArray) < 2 {
					continue
				}
				timestampStr = valueStrArray[0]
				content = valueStrArray[1]
			} else {
				// 尝试通过 JSON 重新解析
				continue
			}

			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				continue
			}

			allLogs = append(allLogs, LogEntry{
				Timestamp: timestamp / 1000000, // 转换为毫秒
				Level:     level,
				Content:   content,
				Labels:    stream.Stream,
			})
		}
	}

	// 分页处理（从新到旧，所以需要反转）
	total := len(allLogs)
	startIdx := total - page*pageSize
	endIdx := total - (page-1)*pageSize

	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > total {
		endIdx = total
	}

	var logs []LogEntry
	if startIdx < endIdx {
		// 反转切片以获取最新的日志在前
		for i := endIdx - 1; i >= startIdx; i-- {
			logs = append(logs, allLogs[i])
		}
	}

	hasMore := startIdx > 0

	return &LogListResponse{
		Success: true,
		Data: &LogListData{
			Logs:     logs,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			HasMore:  hasMore,
		},
	}, nil
}

// durationToPromQL 将 duration 转换为 PromQL 格式
func (s *LogStatsService) durationToPromQL(d time.Duration) string {
	seconds := int(d.Seconds())
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%dh", hours)
	}
	days := hours / 24
	return fmt.Sprintf("%dd", days)
}

// HandleLogStats 处理日志统计查询请求
func (s *LogStatsService) HandleLogStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	nodeID := r.URL.Query().Get("node_id")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	// 解析时间范围
	var start, end time.Time
	now := time.Now()

	if startStr != "" {
		if ts, err := strconv.ParseInt(startStr, 10, 64); err == nil {
			start = time.Unix(ts, 0)
		} else {
			start = now.Add(-24 * time.Hour) // 默认24小时前
		}
	} else {
		start = now.Add(-24 * time.Hour) // 默认24小时前
	}

	if endStr != "" {
		if ts, err := strconv.ParseInt(endStr, 10, 64); err == nil {
			end = time.Unix(ts, 0)
		} else {
			end = now
		}
	} else {
		end = now
	}

	// 查询统计
	response, err := s.GetLogStats(r.Context(), nodeID, start, end)
	if err != nil {
		s.logger.Error("failed to get log stats", "error", err)
		response = &LogStatsResponse{
			Success: false,
			Message: fmt.Sprintf("Internal error: %v", err),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}

// HandleLogList 处理日志列表查询请求
func (s *LogStatsService) HandleLogList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析查询参数
	nodeID := r.URL.Query().Get("node_id")
	level := r.URL.Query().Get("level")
	keyword := r.URL.Query().Get("keyword")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	// 解析时间范围
	var start, end time.Time
	now := time.Now()

	if startStr != "" {
		if ts, err := strconv.ParseInt(startStr, 10, 64); err == nil {
			start = time.Unix(ts, 0)
		} else {
			start = now.Add(-24 * time.Hour)
		}
	} else {
		start = now.Add(-24 * time.Hour)
	}

	if endStr != "" {
		if ts, err := strconv.ParseInt(endStr, 10, 64); err == nil {
			end = time.Unix(ts, 0)
		} else {
			end = now
		}
	} else {
		end = now
	}

	// 解析分页参数
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 50
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 200 {
			pageSize = ps
		}
	}

	// 查询日志列表
	response, err := s.GetLogList(r.Context(), nodeID, level, keyword, start, end, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get log list", "error", err)
		response = &LogListResponse{
			Success: false,
			Message: fmt.Sprintf("Internal error: %v", err),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if !response.Success {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(response)
}
