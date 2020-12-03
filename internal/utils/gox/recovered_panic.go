package gox

import (
	"fmt"
	"runtime/debug"
)

// NewRecoveredPanic creates a RecoveredPanic that captures the current stack
// trace.
//
// If recovered is a RecoveredPanic, it is returned as-is, except with the new
// current stack prepended to the previous one.
func NewRecoveredPanic(recovered interface{}) RecoveredPanic {
	if p, ok := recovered.(RecoveredPanic); ok {
		p.Stack = append(debug.Stack(), append([]byte("\nprevious stack:\n\n"), p.Stack...)...)
		return p
	}
	return RecoveredPanic{
		Recovered: recovered,
		Stack:     debug.Stack(),
	}
}

// RecoveredPanic is a value recovered from a panic with the stack trace where
// the panic happened.
type RecoveredPanic struct {
	Recovered interface{}
	Stack     []byte
}

// Error implements error.
func (p RecoveredPanic) Error() string {
	return fmt.Sprintf("panic in different goroutine: %v", p.Recovered)
}
