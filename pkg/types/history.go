package types

import "time"

type HistoryStatus int

const (
	HistoryStatusSuccess HistoryStatus = iota
	HistoryStatusFailed
	HistoryStatusPending
)

// History defines history table structure
type History struct {
	ID         int           `gorm:"column:id"`
	JobID      int           `gorm:"column:job_id"`
	Status     HistoryStatus `gorm:"column:status"`
	Log        string        `gorm:"column:log"`
	CreatedAt  time.Time     `gorm:"column:created_at"`
	FinishedAt time.Time     `gorm:"column:finished_at"`
}

func (h *History) GetID() int {
	return h.ID
}

func (h *History) GetJobID() int {
	return h.JobID
}

func (h *History) GetStatus() HistoryStatus {
	return h.Status
}

func (h *History) GetLog() string {
	return h.Log
}

func (s HistoryStatus) String() string {
	switch s {
	case HistoryStatusSuccess:
		return "success"
	case HistoryStatusFailed:
		return "failed"
	case HistoryStatusPending:
		return "pending"
	default:
		return "unknown history status"
	}
}
