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
