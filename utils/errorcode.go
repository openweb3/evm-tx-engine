package utils

type TxErrorCode uint

const (
	TxInternalErrorCode           TxErrorCode = 80100 // General Tx
	TxInternalInsufficientBalance TxErrorCode = 80101

	// TODO: consider discarded ———— if a transaction were discarded, sometimes it means something else. How can we tell the transactions' status if it is not possible for us to get the full pending mempool?
	TxPoolErrorCode           TxErrorCode = 80200
	TxPoolDiscarded           TxErrorCode = 80201
	TxPoolInsufficientBalance TxErrorCode = 80202 // will happen?
	TxPoolLowPrice            TxErrorCode = 80203 // check current chain gas price and compare that with the transaction
	TxPoolLowGasLimit         TxErrorCode = 80204 // will happen ?
	// TxPoolFutureNonce         TxErrorCode = 80205 // should be listed here ?

	TxChainErrorCode        TxErrorCode = 80300
	TxChainExecutionFailure TxErrorCode = 80301
)
