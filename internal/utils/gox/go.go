package gox

import (
	"context"
)

// Recover creates a goroutine executing the function given and recovers from any panic in it.
func Recover(f func()) { Impl.Recover(f) }

// Report creates a goroutine executing the function given and reports any panic in it (and then probably your program dies)
func Report(f func()) { Impl.Report(f) }

// CallWithContext calls the given function in a separate goroutine, and returns
// when either it returns or the given context expires.
//
// It returns the error from the function or the context error if the context
// expires.
//
// If the new goroutine panics, the recovered value is wrapped in a
// RecoveredPanic that carries the original stack trace and repanicked.
func CallWithContext(ctx context.Context, f func() error) error { return Impl.CallWithContext(ctx, f) }
