# HELP begin_blocker begin_blocker
# TYPE begin_blocker summary
begin_blocker{module="capability",quantile="0.5"} 0.007153000216931105
begin_blocker{module="capability",quantile="0.9"} 0.007422999944537878
begin_blocker{module="capability",quantile="0.99"} 0.007422999944537878
begin_blocker_sum{module="capability"} 560.8019250514917
begin_blocker_count{module="capability"} 56693
begin_blocker{module="distribution",quantile="0.5"} 0.33373698592185974
begin_blocker{module="distribution",quantile="0.9"} 0.3398439884185791
begin_blocker{module="distribution",quantile="0.99"} 0.3398439884185791
begin_blocker_sum{module="distribution"} 21249.42687267065
begin_blocker_count{module="distribution"} 56693
begin_blocker{module="evidence",quantile="0.5"} 0.002004999900236726
begin_blocker{module="evidence",quantile="0.9"} 0.0022559999488294125
begin_blocker{module="evidence",quantile="0.99"} 0.0022559999488294125
begin_blocker_sum{module="evidence"} 143.07146200269926
begin_blocker_count{module="evidence"} 56693
begin_blocker{module="mint",quantile="0.5"} 0.281823992729187
begin_blocker{module="mint",quantile="0.9"} 0.29623299837112427
begin_blocker{module="mint",quantile="0.99"} 0.29623299837112427
begin_blocker_sum{module="mint"} 19419.40085528791
begin_blocker_count{module="mint"} 56693
begin_blocker{module="slashing",quantile="0.5"} 0.10716599971055984
begin_blocker{module="slashing",quantile="0.9"} 0.12193399667739868
begin_blocker{module="slashing",quantile="0.99"} 0.12193399667739868
begin_blocker_sum{module="slashing"} 6979.286821954709
begin_blocker_count{module="slashing"} 56693
begin_blocker{module="staking",quantile="0.5"} 0.20178300142288208
begin_blocker{module="staking",quantile="0.9"} 0.2094620019197464
begin_blocker{module="staking",quantile="0.99"} 0.2094620019197464
begin_blocker_sum{module="staking"} 12161.243859447539
begin_blocker_count{module="staking"} 56693
begin_blocker{module="upgrade",quantile="0.5"} 0.009522000327706337
begin_blocker{module="upgrade",quantile="0.9"} 0.012569000013172626
begin_blocker{module="upgrade",quantile="0.99"} 0.012569000013172626
begin_blocker_sum{module="upgrade"} 682.8920528916642
begin_blocker_count{module="upgrade"} 56693
# HELP cometbft_abci_connection_method_timing_seconds Timing for each ABCI method.
# TYPE cometbft_abci_connection_method_timing_seconds histogram
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="0.0001"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="0.0004"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="0.002"} 24
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="0.009"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="0.02"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="0.1"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="0.65"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="2"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="6"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="25"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="check_tx",type="async",le="+Inf"} 35
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="check_tx",type="async"} 0.06850520799999998
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="check_tx",type="async"} 35
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="0.0001"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="0.0004"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="0.002"} 657
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="0.009"} 56680
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="0.02"} 56687
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="0.1"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="0.65"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="2"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="6"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="25"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="commit",type="sync",le="+Inf"} 56693
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="commit",type="sync"} 151.22977754900074
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="commit",type="sync"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="0.0001"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="0.0004"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="0.002"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="0.009"} 56552
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="0.02"} 56690
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="0.1"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="0.65"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="2"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="6"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="25"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="finalize_block",type="sync",le="+Inf"} 56693
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="finalize_block",type="sync"} 307.3705516699999
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="finalize_block",type="sync"} 56693
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="0.0001"} 56694
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="0.0004"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="0.002"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="0.009"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="0.02"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="0.1"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="0.65"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="2"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="6"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="25"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="flush",type="sync",le="+Inf"} 56695
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="flush",type="sync"} 0.037293778999999284
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="flush",type="sync"} 56695
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="0.0001"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="0.0004"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="0.002"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="0.009"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="0.02"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="0.1"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="0.65"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="2"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="6"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="25"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="info",type="sync",le="+Inf"} 1
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="info",type="sync"} 8.68e-06
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="info",type="sync"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="0.0001"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="0.0004"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="0.002"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="0.009"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="0.02"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="0.1"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="0.65"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="2"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="6"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="25"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="init_chain",type="sync",le="+Inf"} 1
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="init_chain",type="sync"} 0.625931164
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="init_chain",type="sync"} 1
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="0.0001"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="0.0004"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="0.002"} 14139
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="0.009"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="0.02"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="0.1"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="0.65"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="2"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="6"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="25"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="prepare_proposal",type="sync",le="+Inf"} 14142
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="prepare_proposal",type="sync"} 9.227299544
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="prepare_proposal",type="sync"} 14142
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="0.0001"} 0
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="0.0004"} 9644
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="0.002"} 56690
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="0.009"} 56692
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="0.02"} 56692
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="0.1"} 56692
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="0.65"} 56692
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="2"} 56692
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="6"} 56692
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="25"} 56692
cometbft_abci_connection_method_timing_seconds_bucket{chain_id="biyachain-1",method="process_proposal",type="sync",le="+Inf"} 56692
cometbft_abci_connection_method_timing_seconds_sum{chain_id="biyachain-1",method="process_proposal",type="sync"} 26.40013299599982
cometbft_abci_connection_method_timing_seconds_count{chain_id="biyachain-1",method="process_proposal",type="sync"} 56692
# HELP cometbft_blocksync_syncing Whether or not a node is block syncing. 1 if yes, 0 if no.
# TYPE cometbft_blocksync_syncing gauge
cometbft_blocksync_syncing{chain_id="biyachain-1"} 0
# HELP cometbft_consensus_block_gossip_parts_received Number of block parts received by the node, separated by whether the part was relevant to the block the node is trying to gather or not.
# TYPE cometbft_consensus_block_gossip_parts_received counter
cometbft_consensus_block_gossip_parts_received{chain_id="biyachain-1",matches_current="false"} 651
cometbft_consensus_block_gossip_parts_received{chain_id="biyachain-1",matches_current="true"} 223482
# HELP cometbft_consensus_block_interval_seconds Time between this and the last block.
# TYPE cometbft_consensus_block_interval_seconds histogram
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="0.005"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="0.01"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="0.025"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="0.05"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="0.1"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="0.25"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="0.5"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="1"} 0
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="2.5"} 56689
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="5"} 56691
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="10"} 56691
cometbft_consensus_block_interval_seconds_bucket{chain_id="biyachain-1",le="+Inf"} 56692
cometbft_consensus_block_interval_seconds_sum{chain_id="biyachain-1"} 1.9057151272232267e+06
cometbft_consensus_block_interval_seconds_count{chain_id="biyachain-1"} 56692
# HELP cometbft_consensus_block_parts Number of block parts transmitted by each peer.
# TYPE cometbft_consensus_block_parts counter
cometbft_consensus_block_parts{chain_id="biyachain-1",peer_id="5fc8be4de7f815dd1f2aa607ec9f5b26ef969a7e"} 42118
cometbft_consensus_block_parts{chain_id="biyachain-1",peer_id="741475b39fd741c82b2373d9a3742230a4fc70a2"} 41870
cometbft_consensus_block_parts{chain_id="biyachain-1",peer_id="90ab7585f3f4faed99bd05a5fb663ce82742ef51"} 41882
cometbft_consensus_block_parts{chain_id="biyachain-1",peer_id="a95b67eceaebf1445060ee7fe09523ca7e6cbafc"} 42088
cometbft_consensus_block_parts{chain_id="biyachain-1",peer_id="d392f2adcc3713ec863a76477be7701fb37d8a3a"} 42033
# HELP cometbft_consensus_block_size_bytes Size of the block.
# TYPE cometbft_consensus_block_size_bytes gauge
cometbft_consensus_block_size_bytes{chain_id="biyachain-1"} 910
# HELP cometbft_consensus_byzantine_validators Number of validators who tried to double sign.
# TYPE cometbft_consensus_byzantine_validators gauge
cometbft_consensus_byzantine_validators{chain_id="biyachain-1"} 0
# HELP cometbft_consensus_byzantine_validators_power Total power of the byzantine validators.
# TYPE cometbft_consensus_byzantine_validators_power gauge
cometbft_consensus_byzantine_validators_power{chain_id="biyachain-1"} 0
# HELP cometbft_consensus_chain_size_bytes Size of the chain in bytes.
# TYPE cometbft_consensus_chain_size_bytes counter
cometbft_consensus_chain_size_bytes{chain_id="biyachain-1"} 5.1811062e+07
# HELP cometbft_consensus_duplicate_block_part Number of times we received a duplicate block part
# TYPE cometbft_consensus_duplicate_block_part counter
cometbft_consensus_duplicate_block_part{chain_id="biyachain-1"} 166789
# HELP cometbft_consensus_duplicate_vote Number of times we received a duplicate vote
# TYPE cometbft_consensus_duplicate_vote counter
cometbft_consensus_duplicate_vote{chain_id="biyachain-1"} 1.004934e+06
# HELP cometbft_consensus_full_prevote_delay Interval in seconds between the proposal timestamp and the timestamp of the latest prevote in a round where all validators voted.
# TYPE cometbft_consensus_full_prevote_delay gauge
cometbft_consensus_full_prevote_delay{chain_id="biyachain-1",proposer_address="3210F89D8743393559433C7F9D78ECA3996B8117"} 1.708570495
cometbft_consensus_full_prevote_delay{chain_id="biyachain-1",proposer_address="7764F2E6B8CF70E8BF8BE7C8B48595DCE483C6EF"} 1.7101279169999999
cometbft_consensus_full_prevote_delay{chain_id="biyachain-1",proposer_address="9D3D5D721115DF093D7D1B7B1C4ED3ED8DA886F4"} 1.710862842
cometbft_consensus_full_prevote_delay{chain_id="biyachain-1",proposer_address="DF2AC1640CF3B4CFBEECE6191F1F32341CF72CFE"} 1.710949409
# HELP cometbft_consensus_height Height of the chain.
# TYPE cometbft_consensus_height gauge
cometbft_consensus_height{chain_id="biyachain-1"} 508435
# HELP cometbft_consensus_late_votes LateVotes stores the number of votes that were received by this node that correspond to earlier heights and rounds than this node is currently in.
# TYPE cometbft_consensus_late_votes counter
cometbft_consensus_late_votes{chain_id="biyachain-1",vote_type="precommit"} 651092
cometbft_consensus_late_votes{chain_id="biyachain-1",vote_type="prevote"} 216647
# HELP cometbft_consensus_latest_block_height The latest block height.
# TYPE cometbft_consensus_latest_block_height gauge
cometbft_consensus_latest_block_height{chain_id="biyachain-1"} 508434
# HELP cometbft_consensus_missing_validators Number of validators who did not sign.
# TYPE cometbft_consensus_missing_validators gauge
cometbft_consensus_missing_validators{chain_id="biyachain-1"} 0
# HELP cometbft_consensus_missing_validators_power Total power of the missing validators.
# TYPE cometbft_consensus_missing_validators_power gauge
cometbft_consensus_missing_validators_power{chain_id="biyachain-1"} 0
# HELP cometbft_consensus_num_txs Number of transactions.
# TYPE cometbft_consensus_num_txs gauge
cometbft_consensus_num_txs{chain_id="biyachain-1"} 0
# HELP cometbft_consensus_proposal_create_count ProposalCreationCount is the total number of proposals created by this node since process start. The metric is annotated by the status of the proposal from the application, either 'accepted' or 'rejected'.
# TYPE cometbft_consensus_proposal_create_count counter
cometbft_consensus_proposal_create_count{chain_id="biyachain-1"} 14142
# HELP cometbft_consensus_proposal_receive_count ProposalReceiveCount is the total number of proposals received by this node since process start. The metric is annotated by the status of the proposal from the application, either 'accepted' or 'rejected'.
# TYPE cometbft_consensus_proposal_receive_count counter
cometbft_consensus_proposal_receive_count{chain_id="biyachain-1",status="accepted"} 56692
# HELP cometbft_consensus_proposal_timestamp_difference Difference in seconds between the local time when a proposal message is received and the timestamp in the proposal message.
# TYPE cometbft_consensus_proposal_timestamp_difference histogram
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="-1.5"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="-1"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="-0.5"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="-0.2"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="0"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="0.2"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="0.5"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="1"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="1.5"} 0
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="2"} 56687
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="2.5"} 56691
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="4"} 56692
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="8"} 56692
cometbft_consensus_proposal_timestamp_difference_bucket{chain_id="biyachain-1",is_timely="false",le="+Inf"} 56693
cometbft_consensus_proposal_timestamp_difference_sum{chain_id="biyachain-1",is_timely="false"} 1.8984521262875616e+06
cometbft_consensus_proposal_timestamp_difference_count{chain_id="biyachain-1",is_timely="false"} 56693
# HELP cometbft_consensus_quorum_prevote_delay Interval in seconds between the proposal timestamp and the timestamp of the earliest prevote that achieved a quorum.
# TYPE cometbft_consensus_quorum_prevote_delay gauge
cometbft_consensus_quorum_prevote_delay{chain_id="biyachain-1",proposer_address="3210F89D8743393559433C7F9D78ECA3996B8117"} 1.707502777
cometbft_consensus_quorum_prevote_delay{chain_id="biyachain-1",proposer_address="7764F2E6B8CF70E8BF8BE7C8B48595DCE483C6EF"} 1.7072774910000001
cometbft_consensus_quorum_prevote_delay{chain_id="biyachain-1",proposer_address="9D3D5D721115DF093D7D1B7B1C4ED3ED8DA886F4"} 1.7082260489999999
cometbft_consensus_quorum_prevote_delay{chain_id="biyachain-1",proposer_address="DF2AC1640CF3B4CFBEECE6191F1F32341CF72CFE"} 1.707415765
# HELP cometbft_consensus_round_duration_seconds Histogram of round duration.
# TYPE cometbft_consensus_round_duration_seconds histogram
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="0.1"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="0.2682695795279726"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="0.7196856730011522"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="1.9306977288832508"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="5.1794746792312125"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="13.894954943731381"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="37.27593720314942"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="100.00000000000006"} 56691
cometbft_consensus_round_duration_seconds_bucket{chain_id="biyachain-1",le="+Inf"} 56691
cometbft_consensus_round_duration_seconds_sum{chain_id="biyachain-1"} 55.09734169100039
cometbft_consensus_round_duration_seconds_count{chain_id="biyachain-1"} 56691
# HELP cometbft_consensus_round_voting_power_percent RoundVotingPowerPercent is the percentage of the total voting power received with a round. The value begins at 0 for each round and approaches 1.0 as additional voting power is observed. The metric is labeled by vote type.
# TYPE cometbft_consensus_round_voting_power_percent gauge
cometbft_consensus_round_voting_power_percent{chain_id="biyachain-1",vote_type="precommit"} 0.7505612372162633
cometbft_consensus_round_voting_power_percent{chain_id="biyachain-1",vote_type="prevote"} 0.9999999999999999
# HELP cometbft_consensus_rounds Number of rounds.
# TYPE cometbft_consensus_rounds gauge
cometbft_consensus_rounds{chain_id="biyachain-1"} 0
# HELP cometbft_consensus_step_duration_seconds Histogram of durations for each step in the consensus protocol.
# TYPE cometbft_consensus_step_duration_seconds histogram
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="0.1"} 56692
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="0.2682695795279726"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="0.7196856730011522"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="1.9306977288832508"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="5.1794746792312125"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="13.894954943731381"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="37.27593720314942"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="100.00000000000006"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Commit",le="+Inf"} 56693
cometbft_consensus_step_duration_seconds_sum{chain_id="biyachain-1",step="Commit"} 663.8981980479979
cometbft_consensus_step_duration_seconds_count{chain_id="biyachain-1",step="Commit"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="0.1"} 0
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="0.2682695795279726"} 0
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="0.7196856730011522"} 0
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="1.9306977288832508"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="5.1794746792312125"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="13.894954943731381"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="37.27593720314942"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="100.00000000000006"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewHeight",le="+Inf"} 56693
cometbft_consensus_step_duration_seconds_sum{chain_id="biyachain-1",step="NewHeight"} 84431.04230616974
cometbft_consensus_step_duration_seconds_count{chain_id="biyachain-1",step="NewHeight"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="0.1"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="0.2682695795279726"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="0.7196856730011522"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="1.9306977288832508"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="5.1794746792312125"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="13.894954943731381"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="37.27593720314942"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="100.00000000000006"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="NewRound",le="+Inf"} 56691
cometbft_consensus_step_duration_seconds_sum{chain_id="biyachain-1",step="NewRound"} 49.9714838389999
cometbft_consensus_step_duration_seconds_count{chain_id="biyachain-1",step="NewRound"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="0.1"} 40703
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="0.2682695795279726"} 56689
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="0.7196856730011522"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="1.9306977288832508"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="5.1794746792312125"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="13.894954943731381"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="37.27593720314942"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="100.00000000000006"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Precommit",le="+Inf"} 56693
cometbft_consensus_step_duration_seconds_sum{chain_id="biyachain-1",step="Precommit"} 5720.677044318992
cometbft_consensus_step_duration_seconds_count{chain_id="biyachain-1",step="Precommit"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="0.1"} 20262
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="0.2682695795279726"} 56373
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="0.7196856730011522"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="1.9306977288832508"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="5.1794746792312125"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="13.894954943731381"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="37.27593720314942"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="100.00000000000006"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Prevote",le="+Inf"} 56693
cometbft_consensus_step_duration_seconds_sum{chain_id="biyachain-1",step="Prevote"} 7099.24569648198
cometbft_consensus_step_duration_seconds_count{chain_id="biyachain-1",step="Prevote"} 56693
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="0.1"} 14563
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="0.2682695795279726"} 56674
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="0.7196856730011522"} 56688
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="1.9306977288832508"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="5.1794746792312125"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="13.894954943731381"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="37.27593720314942"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="100.00000000000006"} 56691
cometbft_consensus_step_duration_seconds_bucket{chain_id="biyachain-1",step="Propose",le="+Inf"} 56691
cometbft_consensus_step_duration_seconds_sum{chain_id="biyachain-1",step="Propose"} 4736.816298563027
cometbft_consensus_step_duration_seconds_count{chain_id="biyachain-1",step="Propose"} 56691
# HELP cometbft_consensus_total_txs Total number of transactions.
# TYPE cometbft_consensus_total_txs gauge
cometbft_consensus_total_txs{chain_id="biyachain-1"} 33
# HELP cometbft_consensus_validator_last_signed_height Last height signed by this validator if the node is a validator.
# TYPE cometbft_consensus_validator_last_signed_height gauge
cometbft_consensus_validator_last_signed_height{chain_id="biyachain-1",validator_address="7764F2E6B8CF70E8BF8BE7C8B48595DCE483C6EF"} 508434
# HELP cometbft_consensus_validator_power Power of a validator.
# TYPE cometbft_consensus_validator_power gauge
cometbft_consensus_validator_power{chain_id="biyachain-1",validator_address="7764F2E6B8CF70E8BF8BE7C8B48595DCE483C6EF"} 1000
# HELP cometbft_consensus_validators Number of validators.
# TYPE cometbft_consensus_validators gauge
cometbft_consensus_validators{chain_id="biyachain-1"} 4
# HELP cometbft_consensus_validators_power Total power of all validators.
# TYPE cometbft_consensus_validators_power gauge
cometbft_consensus_validators_power{chain_id="biyachain-1"} 4009
# HELP cometbft_mempool_active_outbound_connections Number of connections being actively used for gossiping transactions (experimental feature).
# TYPE cometbft_mempool_active_outbound_connections gauge
cometbft_mempool_active_outbound_connections{chain_id="biyachain-1"} 5
# HELP cometbft_mempool_already_received_txs Number of duplicate transaction reception.
# TYPE cometbft_mempool_already_received_txs counter
cometbft_mempool_already_received_txs{chain_id="biyachain-1"} 52
# HELP cometbft_mempool_disabled_routes Number of disabled routes.
# TYPE cometbft_mempool_disabled_routes gauge
cometbft_mempool_disabled_routes{chain_id="biyachain-1"} 5
# HELP cometbft_mempool_lane_bytes Number of used bytes per lane.
# TYPE cometbft_mempool_lane_bytes gauge
cometbft_mempool_lane_bytes{chain_id="biyachain-1",lane="default"} 0
# HELP cometbft_mempool_lane_size Number of uncommitted transactions per lane.
# TYPE cometbft_mempool_lane_size gauge
cometbft_mempool_lane_size{chain_id="biyachain-1",lane="default"} 0
# HELP cometbft_mempool_recheck_times Number of times transactions are rechecked in the mempool.
# TYPE cometbft_mempool_recheck_times counter
cometbft_mempool_recheck_times{chain_id="biyachain-1"} 2
# HELP cometbft_mempool_redundancy Redundancy level.
# TYPE cometbft_mempool_redundancy gauge
cometbft_mempool_redundancy{chain_id="biyachain-1"} 0
# HELP cometbft_mempool_size Number of uncommitted transactions in the mempool.  Deprecated: this value can be obtained as the sum of LaneSize.
# TYPE cometbft_mempool_size gauge
cometbft_mempool_size{chain_id="biyachain-1"} 0
# HELP cometbft_mempool_size_bytes Total size of the mempool in bytes.  Deprecated: this value can be obtained as the sum of LaneBytes.
# TYPE cometbft_mempool_size_bytes gauge
cometbft_mempool_size_bytes{chain_id="biyachain-1"} 0
# HELP cometbft_mempool_tx_life_span Duration in ms of a transaction in the mempool.
# TYPE cometbft_mempool_tx_life_span histogram
cometbft_mempool_tx_life_span_bucket{chain_id="biyachain-1",lane="default",le="50"} 33
cometbft_mempool_tx_life_span_bucket{chain_id="biyachain-1",lane="default",le="100"} 33
cometbft_mempool_tx_life_span_bucket{chain_id="biyachain-1",lane="default",le="200"} 33
cometbft_mempool_tx_life_span_bucket{chain_id="biyachain-1",lane="default",le="500"} 33
cometbft_mempool_tx_life_span_bucket{chain_id="biyachain-1",lane="default",le="1000"} 33
cometbft_mempool_tx_life_span_bucket{chain_id="biyachain-1",lane="default",le="+Inf"} 33
cometbft_mempool_tx_life_span_sum{chain_id="biyachain-1",lane="default"} -3.043712772162076e+20
cometbft_mempool_tx_life_span_count{chain_id="biyachain-1",lane="default"} 33
# HELP cometbft_mempool_tx_size_bytes Histogram of transaction sizes in bytes.
# TYPE cometbft_mempool_tx_size_bytes histogram
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="1"} 0
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="3"} 0
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="9"} 0
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="27"} 0
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="81"} 0
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="243"} 0
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="729"} 32
cometbft_mempool_tx_size_bytes_bucket{chain_id="biyachain-1",le="+Inf"} 33
cometbft_mempool_tx_size_bytes_sum{chain_id="biyachain-1"} 14034
cometbft_mempool_tx_size_bytes_count{chain_id="biyachain-1"} 33
# HELP cometbft_p2p_message_receive_bytes_total Number of bytes of each message type received.
# TYPE cometbft_p2p_message_receive_bytes_total counter
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_BlockPart"} 2.02611885e+08
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_HasProposalBlockPart"} 1.700754e+06
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_HasVote"} 2.148028e+07
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_NewRoundStep"} 1.4456466e+07
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_NewValidBlock"} 1.4456511e+07
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_Proposal"} 3.7134159e+07
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_StatusResponse"} 26
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_Vote"} 2.94991268e+08
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_VoteSetBits"} 1.25312e+06
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v1_VoteSetMaj23"} 1.914782e+06
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v2_HaveTx"} 396
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v2_ResetRoute"} 14
cometbft_p2p_message_receive_bytes_total{chain_id="biyachain-1",message_type="v2_Txs"} 35042
# HELP cometbft_p2p_message_send_bytes_total Number of bytes of each message type sent.
# TYPE cometbft_p2p_message_send_bytes_total counter
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_BlockPart"} 2.29791474e+08
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_HasProposalBlockPart"} 1.700754e+06
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_HasVote"} 2.1475862e+07
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_NewRoundStep"} 1.4456572e+07
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_NewValidBlock"} 1.4456409e+07
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_Proposal"} 4.209626e+07
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_StatusResponse"} 26
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_Vote"} 3.28227467e+08
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_VoteSetBits"} 1.223216e+06
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v1_VoteSetMaj23"} 1.925032e+06
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v2_HaveTx"} 396
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v2_ResetRoute"} 20
cometbft_p2p_message_send_bytes_total{chain_id="biyachain-1",message_type="v2_Txs"} 23469
# HELP cometbft_p2p_peer_pending_send_bytes Pending bytes to be sent to a given peer.
# TYPE cometbft_p2p_peer_pending_send_bytes gauge
cometbft_p2p_peer_pending_send_bytes{chain_id="biyachain-1",peer_id="5fc8be4de7f815dd1f2aa607ec9f5b26ef969a7e"} 0
cometbft_p2p_peer_pending_send_bytes{chain_id="biyachain-1",peer_id="741475b39fd741c82b2373d9a3742230a4fc70a2"} 0
cometbft_p2p_peer_pending_send_bytes{chain_id="biyachain-1",peer_id="90ab7585f3f4faed99bd05a5fb663ce82742ef51"} 0
cometbft_p2p_peer_pending_send_bytes{chain_id="biyachain-1",peer_id="a95b67eceaebf1445060ee7fe09523ca7e6cbafc"} 0
cometbft_p2p_peer_pending_send_bytes{chain_id="biyachain-1",peer_id="d392f2adcc3713ec863a76477be7701fb37d8a3a"} 0
# HELP cometbft_p2p_peers Number of peers.
# TYPE cometbft_p2p_peers gauge
cometbft_p2p_peers{chain_id="biyachain-1"} 5
# HELP cometbft_state_block_processing_time Time spent processing FinalizeBlock
# TYPE cometbft_state_block_processing_time histogram
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="1"} 0
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="11"} 56639
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="21"} 56691
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="31"} 56692
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="41"} 56692
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="51"} 56693
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="61"} 56693
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="71"} 56693
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="81"} 56693
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="91"} 56693
cometbft_state_block_processing_time_bucket{chain_id="biyachain-1",le="+Inf"} 56693
cometbft_state_block_processing_time_sum{chain_id="biyachain-1"} 350484.6960490013
cometbft_state_block_processing_time_count{chain_id="biyachain-1"} 56693
# HELP cometbft_state_consensus_param_updates Number of consensus parameter updates returned by the application since process start.
# TYPE cometbft_state_consensus_param_updates counter
cometbft_state_consensus_param_updates{chain_id="biyachain-1"} 56693
# HELP cometbft_state_store_access_duration_seconds The duration of accesses to the state store labeled by which method was called on the store.
# TYPE cometbft_state_store_access_duration_seconds histogram
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load",le="0.0002"} 3
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load",le="0.002"} 3
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load",le="0.02"} 3
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load",le="0.2"} 3
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load",le="2"} 3
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load",le="+Inf"} 3
cometbft_state_store_access_duration_seconds_sum{chain_id="biyachain-1",method="load"} 5.5511e-05
cometbft_state_store_access_duration_seconds_count{chain_id="biyachain-1",method="load"} 3
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_validators",le="0.0002"} 138594
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_validators",le="0.002"} 140710
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_validators",le="0.02"} 140710
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_validators",le="0.2"} 140710
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_validators",le="2"} 140710
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_validators",le="+Inf"} 140710
cometbft_state_store_access_duration_seconds_sum{chain_id="biyachain-1",method="load_validators"} 17.07272521400022
cometbft_state_store_access_duration_seconds_count{chain_id="biyachain-1",method="load_validators"} 140710
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save",le="0.0002"} 6619
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save",le="0.002"} 56652
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save",le="0.02"} 56694
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save",le="0.2"} 56694
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save",le="2"} 56694
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save",le="+Inf"} 56694
cometbft_state_store_access_duration_seconds_sum{chain_id="biyachain-1",method="save"} 15.039665654999974
cometbft_state_store_access_duration_seconds_count{chain_id="biyachain-1",method="save"} 56694
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="saveValidatorsInfo",le="0.0002"} 56695
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="saveValidatorsInfo",le="0.002"} 56695
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="saveValidatorsInfo",le="0.02"} 56695
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="saveValidatorsInfo",le="0.2"} 56695
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="saveValidatorsInfo",le="2"} 56695
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="saveValidatorsInfo",le="+Inf"} 56695
cometbft_state_store_access_duration_seconds_sum{chain_id="biyachain-1",method="saveValidatorsInfo"} 0.15899321600000085
cometbft_state_store_access_duration_seconds_count{chain_id="biyachain-1",method="saveValidatorsInfo"} 56695
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_abci_responses",le="0.0002"} 0
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_abci_responses",le="0.002"} 56645
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_abci_responses",le="0.02"} 56691
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_abci_responses",le="0.2"} 56693
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_abci_responses",le="2"} 56693
cometbft_state_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_abci_responses",le="+Inf"} 56693
cometbft_state_store_access_duration_seconds_sum{chain_id="biyachain-1",method="save_abci_responses"} 19.323406151000007
cometbft_state_store_access_duration_seconds_count{chain_id="biyachain-1",method="save_abci_responses"} 56693
# HELP cometbft_store_block_store_access_duration_seconds The duration of accesses to the state store labeled by which method was called on the store.
# TYPE cometbft_store_block_store_access_duration_seconds histogram
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_block_meta",le="0.0002"} 79425
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_block_meta",le="0.002"} 83064
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_block_meta",le="0.02"} 83064
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_block_meta",le="0.2"} 83064
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_block_meta",le="2"} 83064
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_block_meta",le="+Inf"} 83064
cometbft_store_block_store_access_duration_seconds_sum{chain_id="biyachain-1",method="load_block_meta"} 3.6024640400000263
cometbft_store_block_store_access_duration_seconds_count{chain_id="biyachain-1",method="load_block_meta"} 83064
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_seen_commit",le="0.0002"} 7271
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_seen_commit",le="0.002"} 7272
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_seen_commit",le="0.02"} 7272
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_seen_commit",le="0.2"} 7272
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_seen_commit",le="2"} 7272
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="load_seen_commit",le="+Inf"} 7272
cometbft_store_block_store_access_duration_seconds_sum{chain_id="biyachain-1",method="load_seen_commit"} 0.41553120699999946
cometbft_store_block_store_access_duration_seconds_count{chain_id="biyachain-1",method="load_seen_commit"} 7272
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="new_block_store",le="0.0002"} 0
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="new_block_store",le="0.002"} 1
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="new_block_store",le="0.02"} 1
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="new_block_store",le="0.2"} 1
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="new_block_store",le="2"} 1
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="new_block_store",le="+Inf"} 1
cometbft_store_block_store_access_duration_seconds_sum{chain_id="biyachain-1",method="new_block_store"} 0.001546453
cometbft_store_block_store_access_duration_seconds_count{chain_id="biyachain-1",method="new_block_store"} 1
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block",le="0.0002"} 0
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block",le="0.002"} 56648
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block",le="0.02"} 56692
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block",le="0.2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block",le="2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block",le="+Inf"} 56693
cometbft_store_block_store_access_duration_seconds_sum{chain_id="biyachain-1",method="save_block"} 27.842785415000215
cometbft_store_block_store_access_duration_seconds_count{chain_id="biyachain-1",method="save_block"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_part",le="0.0002"} 56691
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_part",le="0.002"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_part",le="0.02"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_part",le="0.2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_part",le="2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_part",le="+Inf"} 56693
cometbft_store_block_store_access_duration_seconds_sum{chain_id="biyachain-1",method="save_block_part"} 0.45694602600000334
cometbft_store_block_store_access_duration_seconds_count{chain_id="biyachain-1",method="save_block_part"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_to_batch",le="0.0002"} 56686
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_to_batch",le="0.002"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_to_batch",le="0.02"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_to_batch",le="0.2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_to_batch",le="2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_block_to_batch",le="+Inf"} 56693
cometbft_store_block_store_access_duration_seconds_sum{chain_id="biyachain-1",method="save_block_to_batch"} 1.4445943780000068
cometbft_store_block_store_access_duration_seconds_count{chain_id="biyachain-1",method="save_block_to_batch"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_bs_state",le="0.0002"} 0
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_bs_state",le="0.002"} 56650
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_bs_state",le="0.02"} 56692
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_bs_state",le="0.2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_bs_state",le="2"} 56693
cometbft_store_block_store_access_duration_seconds_bucket{chain_id="biyachain-1",method="save_bs_state",le="+Inf"} 56693
cometbft_store_block_store_access_duration_seconds_sum{chain_id="biyachain-1",method="save_bs_state"} 24.156808485999964
cometbft_store_block_store_access_duration_seconds_count{chain_id="biyachain-1",method="save_bs_state"} 56693
# HELP end_blocker end_blocker
# TYPE end_blocker summary
end_blocker{module="crisis",quantile="0.5"} 0.0013599999947473407
end_blocker{module="crisis",quantile="0.9"} 0.0015829999465495348
end_blocker{module="crisis",quantile="0.99"} 0.0015829999465495348
end_blocker_sum{module="crisis"} 86.99632299726363
end_blocker_count{module="crisis"} 56693
end_blocker{module="gov",quantile="0.5"} 0.020021000877022743
end_blocker{module="gov",quantile="0.9"} 0.024594999849796295
end_blocker{module="gov",quantile="0.99"} 0.024594999849796295
end_blocker_sum{module="gov"} 1449.4971818784252
end_blocker_count{module="gov"} 56693
end_blocker{module="staking",quantile="0.5"} 0.13750100135803223
end_blocker{module="staking",quantile="0.9"} 0.16611400246620178
end_blocker{module="staking",quantile="0.99"} 0.16611400246620178
end_blocker_sum{module="staking"} 9388.374241381884
end_blocker_count{module="staking"} 56693
# HELP go_gc_duration_seconds A summary of the wall-time pause (stop-the-world) duration in garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.000124427
go_gc_duration_seconds{quantile="0.25"} 0.000170416
go_gc_duration_seconds{quantile="0.5"} 0.000187259
go_gc_duration_seconds{quantile="0.75"} 0.000221345
go_gc_duration_seconds{quantile="1"} 0.000487341
go_gc_duration_seconds_sum 0.176122674
go_gc_duration_seconds_count 888
# HELP go_gc_gogc_percent Heap size target percentage configured by the user, otherwise 100. This value is set by the GOGC environment variable, and the runtime/debug.SetGCPercent function. Sourced from /gc/gogc:percent.
# TYPE go_gc_gogc_percent gauge
go_gc_gogc_percent 100
# HELP go_gc_gomemlimit_bytes Go runtime memory limit configured by the user, otherwise math.MaxInt64. This value is set by the GOMEMLIMIT environment variable, and the runtime/debug.SetMemoryLimit function. Sourced from /gc/gomemlimit:bytes.
# TYPE go_gc_gomemlimit_bytes gauge
go_gc_gomemlimit_bytes 9.223372036854776e+18
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 189
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.23.9"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated in heap and currently in use. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.737791864e+09
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated in heap until now, even if released already. Equals to /gc/heap/allocs:bytes.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.78110155744e+11
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table. Equals to /memory/classes/profiling/buckets:bytes.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 6.660591e+06
# HELP go_memstats_frees_total Total number of heap objects frees. Equals to /gc/heap/frees:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 2.864121631e+09
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata. Equals to /memory/classes/metadata/other:bytes.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.1730896e+07
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and currently in use, same as go_memstats_alloc_bytes. Equals to /memory/classes/heap/objects:bytes.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.737791864e+09
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used. Equals to /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 7.1221248e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.835057152e+09
# HELP go_memstats_heap_objects Number of currently allocated objects. Equals to /gc/heap/objects:objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 3.0297243e+07
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS. Equals to /memory/classes/heap/released:bytes.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 155648
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system. Equals to /memory/classes/heap/objects:bytes + /memory/classes/heap/unused:bytes + /memory/classes/heap/released:bytes + /memory/classes/heap/free:bytes.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 1.9062784e+09
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.77088428123333e+09
# HELP go_memstats_mallocs_total Total number of heap objects allocated, both live and gc-ed. Semantically a counter version for go_memstats_heap_objects gauge. Equals to /gc/heap/allocs:objects + /gc/heap/tiny/allocs:objects.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 2.894418874e+09
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures. Equals to /memory/classes/metadata/mcache/inuse:bytes.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 9600
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system. Equals to /memory/classes/metadata/mcache/inuse:bytes + /memory/classes/metadata/mcache/free:bytes.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures. Equals to /memory/classes/metadata/mspan/inuse:bytes.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 3.06832e+07
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system. Equals to /memory/classes/metadata/mspan/inuse:bytes + /memory/classes/metadata/mspan/free:bytes.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 3.19056e+07
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place. Equals to /gc/heap/goal:bytes.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 3.31510884e+09
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations. Equals to /memory/classes/other:bytes.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 4.356953e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes obtained from system for stack allocator in non-CGO environments. Equals to /memory/classes/heap/stacks:bytes.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 2.031616e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator. Equals to /memory/classes/heap/stacks:bytes + /memory/classes/os-stacks:bytes.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 2.031616e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system. Equals to /memory/classes/total:byte.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.972979656e+09
# HELP go_sched_gomaxprocs_threads The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. Sourced from /sched/gomaxprocs:threads.
# TYPE go_sched_gomaxprocs_threads gauge
go_sched_gomaxprocs_threads 8
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 15
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 12110.25
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 65536
# HELP process_network_receive_bytes_total Number of bytes received by the process over the network.
# TYPE process_network_receive_bytes_total counter
process_network_receive_bytes_total 2.12452548998e+11
# HELP process_network_transmit_bytes_total Number of bytes sent by the process over the network.
# TYPE process_network_transmit_bytes_total counter
process_network_transmit_bytes_total 1.90525048125e+11
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 67
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.169864192e+09
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.77078162034e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 4.931526656e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 6847
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP runtime_alloc_bytes runtime_alloc_bytes
# TYPE runtime_alloc_bytes gauge
runtime_alloc_bytes 1.737690624e+09
# HELP runtime_free_count runtime_free_count
# TYPE runtime_free_count gauge
runtime_free_count 2.8641216e+09
# HELP runtime_gc_pause_ns runtime_gc_pause_ns
# TYPE runtime_gc_pause_ns summary
runtime_gc_pause_ns{quantile="0.5"} NaN
runtime_gc_pause_ns{quantile="0.9"} NaN
runtime_gc_pause_ns{quantile="0.99"} NaN
runtime_gc_pause_ns_sum 141556
runtime_gc_pause_ns_count 1
# HELP runtime_heap_objects runtime_heap_objects
# TYPE runtime_heap_objects gauge
runtime_heap_objects 3.0296734e+07
# HELP runtime_malloc_count runtime_malloc_count
# TYPE runtime_malloc_count gauge
runtime_malloc_count 2.894418432e+09
# HELP runtime_num_goroutines runtime_num_goroutines
# TYPE runtime_num_goroutines gauge
runtime_num_goroutines 184
# HELP runtime_sys_bytes runtime_sys_bytes
# TYPE runtime_sys_bytes gauge
runtime_sys_bytes 1.972979712e+09
# HELP runtime_total_gc_pause_ns runtime_total_gc_pause_ns
# TYPE runtime_total_gc_pause_ns gauge
runtime_total_gc_pause_ns 1.76122672e+08
# HELP runtime_total_gc_runs runtime_total_gc_runs
# TYPE runtime_total_gc_runs gauge
runtime_total_gc_runs 888
# HELP wasmvm_cache_elements_total Total number of elements in the cache
# TYPE wasmvm_cache_elements_total gauge
wasmvm_cache_elements_total{type="memory"} 0
wasmvm_cache_elements_total{type="pinned"} 0
# HELP wasmvm_cache_hits_total Total number of cache hits
# TYPE wasmvm_cache_hits_total counter
wasmvm_cache_hits_total{type="fs"} 0
wasmvm_cache_hits_total{type="memory"} 0
wasmvm_cache_hits_total{type="pinned"} 0
# HELP wasmvm_cache_misses_total Total number of cache misses
# TYPE wasmvm_cache_misses_total counter
wasmvm_cache_misses_total 0
# HELP wasmvm_cache_size_bytes Total number of elements in the cache
# TYPE wasmvm_cache_size_bytes gauge
wasmvm_cache_size_bytes{type="memory"} 0
wasmvm_cache_size_bytes{type="pinned"} 0
