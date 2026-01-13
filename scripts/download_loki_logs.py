#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Loki 日志下载工具
支持下载完整的日志并导出为 JSON 或文本格式

使用方法:
    # 下载最近24小时的所有日志
    python3 scripts/download_loki_logs.py --query '{job="injective-node"}' --hours 24

    # 下载指定时间范围的日志
    python3 scripts/download_loki_logs.py --query '{job="injective-node"}' --start "2024-01-01T00:00:00Z" --end "2024-01-02T00:00:00Z"

    # 导出为文本格式
    python3 scripts/download_loki_logs.py --query '{job="injective-node"}' --hours 24 --format txt

    # 指定输出文件
    python3 scripts/download_loki_logs.py --query '{job="injective-node"}' --hours 24 --output logs.txt
"""

import argparse
import json
import sys
import time
from datetime import datetime, timedelta
from typing import List, Dict, Any
import requests
from urllib.parse import urlencode


class LokiLogDownloader:
    def __init__(self, base_url: str = "http://localhost:3100", timeout: int = 300):
        """
        初始化 Loki 日志下载器
        
        Args:
            base_url: Loki 服务器地址
            timeout: 请求超时时间（秒）
        """
        self.base_url = base_url.rstrip('/')
        self.timeout = timeout
        self.session = requests.Session()
    
    def query_range(self, query: str, start: datetime, end: datetime, limit: int = 5000, direction: str = "forward") -> Dict[str, Any]:
        """
        查询时间范围内的日志
        
        Args:
            query: LogQL 查询表达式
            start: 开始时间
            end: 结束时间
            limit: 每次查询的最大条目数（默认5000，Loki最大支持10000）
            direction: 查询方向 "forward" 或 "backward"
        
        Returns:
            Loki API 响应
        """
        url = f"{self.base_url}/loki/api/v1/query_range"
        
        params = {
            "query": query,
            "start": int(start.timestamp() * 1e9),  # 转换为纳秒
            "end": int(end.timestamp() * 1e9),
            "limit": limit,
            "direction": direction
        }
        
        try:
            response = self.session.get(url, params=params, timeout=self.timeout)
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException as e:
            print(f"查询失败: {e}", file=sys.stderr)
            if hasattr(e.response, 'text'):
                print(f"响应内容: {e.response.text}", file=sys.stderr)
            raise
    
    def download_all_logs(self, query: str, start: datetime, end: datetime, 
                          limit_per_request: int = 5000, 
                          max_entries: int = None) -> List[Dict[str, Any]]:
        """
        下载指定时间范围内的所有日志（自动处理分页）
        
        Args:
            query: LogQL 查询表达式
            start: 开始时间
            end: 结束时间
            limit_per_request: 每次请求的最大条目数
            max_entries: 最大下载条目数（None 表示不限制）
        
        Returns:
            所有日志条目的列表，每个条目格式: {"timestamp": int, "labels": dict, "content": str}
        """
        all_logs = []
        current_start = start
        total_fetched = 0
        
        print(f"开始下载日志...")
        print(f"查询: {query}")
        print(f"时间范围: {start.isoformat()} 到 {end.isoformat()}")
        print(f"每页限制: {limit_per_request}")
        if max_entries:
            print(f"最大条目数: {max_entries}")
        print("-" * 60)
        
        while current_start < end:
            if max_entries and total_fetched >= max_entries:
                print(f"\n已达到最大条目数限制: {max_entries}")
                break
            
            # 计算本次查询的结束时间（如果接近总结束时间，使用总结束时间）
            query_end = min(current_start + timedelta(hours=1), end)
            
            try:
                response = self.query_range(query, current_start, query_end, limit_per_request)
                
                if response.get("status") != "success":
                    print(f"警告: 查询返回状态: {response.get('status')}", file=sys.stderr)
                    break
                
                result_type = response.get("data", {}).get("resultType")
                if result_type != "streams":
                    print(f"警告: 意外的结果类型: {result_type}", file=sys.stderr)
                    break
                
                streams = response.get("data", {}).get("result", [])
                
                if not streams:
                    # 如果没有更多日志，尝试下一个时间窗口
                    current_start = query_end
                    continue
                
                # 提取日志条目
                page_logs = []
                total_values_count = 0
                for stream in streams:
                    labels = stream.get("stream", {})
                    values = stream.get("values", [])
                    total_values_count += len(values)
                    
                    for value in values:
                        if len(value) < 2:
                            continue
                        timestamp_ns = value[0]
                        content = value[1]
                        
                        try:
                            # 处理时间戳（可能是字符串或数字）
                            if isinstance(timestamp_ns, str):
                                timestamp_ns = int(timestamp_ns)
                            timestamp = timestamp_ns // 1_000_000  # 转换为毫秒
                            
                            page_logs.append({
                                "timestamp": timestamp,
                                "timestamp_ns": timestamp_ns,
                                "labels": labels,
                                "content": content
                            })
                        except (ValueError, TypeError) as e:
                            print(f"警告: 跳过无效的时间戳: {timestamp_ns}", file=sys.stderr)
                            continue
                
                # 按时间戳排序
                page_logs.sort(key=lambda x: x["timestamp"])
                
                # 检查是否达到最大条目数
                remaining = max_entries - total_fetched if max_entries else len(page_logs)
                if max_entries and remaining < len(page_logs):
                    page_logs = page_logs[:remaining]
                
                all_logs.extend(page_logs)
                total_fetched += len(page_logs)
                
                # 显示进度
                if page_logs:
                    start_time_str = datetime.fromtimestamp(page_logs[0]['timestamp']/1000).strftime("%Y-%m-%d %H:%M:%S")
                    end_time_str = datetime.fromtimestamp(page_logs[-1]['timestamp']/1000).strftime("%Y-%m-%d %H:%M:%S")
                    print(f"已下载: {total_fetched} 条日志 (时间: {start_time_str} 到 {end_time_str})")
                else:
                    print(f"已下载: {total_fetched} 条日志 (当前时间窗口无日志)")
                
                # 判断是否需要继续查询
                # 如果返回的总日志数等于limit，可能还有更多日志
                # 使用最后一条日志的时间戳作为下一次查询的起始时间
                if page_logs and total_values_count >= limit_per_request:
                    # 可能还有更多日志，从最后一条日志的时间继续查询
                    last_timestamp = page_logs[-1]["timestamp_ns"]
                    current_start = datetime.fromtimestamp(last_timestamp / 1e9)
                    # 避免重复，稍微增加一点时间（1微秒）
                    current_start += timedelta(microseconds=1)
                else:
                    # 这个时间窗口已经查询完毕，移动到下一个时间窗口
                    current_start = query_end
                
                # 如果达到最大条目数，停止
                if max_entries and total_fetched >= max_entries:
                    break
                
                # 避免请求过快
                time.sleep(0.1)
                
            except Exception as e:
                print(f"查询出错: {e}", file=sys.stderr)
                # 如果出错，尝试下一个时间窗口
                current_start = query_end
                continue
        
        print("-" * 60)
        print(f"下载完成！总共 {len(all_logs)} 条日志")
        
        return all_logs
    
    def export_to_json(self, logs: List[Dict[str, Any]], output_file: str):
        """导出日志为 JSON 格式"""
        with open(output_file, 'w', encoding='utf-8') as f:
            json.dump(logs, f, ensure_ascii=False, indent=2)
        print(f"已导出 JSON 格式到: {output_file}")
    
    def export_to_text(self, logs: List[Dict[str, Any]], output_file: str):
        """导出日志为文本格式（类似日志文件）"""
        with open(output_file, 'w', encoding='utf-8') as f:
            for log in logs:
                timestamp = datetime.fromtimestamp(log["timestamp"] / 1000)
                labels_str = " ".join([f"{k}={v}" for k, v in log["labels"].items()])
                f.write(f"[{timestamp.isoformat()}] [{labels_str}] {log['content']}\n")
        print(f"已导出文本格式到: {output_file}")
    
    def export_to_csv(self, logs: List[Dict[str, Any]], output_file: str):
        """导出日志为 CSV 格式"""
        import csv
        with open(output_file, 'w', encoding='utf-8', newline='') as f:
            writer = csv.writer(f)
            writer.writerow(["timestamp", "timestamp_iso", "labels", "content"])
            for log in logs:
                timestamp = datetime.fromtimestamp(log["timestamp"] / 1000)
                labels_str = json.dumps(log["labels"], ensure_ascii=False)
                writer.writerow([
                    log["timestamp"],
                    timestamp.isoformat(),
                    labels_str,
                    log["content"]
                ])
        print(f"已导出 CSV 格式到: {output_file}")


def parse_time(time_str: str) -> datetime:
    """解析时间字符串（支持多种格式）"""
    formats = [
        "%Y-%m-%dT%H:%M:%S",
        "%Y-%m-%dT%H:%M:%SZ",
        "%Y-%m-%dT%H:%M:%S%z",
        "%Y-%m-%d %H:%M:%S",
        "%Y-%m-%d",
    ]
    
    for fmt in formats:
        try:
            return datetime.strptime(time_str, fmt)
        except ValueError:
            continue
    
    raise ValueError(f"无法解析时间格式: {time_str}")


def main():
    parser = argparse.ArgumentParser(
        description="从 Loki 下载完整的日志",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 下载最近24小时的所有日志
  %(prog)s --query '{job="injective-node"}' --hours 24

  # 下载指定时间范围的日志
  %(prog)s --query '{job="injective-node"}' --start "2024-01-01T00:00:00Z" --end "2024-01-02T00:00:00Z"

  # 导出为文本格式
  %(prog)s --query '{job="injective-node"}' --hours 24 --format txt

  # 下载错误日志
  %(prog)s --query '{job="injective-node"} |~ "(ERR|FATAL|PANIC)"' --hours 24

  # 限制最大下载条目数
  %(prog)s --query '{job="injective-node"}' --hours 24 --max-entries 10000
        """
    )
    
    parser.add_argument(
        "--query",
        required=True,
        help="LogQL 查询表达式，例如: '{job=\"injective-node\"}'"
    )
    
    parser.add_argument(
        "--loki-url",
        default="http://localhost:3100",
        help="Loki 服务器地址 (默认: http://localhost:3100)"
    )
    
    time_group = parser.add_mutually_exclusive_group(required=True)
    time_group.add_argument(
        "--hours",
        type=int,
        help="下载最近 N 小时的日志"
    )
    time_group.add_argument(
        "--start",
        help="开始时间 (格式: 2024-01-01T00:00:00Z 或 2024-01-01 00:00:00)"
    )
    
    parser.add_argument(
        "--end",
        help="结束时间 (格式: 2024-01-01T00:00:00Z 或 2024-01-01 00:00:00，默认: 当前时间)"
    )
    
    parser.add_argument(
        "--format",
        choices=["json", "txt", "csv"],
        default="json",
        help="导出格式 (默认: json)"
    )
    
    parser.add_argument(
        "--output",
        help="输出文件路径 (默认: 根据格式和时间自动生成)"
    )
    
    parser.add_argument(
        "--limit-per-request",
        type=int,
        default=5000,
        help="每次请求的最大条目数 (默认: 5000，最大: 10000)"
    )
    
    parser.add_argument(
        "--max-entries",
        type=int,
        help="最大下载条目数 (默认: 不限制)"
    )
    
    parser.add_argument(
        "--timeout",
        type=int,
        default=300,
        help="请求超时时间（秒）(默认: 300)"
    )
    
    args = parser.parse_args()
    
    # 解析时间范围
    if args.hours:
        end_time = datetime.now()
        start_time = end_time - timedelta(hours=args.hours)
    else:
        start_time = parse_time(args.start)
        if args.end:
            end_time = parse_time(args.end)
        else:
            end_time = datetime.now()
    
    if start_time >= end_time:
        print("错误: 开始时间必须早于结束时间", file=sys.stderr)
        sys.exit(1)
    
    # 确定输出文件
    if args.output:
        output_file = args.output
    else:
        timestamp_str = datetime.now().strftime("%Y%m%d_%H%M%S")
        ext = {
            "json": "json",
            "txt": "txt",
            "csv": "csv"
        }[args.format]
        output_file = f"loki_logs_{timestamp_str}.{ext}"
    
    # 创建下载器并下载日志
    downloader = LokiLogDownloader(args.loki_url, args.timeout)
    
    try:
        logs = downloader.download_all_logs(
            args.query,
            start_time,
            end_time,
            args.limit_per_request,
            args.max_entries
        )
        
        if not logs:
            print("警告: 没有找到任何日志", file=sys.stderr)
            sys.exit(1)
        
        # 导出日志
        if args.format == "json":
            downloader.export_to_json(logs, output_file)
        elif args.format == "txt":
            downloader.export_to_text(logs, output_file)
        elif args.format == "csv":
            downloader.export_to_csv(logs, output_file)
        
        print(f"\n成功！日志已保存到: {output_file}")
        print(f"文件大小: {len(json.dumps(logs).encode('utf-8')) / 1024 / 1024:.2f} MB")
        
    except KeyboardInterrupt:
        print("\n\n下载被用户中断", file=sys.stderr)
        sys.exit(1)
    except Exception as e:
        print(f"\n错误: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
