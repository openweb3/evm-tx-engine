package code

type TxStatus uint

const (
	TxInternalTargetQueue    TxStatus = 80100
	TxInternalGasEnoughQueue TxStatus = 80200
	TxInternalConstructed    TxStatus = 80300
	TxInternalSigned         TxStatus = 80400
	TxInternalError          TxStatus = 80500

	TxPoolPending TxStatus = 81100
	TxPoolError   TxStatus = 81200

	TxChainLatest         TxStatus = 82100
	TxChainLatestError    TxStatus = 82101
	TxChainSafe           TxStatus = 82200
	TxChainSafeError      TxStatus = 82201
	TxChainFinalized      TxStatus = 82300
	TxChainFinalizedError TxStatus = 82301
)

func (ts TxStatus) Name() string {
	switch ts {
	case TxInternalTargetQueue:
		return "TARGET_QUEUE"
	case TxInternalGasEnoughQueue:
		return "GAS_ENOUGH_QUEUE"
	case TxInternalConstructed:
		return "CONSTRUCTED"
	case TxInternalSigned:
		return "SIGNED"
	case TxInternalError:
		return "INTERNAL_ERROR"
	case TxPoolPending:
		return "PENDING"
	case TxPoolError:
		return "POOL_ERROR"
	case TxChainLatest:
		return "LATEST"
	case TxChainLatestError:
		return "LATEST_ERROR"
	case TxChainSafe:
		return "SAFE"
	case TxChainSafeError:
		return "SAFE_ERROR"
	case TxChainFinalized:
		return "FINALIZED"
	case TxChainFinalizedError:
		return "FINALIZED_ERROR"
	default:
		return "UNKNOWN"
	}
}

// // NOTE: in current version, tx_status is not fully in use
// const (
// 	TxInternalTargetQueue    TxStatus = "TARGET_QUEUE"
// 	TxInternalGasEnoughQueue TxStatus = "GAS_ENOUGH_QUEUE"
// 	TxInternalConstructed    TxStatus = "CONSTRUCTED"
// 	TxInternalSigned         TxStatus = "SIGNED"
// 	TxInternalError          TxStatus = "INTERNAL_ERROR" // inernal error means the transaction never enters nodes' mempool

// 	// Transactions are sent or on chain

// 	TxPoolPending TxStatus = "PENDING"
// 	TxPoolError   TxStatus = "POOL_ERROR" // NOTE: pool error means that the transaction ever appeared in nodes' mempool, for example, low gas price or somewhat discarded. Chances are that gas is not sufficient even in this stage

// 	// TxPoolConflictNonceLatestError    TxStatus = "POOL_LATEST_ERROR"    // any transaction with same nonce is latest
// 	// TxPoolConflictNonceSafeError      TxStatus = "POOL_SAFE_ERROR"      // any transaction with same nonce is safe
// 	// TxPoolConflictNonceFinalizedError TxStatus = "POOL_FINALIZED_ERROR" // any transaction with same nonce is finalized

// 	// Transactions are on chain

// 	TxChainLatest      TxStatus = "LATEST"
// 	TxChainLatestError TxStatus = "LATEST_ERROR"

// 	// TODO: Insert Extra Status here

// 	TxChainSafe      TxStatus = "SAFE"
// 	TxChainSafeError TxStatus = "SAFE_ERROR"

// 	TxChainFinalized      TxStatus = "FINALIZED"
// 	TxChainFinalizedError TxStatus = "FINALIZED_ERROR"
// )

func (status TxStatus) IsError() bool {
	return status == TxInternalError || status == TxChainLatestError || status == TxChainSafeError || status == TxChainFinalizedError
}

// Transaction is not sent to the pool
func (status TxStatus) IsInternal() bool {
	return status == TxInternalTargetQueue || status == TxInternalGasEnoughQueue || status == TxInternalConstructed || status == TxInternalSigned || status == TxInternalError
}

func (status TxStatus) IsInPool() bool {
	return status == TxPoolPending || status == TxPoolError
}

func (status TxStatus) IsOnChain() bool {
	return status == TxChainLatest || status == TxChainLatestError || status == TxChainSafe || status == TxChainFinalized || status == TxChainFinalizedError
}

func (status TxStatus) IsSent() bool {
	return status.IsInPool() || status.IsOnChain()
}

func (status TxStatus) IsStable() bool {
	return status == TxChainFinalizedError || status == TxChainFinalized
}
