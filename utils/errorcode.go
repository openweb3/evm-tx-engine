package utils

type TxErrorCode uint

const (
	TxInternalErrorCode           TxErrorCode = 80100 // General Tx
	TxInternalInsufficientBalance TxErrorCode = 80101

	TxPoolErrorCode           TxErrorCode = 80200
	TxPoolInsufficientBalance TxErrorCode = 80201
	TxPoolLowPrice            TxErrorCode = 80202
	TxPoolDiscarded           TxErrorCode = 80203

	TxChainErrorCode        TxErrorCode = 80300
	TxChainExecutionFailure TxErrorCode = 80301
)
