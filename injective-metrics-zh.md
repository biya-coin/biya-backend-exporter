# Injective 节点指标中文说明文档

本文档对 Injective 区块链节点的 Prometheus 监控指标进行详细分类和说明。

## 目录
- [1. 区块链核心指标](#1-区块链核心指标)
- [2. 共识层指标](#2-共识层指标)
- [3. 网络层指标](#3-网络层指标)
- [4. 内存池指标](#4-内存池指标)
- [5. ABCI 接口指标](#5-abci-接口指标)
- [6. 存储层指标](#6-存储层指标)
- [7. 模块性能指标](#7-模块性能指标)
- [8. Go 运行时指标](#8-go-运行时指标)
- [9. 系统资源指标](#9-系统资源指标)
- [10. WebAssembly VM 指标](#10-webassembly-vm-指标)

---

## 1. 区块链核心指标

### 1.1 区块高度
- **cometbft_consensus_height**: 当前区块链高度
  - 类型: gauge
  - 用途: 监控区块链同步状态，检测节点是否在正常出块

- **cometbft_consensus_latest_block_height**: 最新区块高度
  - 类型: gauge
  - 用途: 与 consensus_height 类似，用于交叉验证

### 1.2 区块信息
- **cometbft_consensus_block_size_bytes**: 区块大小（字节）
  - 类型: gauge
  - 用途: 监控区块大小，评估网络负载

- **cometbft_consensus_num_txs**: 当前区块中的交易数量
  - 类型: gauge
  - 用途: 监控区块活跃度

- **cometbft_consensus_total_txs**: 总交易数
  - 类型: gauge
  - 用途: 统计链上总交易量

- **cometbft_consensus_chain_size_bytes**: 区块链总大小（字节）
  - 类型: counter
  - 用途: 监控链数据增长速度，评估存储需求

### 1.3 区块时间
- **cometbft_consensus_block_interval_seconds**: 区块间隔时间分布
  - 类型: histogram
  - 用途: 分析出块速度稳定性，检测异常延迟

---

## 2. 共识层指标

### 2.1 验证者信息
- **cometbft_consensus_validators**: 验证者总数
  - 类型: gauge
  - 用途: 监控验证者集合大小

- **cometbft_consensus_validators_power**: 所有验证者的总权重
  - 类型: gauge
  - 用途: 监控网络总质押量

- **cometbft_consensus_validator_power**: 单个验证者的权重
  - 类型: gauge
  - 标签: validator_address
  - 用途: 监控特定验证者的投票权重

- **cometbft_consensus_validator_last_signed_height**: 验证者最后签名的区块高度
  - 类型: gauge
  - 用途: 检测验证者活跃性

### 2.2 缺失和拜占庭验证者
- **cometbft_consensus_missing_validators**: 未签名的验证者数量
  - 类型: gauge
  - 用途: 监控网络健康度，检测离线验证者

- **cometbft_consensus_missing_validators_power**: 未签名验证者的总权重
  - 类型: gauge
  - 用途: 评估缺失验证者对网络安全的影响

- **cometbft_consensus_byzantine_validators**: 尝试双签的验证者数量
  - 类型: gauge
  - 用途: 检测恶意行为

- **cometbft_consensus_byzantine_validators_power**: 拜占庭验证者的总权重
  - 类型: gauge
  - 用途: 评估安全威胁程度

### 2.3 共识步骤
- **cometbft_consensus_rounds**: 共识轮数
  - 类型: gauge
  - 用途: 监控共识效率，检测共识卡住情况

- **cometbft_consensus_round_duration_seconds**: 共识轮持续时间分布
  - 类型: histogram
  - 用途: 分析共识性能

- **cometbft_consensus_step_duration_seconds**: 各个共识步骤的持续时间
  - 类型: histogram
  - 标签: step (Propose, Prevote, Precommit, Commit, NewHeight, NewRound)
  - 用途: 细粒度分析共识各阶段性能

### 2.4 投票相关
- **cometbft_consensus_round_voting_power_percent**: 当前轮次收到的投票权重百分比
  - 类型: gauge
  - 标签: vote_type (prevote, precommit)
  - 用途: 监控投票进度

- **cometbft_consensus_late_votes**: 延迟到达的投票数量
  - 类型: counter
  - 标签: vote_type
  - 用途: 检测网络延迟或同步问题

- **cometbft_consensus_duplicate_vote**: 收到的重复投票数量
  - 类型: counter
  - 用途: 监控网络冗余消息

### 2.5 区块和投票延迟
- **cometbft_consensus_full_prevote_delay**: prevote 全部完成的延迟时间
  - 类型: gauge
  - 标签: proposer_address
  - 用途: 分析提议者性能

- **cometbft_consensus_quorum_prevote_delay**: 达到 quorum 的 prevote 延迟
  - 类型: gauge
  - 标签: proposer_address
  - 用途: 监控共识效率

### 2.6 区块部分
- **cometbft_consensus_block_parts**: 各个节点传输的区块部分数量
  - 类型: counter
  - 标签: peer_id
  - 用途: 监控区块传播情况

- **cometbft_consensus_block_gossip_parts_received**: 收到的区块部分数量
  - 类型: counter
  - 标签: matches_current (true/false)
  - 用途: 监控区块同步效率

- **cometbft_consensus_duplicate_block_part**: 收到的重复区块部分
  - 类型: counter
  - 用途: 检测网络效率问题

### 2.7 提案相关
- **cometbft_consensus_proposal_create_count**: 节点创建的提案总数
  - 类型: counter
  - 用途: 监控节点作为提议者的活跃度

- **cometbft_consensus_proposal_receive_count**: 收到的提案总数
  - 类型: counter
  - 标签: status (accepted/rejected)
  - 用途: 监控提案接受率

- **cometbft_consensus_proposal_timestamp_difference**: 提案时间戳差异分布
  - 类型: histogram
  - 标签: is_timely (true/false)
  - 用途: 检测时钟偏移问题

### 2.8 区块同步
- **cometbft_blocksync_syncing**: 节点是否在进行区块同步
  - 类型: gauge
  - 值: 0 (未同步) / 1 (同步中)
  - 用途: 监控节点同步状态

---

## 3. 网络层指标

### 3.1 对等节点
- **cometbft_p2p_peers**: 连接的对等节点数量
  - 类型: gauge
  - 用途: 监控网络连接健康度

- **cometbft_p2p_peer_pending_send_bytes**: 待发送给对等节点的字节数
  - 类型: gauge
  - 标签: peer_id
  - 用途: 监控网络拥塞情况

### 3.2 消息传输
- **cometbft_p2p_message_send_bytes_total**: 发送的各类消息总字节数
  - 类型: counter
  - 标签: message_type
  - 消息类型包括:
    - v1_BlockPart: 区块部分
    - v1_Proposal: 提案
    - v1_Vote: 投票
    - v1_NewRoundStep: 新轮次步骤
    - v1_NewValidBlock: 新有效区块
    - v2_Txs: 交易
  - 用途: 分析网络流量分布

- **cometbft_p2p_message_receive_bytes_total**: 接收的各类消息总字节数
  - 类型: counter
  - 标签: message_type
  - 用途: 分析网络流量分布，检测异常流量

---

## 4. 内存池指标

### 4.1 内存池状态
- **cometbft_mempool_size**: 内存池中未提交的交易数量
  - 类型: gauge
  - 用途: 监控内存池负载（已弃用，建议使用 lane_size）

- **cometbft_mempool_size_bytes**: 内存池总大小（字节）
  - 类型: gauge
  - 用途: 监控内存占用（已弃用，建议使用 lane_bytes）

### 4.2 分道机制
- **cometbft_mempool_lane_size**: 每个通道的未提交交易数量
  - 类型: gauge
  - 标签: lane (default等)
  - 用途: 监控各通道负载

- **cometbft_mempool_lane_bytes**: 每个通道使用的字节数
  - 类型: gauge
  - 标签: lane
  - 用途: 监控通道内存占用

### 4.3 交易生命周期
- **cometbft_mempool_tx_life_span**: 交易在内存池中的停留时间（毫秒）
  - 类型: histogram
  - 标签: lane
  - 用途: 分析交易处理速度

- **cometbft_mempool_tx_size_bytes**: 交易大小分布（字节）
  - 类型: histogram
  - 用途: 分析交易大小特征

### 4.4 内存池操作
- **cometbft_mempool_already_received_txs**: 重复接收的交易数量
  - 类型: counter
  - 用途: 检测网络效率问题

- **cometbft_mempool_recheck_times**: 交易重新检查次数
  - 类型: counter
  - 用途: 监控内存池维护活动

### 4.5 连接管理
- **cometbft_mempool_active_outbound_connections**: 用于传播交易的活跃连接数
  - 类型: gauge
  - 用途: 监控交易传播能力

- **cometbft_mempool_disabled_routes**: 禁用的路由数量
  - 类型: gauge
  - 用途: 检测网络连接问题

- **cometbft_mempool_redundancy**: 冗余级别
  - 类型: gauge
  - 用途: 监控交易传播策略

---

## 5. ABCI 接口指标

### 5.1 方法调用时间
- **cometbft_abci_connection_method_timing_seconds**: 各 ABCI 方法的执行时间分布
  - 类型: histogram
  - 标签: method, type
  - 方法类型:
    - **init_chain**: 链初始化
    - **check_tx**: 交易检查（异步）
    - **prepare_proposal**: 准备提案（同步）
    - **process_proposal**: 处理提案（同步）
    - **finalize_block**: 完成区块（同步）
    - **commit**: 提交区块（同步）
    - **flush**: 刷新（同步）
    - **info**: 查询信息（同步）
  - 用途: 分析应用层性能瓶颈

---

## 6. 存储层指标

### 6.1 状态存储
- **cometbft_state_store_access_duration_seconds**: 状态存储访问时间
  - 类型: histogram
  - 标签: method
  - 方法类型:
    - **load**: 加载状态
    - **save**: 保存状态
    - **load_validators**: 加载验证者信息
    - **saveValidatorsInfo**: 保存验证者信息
    - **save_abci_responses**: 保存 ABCI 响应
  - 用途: 监控状态存储性能

### 6.2 区块存储
- **cometbft_store_block_store_access_duration_seconds**: 区块存储访问时间
  - 类型: histogram
  - 标签: method
  - 方法类型:
    - **new_block_store**: 初始化区块存储
    - **load_block_meta**: 加载区块元数据
    - **load_seen_commit**: 加载已见提交
    - **save_block**: 保存区块
    - **save_block_part**: 保存区块部分
    - **save_block_to_batch**: 批量保存区块
    - **save_bs_state**: 保存区块存储状态
  - 用途: 监控区块存储性能，识别 I/O 瓶颈

### 6.3 区块处理
- **cometbft_state_block_processing_time**: FinalizeBlock 处理时间
  - 类型: histogram
  - 用途: 监控区块处理性能

- **cometbft_state_consensus_param_updates**: 共识参数更新次数
  - 类型: counter
  - 用途: 监控链参数变更

---

## 7. 模块性能指标

### 7.1 BeginBlocker
- **begin_blocker**: 各模块 BeginBlocker 执行时间
  - 类型: summary
  - 标签: module
  - 模块类型:
    - **capability**: 能力模块
    - **distribution**: 分配模块（奖励分配）
    - **evidence**: 证据模块（处理双签等恶意行为）
    - **mint**: 铸币模块（通胀）
    - **slashing**: 惩罚模块
    - **staking**: 质押模块
    - **upgrade**: 升级模块
  - 用途: 识别 BeginBlock 阶段的性能瓶颈

### 7.2 EndBlocker
- **end_blocker**: 各模块 EndBlocker 执行时间
  - 类型: summary
  - 标签: module
  - 模块类型:
    - **crisis**: 危机模块
    - **gov**: 治理模块
    - **staking**: 质押模块（验证者更新）
  - 用途: 识别 EndBlock 阶段的性能瓶颈

---

## 8. Go 运行时指标

### 8.1 垃圾回收
- **go_gc_duration_seconds**: GC 暂停时间分布
  - 类型: summary
  - 用途: 监控 GC 对性能的影响

- **go_gc_gogc_percent**: GOGC 堆目标百分比
  - 类型: gauge
  - 用途: 查看 GC 配置

- **go_gc_gomemlimit_bytes**: Go 运行时内存限制
  - 类型: gauge
  - 用途: 监控内存限制设置

### 8.2 内存统计
- **go_memstats_alloc_bytes**: 当前分配的堆内存
  - 类型: gauge
  - 用途: 监控实时内存使用

- **go_memstats_alloc_bytes_total**: 累计分配的堆内存
  - 类型: counter
  - 用途: 分析内存分配速度

- **go_memstats_heap_alloc_bytes**: 堆分配的字节数
  - 类型: gauge
  - 用途: 监控堆使用情况

- **go_memstats_heap_idle_bytes**: 空闲堆内存
  - 类型: gauge
  - 用途: 评估内存利用率

- **go_memstats_heap_inuse_bytes**: 使用中的堆内存
  - 类型: gauge
  - 用途: 监控实际堆占用

- **go_memstats_heap_objects**: 堆对象数量
  - 类型: gauge
  - 用途: 监控对象分配情况

- **go_memstats_heap_released_bytes**: 释放给操作系统的堆内存
  - 类型: gauge
  - 用途: 监控内存回收效率

- **go_memstats_heap_sys_bytes**: 从系统获取的堆内存总量
  - 类型: gauge
  - 用途: 监控系统内存分配

- **go_memstats_next_gc_bytes**: 下次 GC 触发的堆大小阈值
  - 类型: gauge
  - 用途: 预测 GC 触发时机

### 8.3 内存分配统计
- **go_memstats_mallocs_total**: 累计分配的对象数
  - 类型: counter
  - 用途: 分析对象分配速度

- **go_memstats_frees_total**: 累计释放的对象数
  - 类型: counter
  - 用途: 分析对象释放速度

- **go_memstats_last_gc_time_seconds**: 最后一次 GC 的时间戳
  - 类型: gauge
  - 用途: 监控 GC 频率

### 8.4 其他内存区域
- **go_memstats_stack_inuse_bytes**: 栈内存使用量
  - 类型: gauge
  - 用途: 监控栈使用情况

- **go_memstats_stack_sys_bytes**: 栈内存系统分配量
  - 类型: gauge
  - 用途: 监控栈内存总量

- **go_memstats_sys_bytes**: 从系统获取的总内存
  - 类型: gauge
  - 用途: 监控总内存占用

- **go_memstats_mcache_inuse_bytes / go_memstats_mcache_sys_bytes**: mcache 内存统计
- **go_memstats_mspan_inuse_bytes / go_memstats_mspan_sys_bytes**: mspan 内存统计
- **go_memstats_buck_hash_sys_bytes**: 性能分析哈希表内存
- **go_memstats_gc_sys_bytes**: GC 元数据内存
- **go_memstats_other_sys_bytes**: 其他系统分配内存

### 8.5 协程和线程
- **go_goroutines**: 当前 goroutine 数量
  - 类型: gauge
  - 用途: 监控并发度，检测 goroutine 泄漏

- **go_threads**: 操作系统线程数
  - 类型: gauge
  - 用途: 监控线程使用情况

- **go_sched_gomaxprocs_threads**: GOMAXPROCS 设置
  - 类型: gauge
  - 用途: 查看并发配置

### 8.6 Go 版本信息
- **go_info**: Go 版本信息
  - 类型: gauge
  - 标签: version
  - 用途: 记录运行时版本

### 8.7 Runtime 指标（自定义）
- **runtime_alloc_bytes**: 运行时分配字节数
- **runtime_free_count**: 释放次数
- **runtime_heap_objects**: 堆对象数
- **runtime_malloc_count**: 分配次数
- **runtime_num_goroutines**: goroutine 数量
- **runtime_sys_bytes**: 系统字节数
- **runtime_total_gc_pause_ns**: 总 GC 暂停时间（纳秒）
- **runtime_total_gc_runs**: GC 运行总次数
- **runtime_gc_pause_ns**: GC 暂停时间分布

---

## 9. 系统资源指标

### 9.1 CPU
- **process_cpu_seconds_total**: 进程 CPU 总使用时间（用户态+系统态）
  - 类型: counter
  - 用途: 监控 CPU 使用情况

### 9.2 内存
- **process_resident_memory_bytes**: 常驻内存大小（RSS）
  - 类型: gauge
  - 用途: 监控物理内存占用

- **process_virtual_memory_bytes**: 虚拟内存大小
  - 类型: gauge
  - 用途: 监控虚拟内存占用

- **process_virtual_memory_max_bytes**: 最大虚拟内存限制
  - 类型: gauge
  - 用途: 查看系统限制

### 9.3 文件描述符
- **process_open_fds**: 当前打开的文件描述符数
  - 类型: gauge
  - 用途: 监控文件句柄使用情况，防止耗尽

- **process_max_fds**: 最大文件描述符数限制
  - 类型: gauge
  - 用途: 查看系统限制

### 9.4 网络流量
- **process_network_receive_bytes_total**: 进程接收的网络字节总数
  - 类型: counter
  - 用途: 监控网络入站流量

- **process_network_transmit_bytes_total**: 进程发送的网络字节总数
  - 类型: counter
  - 用途: 监控网络出站流量

### 9.5 进程启动时间
- **process_start_time_seconds**: 进程启动时间戳（Unix 时间）
  - 类型: gauge
  - 用途: 计算进程运行时长

---

## 10. WebAssembly VM 指标

### 10.1 缓存状态
- **wasmvm_cache_elements_total**: 缓存中的元素总数
  - 类型: gauge
  - 标签: type (memory, pinned)
  - 用途: 监控 WASM 缓存使用情况

- **wasmvm_cache_size_bytes**: 缓存大小（字节）
  - 类型: gauge
  - 标签: type (memory, pinned)
  - 用途: 监控 WASM 缓存内存占用

### 10.2 缓存命中率
- **wasmvm_cache_hits_total**: 缓存命中总数
  - 类型: counter
  - 标签: type (fs, memory, pinned)
  - 用途: 评估缓存效率

- **wasmvm_cache_misses_total**: 缓存未命中总数
  - 类型: counter
  - 用途: 评估缓存效率

---

## 11. Prometheus HTTP 处理器指标

- **promhttp_metric_handler_requests_in_flight**: 当前正在处理的抓取请求数
  - 类型: gauge
  - 用途: 监控 Prometheus 抓取并发度

- **promhttp_metric_handler_requests_total**: 按状态码分类的抓取请求总数
  - 类型: counter
  - 标签: code (200, 500, 503)
  - 用途: 监控指标端点健康状态

---

## 告警建议

### 高优先级告警

1. **节点同步状态异常**
   - 条件: `cometbft_blocksync_syncing == 1` 持续超过 30 分钟
   - 影响: 节点未与网络同步，无法参与共识

2. **验证者离线**
   - 条件: `cometbft_consensus_missing_validators > 0`
   - 影响: 可能导致惩罚和收益损失

3. **区块高度停滞**
   - 条件: `rate(cometbft_consensus_height[5m]) == 0`
   - 影响: 链停止出块

4. **内存使用过高**
   - 条件: `process_resident_memory_bytes > 阈值`
   - 影响: 可能导致 OOM

5. **文件描述符即将耗尽**
   - 条件: `process_open_fds / process_max_fds > 0.8`
   - 影响: 进程可能崩溃

### 中优先级告警

6. **区块处理时间过长**
   - 条件: `histogram_quantile(0.99, cometbft_state_block_processing_time_bucket) > 10s`
   - 影响: 影响出块效率

7. **内存池堆积**
   - 条件: `cometbft_mempool_size > 1000`
   - 影响: 交易处理延迟

8. **GC 暂停时间过长**
   - 条件: `go_gc_duration_seconds{quantile="0.99"} > 0.001`
   - 影响: 性能抖动

9. **网络连接数异常**
   - 条件: `cometbft_p2p_peers < 3`
   - 影响: 网络孤立风险

10. **共识轮数异常**
    - 条件: `cometbft_consensus_rounds > 0`
    - 影响: 共识效率降低

### 低优先级告警

11. **Goroutine 泄漏**
    - 条件: `go_goroutines > 500` 且持续增长
    - 影响: 内存泄漏风险

12. **存储访问延迟**
    - 条件: `histogram_quantile(0.99, cometbft_store_block_store_access_duration_seconds_bucket) > 0.1`
    - 影响: I/O 性能下降

---

## 仪表盘建议

### 1. 区块链概览仪表盘
- 当前区块高度
- 区块间隔时间
- 验证者数量和权重
- 交易总数和区块大小
- 网络连接数

### 2. 共识性能仪表盘
- 各共识步骤耗时
- 投票完成率
- 提案创建/接收统计
- 区块部分传输情况

### 3. 系统资源仪表盘
- CPU 使用率
- 内存使用（RSS 和虚拟内存）
- 文件描述符使用率
- 网络流量（入站/出站）
- Goroutine 数量

### 4. 存储性能仪表盘
- 状态存储访问延迟
- 区块存储访问延迟
- 各存储操作的 P99 延迟

### 5. 应用层仪表盘
- ABCI 方法调用延迟
- BeginBlocker/EndBlocker 各模块耗时
- 区块处理时间分布

### 6. 内存池仪表盘
- 内存池大小和字节数
- 交易生命周期分布
- 交易大小分布
- 重复交易统计

---

## 性能优化建议

根据这些指标，可以进行以下优化：

1. **降低区块处理时间**
   - 监控 `cometbft_state_block_processing_time`
   - 优化慢模块（通过 `begin_blocker` 和 `end_blocker` 识别）

2. **减少存储延迟**
   - 监控 `cometbft_store_block_store_access_duration_seconds`
   - 使用 SSD、调整数据库配置、启用缓存

3. **优化内存使用**
   - 监控 `go_memstats_heap_alloc_bytes` 和 `process_resident_memory_bytes`
   - 调整 GOGC 参数、减少对象分配

4. **改善网络性能**
   - 监控 `cometbft_p2p_message_*_bytes_total`
   - 优化消息大小、增加带宽、调整连接数

5. **提升共识效率**
   - 监控 `cometbft_consensus_step_duration_seconds`
   - 优化时钟同步、减少网络延迟、调整超时参数

---

## 总结

本文档涵盖了 Injective 节点暴露的所有 Prometheus 指标，按功能模块进行了详细分类和说明。这些指标对于监控节点健康状态、识别性能瓶颈、设置告警规则至关重要。

建议运维人员根据实际情况定制监控仪表盘和告警策略，持续优化节点性能和稳定性。
