package executor

type EventType int

const (
	EventPrepare EventType = iota
	EventAnsibleDownload
	EventAnsibleStop
	EventAnsibleClearLog
	EventAnsibleSwapOff
	EventAnsibleNTP
	EventAnsibleCPU
	EventAnsibleDeploy
	EventAnsibleStart
	EventLoadStart
	EventLoadStop
	EventMetricsStart
	EventMetricsStop
	EventBlueprintStart
	EventBlueprintStop
	EventFetchLog
	EventStop
	EventStopSupervise
)

// String change the event into string
func (e EventType)String() string {
	switch e {
	case EventPrepare:
		return "prepare"
	case EventAnsibleDownload:
		return "ansible download"
	case EventAnsibleStop:
		return "ansible stop"
	case EventAnsibleClearLog:
		return "ansible clear log"
	case EventAnsibleSwapOff:
		return "ansible swapoff"
	case EventAnsibleNTP:
		return "ansible deploy ntp"
	case EventAnsibleCPU:
		return "ansible governor CPU mode"
	case EventAnsibleDeploy:
		return "ansible deploy"
	case EventAnsibleStart:
		return "ansible cluster start"
	case EventLoadStart:
		return "load start"
	case EventLoadStop:
		return "load stop"
	case EventMetricsStart:
		return "metrics start"
	case EventMetricsStop:
		return "metrics stop"
	case EventBlueprintStart:
		return "blue print start"
	case EventBlueprintStop:
		return "blue print stop"
	case EventFetchLog:
		return "fetch log"
	case EventStop:
		return "stop"
	case EventStopSupervise:
		return "stop supervise"
	default:
		return "unknown status"
	}
}

