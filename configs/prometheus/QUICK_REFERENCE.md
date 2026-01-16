# 告警规则快速参考卡片

## 核心告警规则一览

### 🔴 紧急告警（Emergency）

| 告警名称 | 触发条件 | 持续时间 | 说明 |
|---------|---------|---------|------|
| **PerformanceComprehensiveAbnormal** | TPS<50%历史均值 AND 延迟>10s AND 成功率<98% | 3分钟 | 多指标同时异常 |

**处理流程**：
1. 立即启动应急预案
2. 召集技术团队会商
3. 全面排查网络问题
4. 必要时暂停服务维护

---

### 🟠 严重告警（Critical）

| 告警名称 | 触发条件 | 持续时间 | 说明 |
|---------|---------|---------|------|
| **TransactionSuccessRateDrop** | 成功率 < 98% | 5分钟 | 交易成功率下降 |
| **TPSZero** | TPS = 0 | 2分钟 | 链可能停止出块 |

**处理流程**：
1. 立即查看失败交易原因
2. 检查是否有合约攻击
3. 分析失败交易的共同特征
4. 紧急处理高危问题

---

### 🟡 警告告警（Warning）

| 告警名称 | 触发条件 | 持续时间 | 说明 |
|---------|---------|---------|------|
| **TPSAbnormalDrop** | TPS < 50%历史平均值 | 5分钟 | TPS异常下降 |
| **NetworkLatencyAbnormalIncrease** | 延迟 > 10秒 | 10分钟 | 网络延迟异常 |
| **MempoolCongestion** | 待处理交易 > 10000 | 5分钟 | 交易池拥堵 |
| **HighGasUtilization** | Gas利用率 > 95% | 10分钟 | Gas利用率过高 |
| **BlockTimeAbnormal** | 出块时间 > 10秒 | 5分钟 | 出块时间异常 |
| **NetworkCongestion** | 拥堵指数 > 80% | 5分钟 | 网络拥堵 |

**处理流程**：
- 检查节点运行状态
- 查看网络分区情况
- 分析交易来源
- 评估是否需要扩容

---

## 快速命令

### 查看告警状态

```bash
# 查看所有活跃告警
curl http://localhost:9090/api/v1/alerts | jq '.data.alerts[] | select(.state=="firing")'

# 查看特定告警
curl "http://localhost:9090/api/v1/alerts" | jq '.data.alerts[] | select(.labels.alertname=="TPSAbnormalDrop")'

# 使用测试脚本
./scripts/test_alerts.sh
```

### 静默告警

```bash
# 静默 2 小时
amtool silence add alertname=TPSAbnormalDrop --duration=2h --comment="维护窗口"

# 查看所有静默
amtool silence query

# 删除静默
amtool silence expire <silence-id>
```

### 测试告警

```bash
# 验证规则语法
promtool check rules configs/prometheus/alert_rules.yml

# 测试 PromQL 表达式
curl "http://localhost:9090/api/v1/query?query=biya_tps_current"

# 重新加载配置
curl -X POST http://localhost:9090/-/reload
```

---

## 关键指标快查

| 指标名称 | 类型 | 说明 | 正常范围 |
|---------|------|------|---------|
| `biya_tps_current` | Gauge | 当前TPS | > 历史均值的50% |
| `biya_tps_24h_avg` | Gauge | 24小时平均TPS | - |
| `biya_tx_success_rate` | Gauge | 交易成功率 | > 0.98 (98%) |
| `biya_tx_confirm_time_seconds_avg` | Gauge | 平均确认延迟 | < 10秒 |
| `biya_chain_mempool_pending_txs` | Gauge | 待处理交易数 | < 10000 |
| `biya_chain_block_gas_utilization_ratio_avg` | Gauge | Gas利用率 | < 0.95 (95%) |
| `biya_chain_congestion_ratio` | Gauge | 拥堵指数 | < 0.8 (80%) |
| `biya_avg_block_time` | Gauge | 平均出块时间 | < 10秒 |

---

## 告警级别说明

| 级别 | 响应时间 | 通知方式 | 典型场景 |
|------|---------|---------|---------|
| 🔴 **Emergency** | 立即 | 电话+短信+IM+邮件 | 系统级故障 |
| 🟠 **Critical** | < 5分钟 | 短信+IM+邮件 | 关键功能异常 |
| 🟡 **Warning** | < 30分钟 | IM+邮件 | 性能下降 |
| 🔵 **Info** | < 2小时 | 邮件 | 一般性通知 |

---

## 常用 PromQL 表达式

```promql
# 查询当前 TPS
biya_tps_current

# 查询 7 天平均 TPS
avg_over_time(biya_tps_24h_avg[7d])

# 计算 TPS 下降百分比
((avg_over_time(biya_tps_24h_avg[7d]) - biya_tps_current) / avg_over_time(biya_tps_24h_avg[7d])) * 100

# 查询过去 1 小时的成功率
avg_over_time(biya_tx_success_rate[1h])

# 查询过去 5 分钟的交易增长率
rate(biya_tx_total[5m])

# 检查指标是否缺失
absent(biya_tps_current)
```

---

## Web UI 访问

| 服务 | 地址 | 说明 |
|------|------|------|
| **Prometheus** | http://localhost:9090 | 主界面 |
| **告警页面** | http://localhost:9090/alerts | 查看所有告警规则和状态 |
| **规则页面** | http://localhost:9090/rules | 查看所有规则（告警+记录） |
| **目标页面** | http://localhost:9090/targets | 查看采集目标状态 |
| **Alertmanager** | http://localhost:9093 | 告警管理和静默 |

---

## 故障排查检查清单

### ❌ 告警未触发

- [ ] 表达式是否正确？→ 在 Prometheus UI 中测试
- [ ] 数据是否存在？→ 查询基础指标
- [ ] 持续时间是否太长？→ 检查 `for` 字段
- [ ] 规则是否已加载？→ 访问 `/rules` 页面
- [ ] Alertmanager 是否连接？→ 访问 `/alertmanagers` 页面

### ❌ 通知未收到

- [ ] Alertmanager 是否运行？→ `curl http://localhost:9093/-/healthy`
- [ ] 路由配置是否正确？→ 检查 `alertmanager.yml` 中的 `route`
- [ ] 接收者配置是否正确？→ 检查 `receivers` 配置
- [ ] 是否被静默？→ `amtool silence query`
- [ ] 是否被抑制？→ 检查 `inhibit_rules`

### ❌ 告警太多（告警疲劳）

- [ ] 阈值是否过于敏感？→ 调整表达式中的阈值
- [ ] 持续时间是否太短？→ 增加 `for` 时间
- [ ] 是否需要聚合？→ 配置告警分组
- [ ] 重复间隔是否太短？→ 调整 `repeat_interval`

---

## 联系方式

| 场景 | 联系方式 |
|------|---------|
| 紧急故障 | oncall@biya.chain |
| 运维支持 | ops@biya.chain |
| 技术支持 | support@biya.chain |

---

## 相关文档

- 📖 [完整配置文档](./README.md)
- 📖 [告警规则编写指南](./ALERT_RULES_GUIDE.md)
- 📖 [变更日志](../../CHANGELOG_ALERTS.md)
- 📖 [产品需求文档](../../公链后台需求/公链后台管理系统需求文档%20-%20运行状态监控.md)

---

**最后更新**: 2025-12-24

