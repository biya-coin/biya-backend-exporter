# Injective 链指标说明（中文版）

本文档对 Injective 节点暴露的 Prometheus 指标进行归类与中文说明。指标来源于 CometBFT（原 Tendermint）共识引擎及 Go 运行时等。

---

## 一、ABCI 连接与调用

与应用层（ABCI）交互的耗时与调用情况。

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_abci_connection_method_timing_seconds` | histogram | 各 ABCI 方法的调用耗时（秒）。标签：`method`（如 commit、finalize_block、flush、info、prepare_proposal、process_proposal）、`type`（如 sync）。用于分析应用层处理延迟。 |

**常见 method 含义：**
- **commit**：提交区块
- **finalize_block**：敲定区块
- **flush**：刷新连接
- **info**：节点信息查询
- **prepare_proposal**：准备提案（提议者）
- **process_proposal**：处理提案（验证者）

---

## 二、共识（Consensus）

区块生产、投票、验证者与链高相关的指标。

### 2.1 区块与链

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_consensus_block_interval_seconds` | histogram | 当前区块与上一区块之间的时间间隔（秒）。可用来观察出块节奏与稳定性。 |
| `cometbft_consensus_block_size_bytes` | gauge | 当前区块大小（字节）。 |
| `cometbft_consensus_block_gossip_parts_received` | counter | 节点收到的区块部分数量。标签 `matches_current` 表示是否与当前正在收集的区块相关。 |
| `cometbft_consensus_chain_size_bytes` | counter | 链的累计大小（字节）。 |
| `cometbft_consensus_height` | gauge | 当前共识的链高度。 |
| `cometbft_consensus_latest_block_height` | gauge | 最新区块高度。 |
| `cometbft_consensus_num_txs` | gauge | 当前区块中的交易数量。 |
| `cometbft_consensus_total_txs` | gauge | 总交易数（当前上下文中）。 |

### 2.2 提案（Proposal）

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_consensus_proposal_create_count` | counter | 自进程启动以来本节点**创建**的提案总数。可带状态标签（如 accepted/rejected）。 |
| `cometbft_consensus_proposal_receive_count` | counter | 自进程启动以来本节点**收到**的提案总数。标签 `status`：accepted/rejected。 |
| `cometbft_consensus_proposal_timestamp_difference` | histogram | 收到提案消息时的本地时间与提案消息中时间戳的差值（秒）。标签 `is_timely` 表示是否在时间窗内。用于监控提案时效性。 |

### 2.3 投票与轮次

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_consensus_round_duration_seconds` | histogram | 共识轮次持续时长（秒）。 |
| `cometbft_consensus_round_voting_power_percent` | gauge | 当前轮次内已收到的投票权占全体投票权的百分比。标签 `vote_type`：precommit、prevote。 |
| `cometbft_consensus_rounds` | gauge | 当前高度下的共识轮次数（0 表示首轮即达成）。 |
| `cometbft_consensus_step_duration_seconds` | histogram | 共识各步骤的耗时（秒）。标签 `step`：NewHeight、NewRound、Propose、Prevote、Precommit、Commit。 |
| `cometbft_consensus_full_prevote_delay` | gauge | 从提案时间戳到「所有验证者都完成预投票」的那一票的时间间隔（秒）。按 `proposer_address` 区分。 |
| `cometbft_consensus_quorum_prevote_delay` | gauge | 从提案时间戳到「达到法定人数预投票」的最早那一票的时间间隔（秒）。按 `proposer_address` 区分。 |

