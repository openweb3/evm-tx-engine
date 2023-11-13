package code

// NOTE: in current version task status is not in use
// TaskStatus is the status used for user interface
// from the user's perspective, he/she could not know the actual task status from the field, but can somewhat know
// type TaskStatus struct {
// 	Code uint
// 	Name string
// }

type TaskStatus uint

const (
	TaskWaiting       TaskStatus = 90100
	TaskProcessing    TaskStatus = 90200
	TaskSuccess       TaskStatus = 90300
	TaskStableSuccess TaskStatus = 90309
	TaskFailure       TaskStatus = 90400
	TaskStableFailure TaskStatus = 90409

	TaskCancelling    TaskStatus = 91100
	TaskCancelled     TaskStatus = 91200
	TaskCancelFailure TaskStatus = 91300

	TaskUnexpected TaskStatus = 99900
)

func (ts TaskStatus) Name() string {
	switch ts {
	case TaskWaiting:
		return "WAITING"
	case TaskProcessing:
		return "PROCESSING"
	case TaskSuccess:
		return "SUCCESS"
	case TaskStableSuccess:
		return "STABLE_SUCCESS"
	case TaskFailure:
		return "FAILURE"
	case TaskStableFailure:
		return "STABLE_FAILURE"
	case TaskCancelling:
		return "CANCELLING"
	case TaskCancelled:
		return "CANCELLED"
	case TaskCancelFailure:
		return "CANCEL_FAILURE"
	case TaskUnexpected:
		return "UNEXPECTED"
	default:
		return "UNKNOWN"
	}
}
