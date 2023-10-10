package utils

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

	// Transactions are on chain

	TxChainLatest    TxStatus = "LATEST"
	TxChainSafe      TxStatus = "SAFE"
	TxChainFinalized TxStatus = "FINALIZED"
	TxChainError     TxStatus = "ONCHAIN_ERROR"
)

// TaskStatus is the status used for user interface
// from the user's perspective, he/she could not know the actual task status from the field, but can somewhat know
type TaskStatus string

const (
	TaskWaiting       TaskStatus = "WAITING" // Task is in waitlist
	TaskProcessing    TaskStatus = "PROCESSING"
	TaskSuccess       TaskStatus = "SUCCESS" // The transaction succeeds, but the result would revert
	TaskStableSuccess TaskStatus = "STABLE_SUCCESS"
	TaskFailure       TaskStatus = "FAILURE"
	TaskStableFailure TaskStatus = "STABLE_FAILURE"

	TaskCancelling    TaskStatus = "CANCELLING"
	TaskCancelled     TaskStatus = "CANCELLED"
	TaskCancelFailure TaskStatus = "CANCEL_FAILURE"

	TaskUnexpected TaskStatus = "UNEXPECTED" // The most worst case. The internal implementation should avoid this case
)
