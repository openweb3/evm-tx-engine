# Interfaces

你是一位精通web3与web2的系统架构师。
我正在进行交易上链引擎的接口设计工作，该引擎的技术目标是屏蔽web3的底层细节，自动完成账户管理、交易组装、上链、监控、重试等操作，同时支持多链。现在请你协助我完成接口设计工作，指出我设计中存在的潜在问题，并给我相应建议。
目前我正在设计的是 1）交易上链任务创建接口 与 2）交易状态接口

## tasks

### signle task

```bash
[post] /task

{
    "chain": "ethereum",
    "from" : "0x1234567890123456789012345678901234567890",
    "fields": 
    {
        "to": "0x12345678901234567890123456789",
        "function": "transfer(address,uint256)",
        "params": ["0x12345678901234567890123456789",1000000000000000000],
        "maxFeePerGas": 1000000000000000000,
    }
    "fee": "normal", // gas price strategy
    "requestId": "2345678945678kjhgf", // Optional, the identifier for the request, should be unique for each new request. It will prevent commit same request twice or more, the request id will expire after a certain period, for example, 1 day
    "retry": 
    {
        "maxAttempts": 3 // Optional, the max attempts that one transaction would retry
        "deadline": 1556842100000 // Optional, the deadline that the request will not be retried. Could be timestamp, blockNumber, ...
    }
}

returns

"345678987654345678" // taskId
```

error:

- duplicate request: the requestId already exists in the database
- invalid request during the first step check: ...
  - invalid chain: ...
  - invalid from: ...
  - invalid function/parameter

```bash
[get] /task/:taskId

returns

{
    "id": "345678987654345678",
    "chain": "ethereum",
    "from" : "0x1234567890123456789012345678901234567890",
    "status": "pending",
    "fields": // content fields are chain-specific fields
    {
        "to": "0x12345678901234567890123456789",
        "value": 0,
        "function": "transfer(address,uint256)", // function, selector, params and data are duplicate fields
        "selector": "0x41c3f5",
        "params": ["0x12345678901234567890123456789",1000000000000000000],
        "data": "0x......"
    },
    "fee": "normal",
    "hash": null,
    "previousTask": null,
    "history": [
        {
            "attempt": 1,
            "timestamp": 1556842100000,
            "errorCode": 2001,
            "errorMessage": "Transaction failed - chain error",
            "transactionHash": "0x123..."
        },
        {
            "attempt": 2,
            "timestamp": 1556843100000,
            "errorCode": null,
            "errorMessage": null,
            "transactionHash": "0x456..."
        }
    ],
    "retry": 
    {
        maxAttempts: 3,
        deadline: 1556842100000

    }
}
```

### task batch

``` bash
[post] /taskBatch

{
    "type": "unordered", // could be unordered or ordered. If ordered, all task should have same `from` and `chain` field. All tasks in the request will be executed in order, and subsequent task will not be on chain unless the previous one succeeds
    "tasks": [ // same parameters except for request id field
        {
            "chain": "ethereum",
            "from" : "0x1234567890123456789012345678901234567890",
            "fields": 
            {
                "to": "0x12345678901234567890123456789",
                "function": "transfer(address,uint256)",
                "params": ["0x12345678901234567890123456789",1000000000000000000],
            }
        },
        {
            "chain": "ethereum",
            "from" : "0x1234567890123456789012345678901234567890",
            {
                "to": "0x12345678901234567890123456789",
                "function": "transfer(address,uint256)",
                "params": ["0x12345678901234567890123456789",1000000000000000000],
            }
        }
    ],
    "requestId": "2345678945678kjhgf", // Optional, the identifier for the request, should be unique for each new request. It helps clients to prevent commiting same request twice or more. The request id will expire after a certain period, for example, 1 day
}

returns

{
    "batchId": "34567890987654345678", // the id of the batch
    "taskIds": [ // the task ids of the batch
        "4567887654456",
        "456e3456",
        ...
    ]
}
```

error:

- duplicate request: the requestId already exists in the database
- invalid request during the first step check: ...
  - invalid chain: ...
  - invalid from: ...
  - invalid function/parameter

Note that the entire request will fail if any of the requests fail.

``` bash
[get] /taskBatch/:batchId

{
    "type": "unordered", // could be unordered or ordered. If ordered, all task should have same `from` and `chain` field. All tasks in the request will be executed in order, and subsequent task will not be on chain unless the previous one succeeds
    "batchId": "dfghjkjhgfghjk",
    "batchStatus": "pending",
    "tasks": [ // same parameters except for request id field
        {
            "id": "345678987654345678",
            "chain": "ethereum",
            "from" : "0x1234567890123456789012345678901234567890",
            "status": "pending",
            "content": // content fields are chain-specific fields
            {
                "to": "0x12345678901234567890123456789",
                "value": 0,
                "function": "transfer(address,uint256)", // function, selector, params and data are duplicate fields
                "selector": "0x41c3f5",
                "params": ["0x12345678901234567890123456789",1000000000000000000],
                "data": "0x......"
            },
            "hash": null,
            "previousTask": null,
            "history": [
                {
                    "attempt": 1,
                    "timestamp": 1556842100000,
                    "errorCode": 2001,
                    "errorMessage": "Transaction failed - chain error",
                    "transactionHash": "0x123..."
                },
                {
                    "attempt": 2,
                    "timestamp": 1556843100000,
                    "errorCode": null,
                    "errorMessage": null,
                    "transactionHash": "0x456..."
                }
            ],
            "retry": 
            {
                maxAttempts: 3,
                deadline: 1556842100000

            }
        },
        ...
    ]
}
```

### modify specific task

add retry

## cancel task

cancel specific task in queue by taskId. It should be noticed that the stop might not succeed

```bash
[post] /task/:taskId/cancel

returns cancel status (potential async cancel)

{
    "status": "processing",
    "message": "Cancellation request is being processed.",
}

{
    "status": "success",
    "message": "Task successfully cancelled.",
}
```

```bash
[get] /task/:taskId/cancel

returns

{
    "status": "failure",
    "errorCode": "3002",
    "errorMessage": "Task already executed on chain.",
    "taskId": "345678987654345678"
}
```

## Account Manager Adapter

the account is an offchain module provides universal apis for the adapted account implementation.

The actual account implementation could be a simple database-based implementation or AWS/aliyun/... KMS

It should be noted that the adaptor might only provide limited methods to interact with the actual implementation.

```bash
get /accounts/wrappingPublicKey

returns

{
    "wrappingPublicKey": "fnewoafnawoge",
    "wrappingAlg": "fjweognawiowa",
}
```

```bash
post /accounts

{
    "chain": "ethereum", // polygon, all, ...
    "wrappedPrivateKey": "56789876bcdef...", // optional key import, should be encrypted by wrapping key
    "wrappingPublicKey": "fnewoafnawoge" // optional the public key used to encrypt the private key,
    "wrappingAlg": "fjweognawiowa",
}

returns

"0x34567898765abcdef..."
```

```bash
get /accounts?chain=ethereum

returns

[
    {
        "address": "0x34567898765abcdef...",
        "publicKey": "fjeoiwanfoawfw...",
        "chains": ["Ethereum", "Polygon", "Conflux eSpace"], // which chains could the account be used
        "alias": "my-account",
        "status": "active",
        "alg": "ecc-secp256k1"
    },
    {
        "address": "0x34567898765abcdee...",
        "publicKey": "fjeoiwanfoawfwfwenfenwio...",
        "chains": ["Ethereum"],
        "alias": "my-account",
        "alg": "ecc-secp256k1"
        "status": "active"
    }
]
```

### internal api

```bash
post /accounts/:address/sign_raw
```

```bash
post /accounts/:address/sign_transaction
```

```bash
post /accounts/:address/encrypt
```

```bash
post /accounts/:address/verify
```

## enums

### HTTP Error Code