### 2.4 验证者与拜占庭

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_consensus_validators` | gauge | 验证者数量。 |
| `cometbft_consensus_validators_power` | gauge | 所有验证者的总投票权。 |
| `cometbft_consensus_validator_power` | gauge | 单个验证者的投票权。标签 `validator_address`。 |
| `cometbft_consensus_validator_last_signed_height` | gauge | 若本节点是验证者，其最后一次签名的区块高度。标签 `validator_address`。 |
| `cometbft_consensus_missing_validators` | gauge | 未在本轮签名的验证者数量。 |
| `cometbft_consensus_missing_validators_power` | gauge | 未签名验证者的总投票权。 |
| `cometbft_consensus_byzantine_validators` | gauge | 被判定为双签等拜占庭行为的验证者数量。 |
| `cometbft_consensus_byzantine_validators_power` | gauge | 拜占庭验证者的总投票权。 |

---

## 三、内存池（Mempool）

待打包的未确认交易与内存占用。

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_mempool_size` | gauge | 内存池中未提交交易总数。（已弃用：可由各 lane 的 size 求和得到。） |
| `cometbft_mempool_size_bytes` | gauge | 内存池总占用字节数。（已弃用：可由各 lane 的 bytes 求和得到。） |
| `cometbft_mempool_lane_size` | gauge | 每个通道（lane）的未提交交易数。标签 `lane`（如 default）。 |
| `cometbft_mempool_lane_bytes` | gauge | 每个通道已使用的字节数。标签 `lane`。 |

---

## 四、状态与区块处理

应用状态更新与区块处理耗时。

