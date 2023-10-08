# States & Pipeline

## states

### Current Design

TX_STATE_PENDING_WAIT                TxState = iota - 10 // -10; pending because of notEnoughCash or futureNonce
_RESERVE_1                                               // -9
_RESERVE_2                                               // -8
TX_STATE_PENDING_TX_EMPTY_RETRY                          // -7
TX_STATE_PENDING_RETRY_UPPER_GAS                         // -6
TX_STATE_PENDING_RETRY                                   // -5
TX_STATE_SEND_FAILED_RETRY_UPPER_GAS                     // -4
TX_STATE_SEND_FAILED_RETRY                               // -3
TX_STATE_EXECUTE_FAILED                                  // -2
TX_STATE_SEND_FAILED                                     // -1
TX_STATE_INIT                                            // 0
TX_STATE_POPULATED                                       // 1
TX_STATE_READY_OR_INSERTED                               // 2
TX_STATE_EXECUTED                                        // 3
TX_STATE_CONFIRMED                                       // 4

### 需要考虑的情况

理想状态下交易应该经历以下几个状态：

`QUEUE`：交易在数据库的待处理队列中，尚未发送

`SENT` 交易已发送，但尚未被打包进区块
`LATEST` 交易已被打包进区块
`SAFE` 交易已被打包进区块 且已处于safe（～5 minutes or more）
`FINALIZED` 交易已被打包进区块 且已处于finalized（～15 minutes）

到达 `FINALIZED` 状态后，交易的状态达到稳定并不会再变更。

但实际上可能交易会因为种种原因等待/失败，因此状态会更为复杂。我希望交易的状态设计能考虑到以下问题：

1. 失败的原因
2. 交易的nonce是否分配
3. 是否有关联交易
4. 是否reorg
5. 是否已上链

### new design

tx_status and task_status are isolated
task_status is the aggregation of all relating tx_status

为什么？
task会和多条交易关联（重试等），而多条交易的状态变化是难以预测的(或者说，想要完全的确定性会导致性能可用性极大幅损失)，甚至可能不是有限的，想要精确描述task的状态是非常困难的。

因此不使用status来描述task处于的具体状态，或者说，接口展示的task status并不是交易的真正状态，而是通过接口来判断task状态，从而进行处理。

task 包含一个cancel字段，代表是否用户希望取消交易。该字段相当于更改了交易的元数据 （data, from...）。

定义不同的决策，实现以下接口，决策器决定了交易发送策略。保守的决策与激进的决策对于特定接口的返回策略不一致。

strategy接收task、关联交易列表

- in_error() 有关联交易，且没有任何一条关联交易成功
- in_queue() 没有关联交易或者所有关联交易都处于in_queue
- latest_error() 读取最新一条关联交易的error
- is_stable() 所有关联交易都处于stable状态
- is_stable_enough() 所有关联交易都处于某级别的stable状态（不同决策行为不同）
- succeed() 有且仅有一条关联交易成功
- in_unexpeced_result() 依赖的task失败但本task成功；task关联交易中有复数条成功;...
- needs_sponsor() 任意一条关联交易处于 insufficient_balance
- needs_retry() needs_sponsor() 为 false 且最后一条关联交易为error且没有超过最大重试次数
- ready_to_send() 依赖的task已成功/没有失败(不同决策器行为可能不同)

- is_cancelling()
- cancelled()

user interface task status

PROCESSING 处理中
LIKELY_SUCCESS 大概率成功
SUCCESS 成功
FAILURE 失败（retry 次数用尽）

CANCELLING 取消中
CANCELLED 取消成功
CANCEL_FAILURE 取消失败

#### tx status

关联交易类型

- type
  - task
  - replacement
- chain_status
  - queue(需要引入queue吗？)
  - pending
  - latest
  - safe
  - finalized
- error
  - discarded
  - insufficient_balance
  - pending_too_long
  - ...

interfaces

- is_stable()
  - 交易处于 finalized 状态
  - 交易处于 pending 但 nonce 相同的交易已经为finalized状态（交易已经确定fail）

- is_safe_stable()
  - 交易处于 safe/finalized 状态
  - 交易处于 pending 但nonce相同的交易已经为safe/finalized状态

- is_somewhat_stable()
  - 交易处于 safe/finalized 状态 （交易已上链一分钟。。。）
  - 交易处于 pending 但nonce相同的交易已经为safe/finalized状态

...

## pipeline

service 1:

tx chain status updater

```py

需要大范围的锁？

while True:
    db.lock()
    # 有很多优化的空间
    for account in accounts:
        stable_nonce, somewhat_nonce, latest_nonce = chain_fetcher.get_nonces(account)
        account.write(account, stable_nonce, somewhat_nonce, latest_nonce)

    for tx in db.tx_table if is_not_stable(tx):
        chain_fetcher.update_tx(tx)
    db.unlock()
    time.sleep(xx)
```

service 2:

chain_writer

```py

batch = Batch()

while True:
    # 余额不足不是need retry
    for task in db.task_table if task.strategy.needs_retry(task):
        if batch.is_full():
            break
        batch.add(retryer.get_next_retry_tx(task))
    for task in order_by_priority(db.task_table) if task.strategy.in_queue(task):
        if batch.is_full():
            break
        batch.add(transaction_builder.build(task))
    txs = batch.send()
    write_tx(txs, database.tx_table)
    time.sleep(xx)
```

service 3:

sponsor/paymaster

```py
while True:
    address_to_sponsor = set()
    for tx in db.tx_table if sent(tx):
        if tx.error.is_insufficient_balance():
            set.add(tx.from_address)
    ... # 可能有其他的sponsor方法，比如计算in_queue的task开销，并和余额比较
    sponsor(address_to_sponsor)
    time.sleep(xx)
```
