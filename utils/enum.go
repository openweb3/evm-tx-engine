package utils

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
