package object

// Log represents a log object
type Log struct {
	Message Object
	Line    uint64
}

// Type .
func (l *Log) Type() ObjectType { return LogObject }

// Inspect null
func (l *Log) Inspect() string { return l.Message.Inspect() }
