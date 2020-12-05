package gox

import (
	"context"

	"github.com/edusalguero/roteiro.git/internal/utils/paniccatcher"
)

// AsyncRoutine is the default `go` based implementation
type AsyncRoutine struct{}

// Recover creates a goroutine executing the function given and recovers from any panic in it.
func (AsyncRoutine) Recover(f func()) {
	go func() {
		defer paniccatcher.Catcher()
		f()
	}()
}

// Report creates a goroutine executing the function given and reports any panic in it (and then probably your program dies)
func (AsyncRoutine) Report(f func()) {
	go func() {
		defer paniccatcher.Reporter()
		f()
	}()
}

// CallWithContext calls the given function in a separate goroutine, and returns
// when either it returns or the given context expires.
//
// It returns the error from the function or the context error if the context
// expires.
//
// If the new goroutine panics, the recovered value is wrapped in a
// RecoveredPanic that carries the original stack trace and repanicked.
func (AsyncRoutine) CallWithContext(ctx context.Context, f func() error) error {
	done := make(chan RecoveredPanic)

	var err error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- NewRecoveredPanic(r)
			}
		}()

		err = f()
		close(done)
	}()

	select {
	case recovered, ok := <-done:
		if !ok {
			return err
		}
		panic(recovered)
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Routine is the main gox implementation interface
type Routine interface {
	Recover(f func())
	Report(f func())
	CallWithContext(ctx context.Context, f func() error) error
}

//go:generate mockery -inpkg -testonly -case underscore -name Routine
//go:generate mockery -outpkg goxmock -output goxmock -case underscore -name Routine

// Impl is the Routine singleton
var Impl Routine = AsyncRoutine{}
