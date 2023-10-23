package utils

import "fmt"

type TxStatus string

const (
	TxInternalTargetQueue    TxStatus = "TARGET_QUEUE"
	TxInternalGasEnoughQueue TxStatus = "GAS_ENOUGH_QUEUE"
	TxInternalConstructed    TxStatus = "CONSTRUCTED"
	TxInternalSigned         TxStatus = "SIGNED"
	TxInternalError          TxStatus = "INTERNAL_ERROR" // inernal error means the transaction never enters nodes' mempool

	// Transactions are sent or on chain

	TxPoolPending TxStatus = "PENDING"
	TxPoolError   TxStatus = "POOL_ERROR" // NOTE: pool error means that the transaction ever appeared in nodes' mempool, for example, low gas price or somewhat discarded. Chances are that gas is not sufficient even in this stage

	// TxPoolConflictNonceLatestError    TxStatus = "POOL_LATEST_ERROR"    // any transaction with same nonce is latest
	// TxPoolConflictNonceSafeError      TxStatus = "POOL_SAFE_ERROR"      // any transaction with same nonce is safe
	// TxPoolConflictNonceFinalizedError TxStatus = "POOL_FINALIZED_ERROR" // any transaction with same nonce is finalized

	// Transactions are on chain

	TxChainLatest      TxStatus = "LATEST"
	TxChainLatestError TxStatus = "LATEST_ERROR"

	// TODO: Insert Extra Status here

	TxChainSafe      TxStatus = "SAFE"
	TxChainSafeError TxStatus = "SAFE_ERROR"

	TxChainFinalized      TxStatus = "FINALIZED"
	TxChainFinalizedError TxStatus = "FINALIZED_ERROR"
)

func InferSentTransactionStatus(meta TxStatusMeta, legacyTxStatus TxStatus, cachedTaggedBlockNumbers TaggedBlockNumbers) (TxStatus, error) {
	// TODO: consider errors
	if !legacyTxStatus.IsSent() {
		return legacyTxStatus,
			fmt.Errorf("transaction status is not sent, but it is sent")
	}
	if meta.BlockNumber == nil {
		return TxPoolPending, nil
	}
	if meta.BlockNumber.ToInt().Uint64() > cachedTaggedBlockNumbers.FinalizedBlock {
		return TxChainFinalized, nil
	}
	if meta.BlockNumber.ToInt().Uint64() > cachedTaggedBlockNumbers.SafeBlock {
		return TxChainSafe, nil
	}
	if meta.BlockNumber.ToInt().Uint64() > cachedTaggedBlockNumbers.LatestBlock {
		return TxChainLatest, nil
	}
	return legacyTxStatus,
		fmt.Errorf("transaction not in expected status")
}

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