### 4.1 状态（State）

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_state_block_processing_time` | histogram | 处理 FinalizeBlock（敲定区块）所花费的时间。单位与实现相关（常为秒或毫秒）。 |
| `cometbft_state_consensus_param_updates` | counter | 自进程启动以来，应用返回的共识参数更新次数。 |
| `cometbft_state_store_access_duration_seconds` | histogram | 访问状态存储的耗时（秒）。标签 `method`：load、save、load_abci_responses、save_abci_responses、load_validators、saveValidatorsInfo 等。 |

### 4.2 区块存储（Block Store）

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `cometbft_store_block_store_access_duration_seconds` | histogram | 访问区块存储的耗时（秒）。标签 `method`：load_block、load_block_meta、load_block_part、load_seen_commit、save_block、save_block_part、save_block_to_batch、save_bs_state、new_block_store 等。 |

---

## 五、Go 运行时

Go 进程的 GC、调度与内存统计。

### 5.1 垃圾回收（GC）

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `go_gc_duration_seconds` | summary | 每次 GC 的停顿时间（秒），即 stop-the-world 时长。 |
| `go_gc_gogc_percent` | gauge | 用户配置的堆目标百分比（如 GOGC 或 runtime/debug.SetGCPercent）。 |
| `go_gc_gomemlimit_bytes` | gauge | 用户配置的 Go 内存上限（如 GOMEMLIMIT），未设置时为极大值。 |

### 5.2 Goroutine 与线程

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `go_goroutines` | gauge | 当前存在的 goroutine 数量。 |
| `go_threads` | gauge | 已创建的 OS 线程数。 |
| `go_sched_gomaxprocs_threads` | gauge | 当前 GOMAXPROCS，即可同时执行用户 Go 代码的 OS 线程数。 |
| `go_info` | gauge | Go 环境信息。标签如 `version`（如 go1.23.9）。 |

### 5.3 内存（memstats）

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `go_memstats_alloc_bytes` | gauge | 堆上当前已分配且仍在使用的字节数。 |
| `go_memstats_alloc_bytes_total` | counter | 堆上累计分配过的字节数（含已释放）。 |
| `go_memstats_heap_alloc_bytes` | gauge | 堆上已分配且在用的字节数（与 alloc_bytes 一致）。 |
| `go_memstats_heap_idle_bytes` | gauge | 堆上空闲、待使用的字节数。 |
| `go_memstats_heap_inuse_bytes` | gauge | 堆上正在使用的字节数。 |
| `go_memstats_heap_sys_bytes` | gauge | 从系统申请的堆总字节数。 |
| `go_memstats_heap_objects` | gauge | 当前堆上存活对象数量。 |
| `go_memstats_heap_released_bytes` | gauge | 已归还给操作系统的堆字节数。 |
| `go_memstats_frees_total` | counter | 堆对象累计释放次数。 |
| `go_memstats_mallocs_total` | counter | 堆对象累计分配次数（含已被 GC 回收的）。 |
| `go_memstats_next_gc_bytes` | gauge | 触发下一次 GC 的堆大小目标（字节）。 |
| `go_memstats_last_gc_time_seconds` | gauge | 上次 GC 的 Unix 时间戳（秒）。 |
| `go_memstats_gc_sys_bytes` | gauge | GC 元数据占用的系统内存（字节）。 |
| `go_memstats_buck_hash_sys_bytes` | gauge |  profiling 桶哈希表占用字节数。 |
| `go_memstats_mcache_inuse_bytes` | gauge | mcache 正在使用的字节数。 |
| `go_memstats_mcache_sys_bytes` | gauge | mcache 从系统申请的字节数。 |
| `go_memstats_mspan_inuse_bytes` | gauge | mspan 正在使用的字节数。 |
| `go_memstats_mspan_sys_bytes` | gauge | mspan 从系统申请的字节数。 |
| `go_memstats_stack_inuse_bytes` | gauge | 栈分配器正在使用的字节数。 |
| `go_memstats_stack_sys_bytes` | gauge | 栈分配器从系统申请的字节数。 |
| `go_memstats_other_sys_bytes` | gauge | 其他系统分配占用的字节数。 |
| `go_memstats_sys_bytes` | gauge | 从系统申请的总字节数。 |

---

## 六、进程与系统资源

节点进程的 CPU、内存、网络与文件描述符。

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `process_cpu_seconds_total` | counter | 进程累计消耗的 CPU 时间（秒，含用户态+内核态）。 |
| `process_resident_memory_bytes` | gauge | 进程常驻内存（RSS）字节数。 |
| `process_virtual_memory_bytes` | gauge | 进程虚拟内存占用字节数。 |
| `process_virtual_memory_max_bytes` | gauge | 进程可用的最大虚拟内存（字节）。 |
| `process_start_time_seconds` | gauge | 进程启动时的 Unix 时间戳（秒）。 |
| `process_open_fds` | gauge | 当前打开的文件描述符数量。 |
| `process_max_fds` | gauge | 进程可打开的文件描述符上限。 |
| `process_network_receive_bytes_total` | counter | 进程累计接收的网络字节数。 |
| `process_network_transmit_bytes_total` | counter | 进程累计发送的网络字节数。 |

---

## 七、Prometheus HTTP 暴露

指标 HTTP 接口的请求情况。

| 指标名 | 类型 | 说明 |
|--------|------|------|
| `promhttp_metric_handler_requests_in_flight` | gauge | 当前正在被处理的抓取（scrape）请求数。 |
| `promhttp_metric_handler_requests_total` | counter | 按 HTTP 状态码统计的抓取总次数。标签 `code`（如 200、500、503）。 |

---

## 归类总览

| 类别 | 含义 |
|------|------|
| **ABCI** | 与应用的调用延迟与方法分布。 |
| **共识** | 区块、提案、投票、验证者、链高与拜占庭检测。 |
| **内存池** | 待打包交易数量与内存占用。 |
| **状态 / 存储** | 状态库与区块库的访问延迟、区块处理时间。 |
| **Go 运行时** | GC、goroutine、线程、堆/栈与系统内存。 |
| **进程** | CPU、内存、网络、文件描述符与启动时间。 |
| **HTTP** | 指标接口的并发与请求统计。 |

---

## 常用标签说明

- **chain_id**：链 ID（如 `injective-666` 表示测试网）。  
- **method**：ABCI 方法名或 store 方法名。  
- **type**：连接类型（如 sync）。  
- **step**：共识步骤（NewHeight、NewRound、Propose、Prevote、Precommit、Commit）。  
- **vote_type**：投票类型（prevote、precommit）。  
- **lane**：内存池通道（如 default）。  
- **validator_address** / **proposer_address**：验证者/提议者地址。  
- **status** / **is_timely**：提案状态或是否及时。  
- **le**：直方图桶上界（≤）。  
- **code**：HTTP 状态码。

如需对某一类指标做告警或大盘，可优先从「共识」「ABCI」「状态/存储」和「进程」四类入手。