- duplicate request: the requestId already exists in the database
- invalid request during the first step check: ...
  - invalid chain: ...
  - invalid from: ...
  - invalid function/parameter
- internal error

### Error Code

错误码分为五位。

第一位为错误的大类

- 00000 无错误
- 1xxxx. 请求参数错误
- 2xxxx. tx engine 错误（如内部逻辑约束导致的错误）
- 3xxxx. web3 相关错误。如sdk相关错误（json-rpc错误）/链上错误（交易执行错误）
- 4xxxx. 数据库操作错误/异常
- 5xxxx. Account Adapter 相关错误

#### 基础错误码

#### 1xxxx. 请求参数错误

- 10001. 缺少必要参数
- 10002. 参数格式无效
- 10003. 参数值不在允许的范围内
- 10004. JSON体解析错误
- 10005. 请求方法不允许（例如GET/POST/PUT错误）
- 10006. 重复的requestId

#### 2xxxx. Tx Engine 错误

- 20001. 交易重试达到上限
- 20002. 余额不足
- 20003. 
- 20004. 
- 20101. 交易取消失败

#### 3xxxx. Web3 相关错误

- 30001. JSON-RPC通信错误
- 30002. 链上状态查询错误
- 30003. SDK初始化错误
- 30004. 无效的链连接信息（例如无效的节点URL）
- 30005. estimate失败
- 30101. 交易发送失败
- 30201. 链上交易执行失败

#### 4xxxx. 数据库操作错误/异常

- 40001. 数据库连接失败
- 40002. 数据查询失败
- 40003. 数据插入失败
- 40004. 数据更新失败
- 40005. 数据删除失败

#### 5xxxx. Account Adapter 相关错误

- 50001. 账户创建失败
- 50002. 账户导入失败
- 50101. 交易签名错误（通用）
- 50102. 错误的交易签名参数
- 50103. 错误的交易签名chain参数
- 50104. 指定的发送者不存在

这样的设计使得开发人员和用户能够更快地理解和定位问题，同时为未来可能出现的新错误留有足够的空间。同时，每个错误码都应该有对应的详细错误信息，以帮助解决问题。

### webhook

中划线
版本号
回调函数（先不管，但是留出实现空间）
优先级更改
错误设计
pipeline & state machine

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
- is_likely_stable() 所有关联交易都处于likely_stable状态（不同决策行为不同）
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

- is_likely_stable()
  - 交易处于 safe/finalized 状态
  - 交易处于 pending 但nonce相同的交易已经为safe/finalized状态
- is_stable()
  - 交易处于 finalized 状态
  - 交易处于 pending 但 nonce 相同的交易已经为finalized状态（交易已经确定fail）

#### QUEUE

交易在数据库的待处理队列中，尚未发送

- INIT
  - WAITING_FOR_PREVIOUS_TRANSACTION
- POPULATED

#### SENT

交易已发送，在矿工交易池内，但尚未被打包进区块

- WAITING_FOR_MINER

#### ONCHAIN

- LATEST
- SAFE
- FINALIZED

#### CANCELLED

对于取消详情，通过 tasks/:id/cancel 获取
此时如果hash不为空

- QUEUE：在队列中时被取消
- REPLACED：发送到链上后被取消

#### ERROR

- QUEUE
  - ESTIMATION_ERROR 在estimate时出现了错误
  - PREVIOUS_TX_FAILED 上一笔交易失败
  - INSUFFICIENT_BALANCE 余额不足因此未发送
  - EXCEEDING_MAX_RETRY 超过最大尝试次数

- SENT
  - PENDING_TIMEOUT 在发送后未上链的等待时间超时
    - FUTURE_NONCE
    - LOW_PRICE
  - INSUFFICIENT_BALANCE 余额不足

- ONCHAIN
  - EXECUTION_ERROR 发生了执行错误

#### UNCERTAIN

