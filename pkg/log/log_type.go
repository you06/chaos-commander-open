package log

type Type int

const (
	TypeManager Type = iota
	TypeLoad
	TypeNode
)

// String implement string interface
func (l Type)String() string {
	switch l {
	case TypeManager:
		return "manager log"
	case TypeLoad:
		return "load log"
	case TypeNode:
		return "node log"
	default:
		return "unknown log"
	}
}
