package executor

type ExecStatus int

const (
	StatusPending ExecStatus = iota
	StatusPrepare
	StatusRunning
)

func (s ExecStatus) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusPrepare:
		return "prepare"
	case StatusRunning:
		return "running"
	default:
		return "unknown status"
	}
}
