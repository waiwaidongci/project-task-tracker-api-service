package model

const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
)

const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
)

func ValidStatus(status string) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	default:
		return false
	}
}

func ValidPriority(priority string) bool {
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}
}
