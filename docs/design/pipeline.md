# pipeline

## tx life cycle

tasklist
|

| certain task will init a bare transaction and put into targetQueue

|
targetQueue: no bigger than several hundreds(e.g. 200 ?)
|
| execution simulator
| gasLimit 在这里分配
|
| transactions will go to gasEnoughQueue if gas is enough
| the database（redis?） maintains a field named `internalPendingBalance`
| 全局变量 MAX_GAS_PRICE 会被用来估算gas开销
| 使用 eth_getBalance(..., 'finalized') 来估算余额
|
| 交易的费用用粗粒度的方式估算？（不打算在这一步完全避免gas不足的问题？）
|
gasEnoughQueue
|
| nonce manager(是否要移到 targetQueue -> gasEnoughQueue ?)
| 理论上来说，nonce manager知道所有用户下nonce的分布状况
| 一个简单的实现——数据库会维护单个字段`internalNextNonce`。代表下次分配的nonce值。该值考虑了constructedQueue中的nonce分配情况
|
| & gas manager：根据交易/task配置为交易设置 price
|
constructed：所有交易参数都已经确定, 但是没有被签名发出
|
| 签名服务签名交易
|
signed: 已签名的 raw_transaction
|
| send to nodes
|
| 发送的结果分为两种: 1. 交易被节点接收。 2. 报错，交易未被接收（余额不足/超出rate limit/。。。）
|
pending -> latest -> safe -> finalized
｜

｜status watcher 观察pending交易状态, 观察链上交易状态
｜注意：交易一旦被节点接收，就需要一直维护该交易的状态，因为我们无法保证该交易会被丢弃

｜

｜
error
| - 简单的需要重发，或在条件满足后重发（如余额不足）
| - 需要构造新交易并签名

**如果 tx 进入 error 状态，不进行任何处理，监听task的服务会提交新的tx**

## alg

**数据一致性？**

### 

### signer service

将constructed状态的交易提交给signer

监听 constructed：

for tx in constructed:
    raw = signer.sign(tx)
    tx.raw = raw
    tx.status = "signed"

### sender service

从signed状态的交易中不断取出batch，并批量发送至链上

### chain watcher

监听链上不处于 stable 状态的交易状态

### sponsor

余额不足监听：

1. 余额不足
2. 对应From账户没有进行中的sponsor交易