- SENT 交易已发送
  - RELYING_TASK_FAILED ordered batch 中上一笔任务失败
  - REPLACING 正在发生交易替换，复数笔相同nonce的交易在交易池中
  - OFFCHAIN_AFTER_REORG // 交易在 REORG 前已上链，在 REORG 后进入了 SENT 状态

#### UNEXPECTED

- ONCHAIN
  - RELYING_TASK_FAILED 上一条交易
  - DISCARDED_HISTORY_SUCCESS 被丢弃的 history 交易成功了
  - WRONG_ORDER 上链交易的顺序出错（ordered batch）

## legacy_draft

- PENDING_TIMEOUT：交易在发送后未上链的时间超时

- INSUFFICIENT_BALANCE: 余额不足
  - QUEUE
  - SENT
  - ONCHAIN（？）
- EXECUTION_ERROR：发生了执行错误
  - ONCHAIN（已上链时出现错误）
- RELYING_TX_FAILED：上一个交易失败
  - QUEUE: 本条交易尚未发送
  - ~~SENT：本条交易已发送，但尚未上链（这里应该不是error）~~
  - ~~ONCHAIN_FAILED：本条交易已上链，但执行失败~~ （这里不需要单独的状态，如果fail了就加到EXECUTION_ERROR中去，此时的行为是符合预期的，如果不合预期就加入到UNEXPECTED中去）
- EXCEEDS_MAX_RETRY

UNEXPECTED

- PREVIOUS_TX_FAILED 上一条交易
  - ONCHAIN_SUCCESS
- DISCARDED_HISTORY_SUCCESS 被丢弃的 history 交易成功了
- WRONG_ORDER 上链交易的顺序出错（ordered batch）

<https://github.com/pytransitions/transitions>

unexpected 说明：

1. 依赖的交易失败了，但本条交易却成功了
2. 相同task对应多笔交易，且nonce不同，出于各种原因，复数交易均成功上链

实现逻辑会尽可能避免上述情况发生。但由于存在reorg等情况，除非放弃交易效率，使用finalized标签进行交易确认，否则总存在发生的概率。

进行交易替换

1. 存在可能的reorg，可能会有预期以外的行为，这些行为可能不是单纯的成功或失败
2. ~~替换交易时，旧交易无法保证被丢弃~~（发生替换交易时应该怎么设计？）
3. 存在reorg，且存在交易替换
4. 存在reorg

特殊情况
ordered batch
A-B-C

A 失败，但是 B 成功

这种情况需要尽可能避免，但是仍有极小的可能出现。如

A 成功后 B 成功

发生 reorg 后 A 失败，但是 B 成功

sent -> cancelling -> onchain
                  |-> replaced

1. 发送到链上 -> 发送替代交易 -> 替代交易成功
2. 发送到链上 -> 发送替代交易 -> 替代交易失败
3. 发送到链上 -> 发送替代交易

## design

task: status, state

1. database: 存储所有数据与缓存的链上状态
   1. task_table
   2. tx_table
2. nonce_manager: 管理nonce
3. gas_manager: 管理gas
4. retryer: 管理交易重试
5. chain_fetcher: 轮询链上状态
6. transaction_builder: 构造交易


service 1:

tx chain status updater

```py

while True:
    for tx in db.tx_table if is_not_stable(tx):
        task = tx.relating_task
        relying_task, same_nonce_txs = relating(txs)
        # updated state will have a related timestamp
        # apart from querying on chain tx status
        # 为什么需要relying_task: 为了保证交易顺序，需要保证交易顺序与链上一致
        # 为什么需要same_nonce_txs: 确保交易被丢弃的状态且不会再上链

        # 有很多优化的空间
        chain_fetcher.update_tx_and_relating_task(tx, task, relying_task, same_nonce_txs)
    time.sleep(xx)

```

service 2:

chain_writer

```py

batch = Batch()

while True:
    # 余额不足不是need retry
    for task in db.task_table if needs_retry(task):
        if batch.is_full():
            break
        batch.add(retryer.get_next_retry_tx(task))
    for task in order_by_priority(db.task_table) if in_queue(task):
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
