package paniccatcher

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	log "github.com/cabify/go-logging"
	"github.com/google/uuid"
)

type PanicReporterFunc func(id, stackTrace, err, message string)

// PanicReporter is called during a panic to report panics.
var PanicReporter PanicReporterFunc

// Catcher defers Default.Catcher and repanics.
//
// See Interface.Catcher.
func Catcher() {
	defer Default.Catcher()
	if r := recover(); r != nil {
		panic(r)
	}
}

// Reporter defers Default.Reporter and repanics.
//
// See Interface.Reporter.
func Reporter() {
	defer Default.Reporter()
	if r := recover(); r != nil {
		panic(r)
	}
}

// CatcherWithContext defers Default.CatcherWithContext and repanics.
//
// See Interface.CatcherWithContext.
func CatcherWithContext(contextData func() string, afterRecovery func()) {
	defer Default.CatcherWithContext(contextData, afterRecovery)
	if r := recover(); r != nil {
		panic(r)
	}
}

// LogWithContext calls Default.LogWithContext.
//
// See Interface.LogWithContext.
func LogWithContext(err interface{}, context string) {
	Default.LogWithContext(err, context)
}

// Interface is the interface that this package exposes.
type Interface interface {
	// Catcher can be deferred in any place in order to recover gracefully from a panic, logging it in our logging systems
	Catcher()

	// Reporter can be deferred in any place in order to report panics, logging it in our logging systems, but then panicking again
	Reporter()

	// CatcherWithContext can be deferred just like Catcher, but will invoke the provided function
	// in order to recover more information about the context of the panic. If successfully recovered from a panic,
	// it will call the afterRecovery function
	CatcherWithContext(contextData func() string, afterRecovery func())

	// LogWithContext logs the provided panic and the context data (if not empty)
	LogWithContext(err interface{}, context string)
}

// Default is a default Interface instance that uses log.DefaultLogger.
var Default Interface

func init() {
	Default = With(log.DefaultLogger)
}

// Logger is the subset of log.Logger that this package may depend on.
type Logger interface {
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
	Criticalln(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Warningln(args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Noticeln(args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})
}

// With returns an Interface that uses the provided Logger to log things on panics.
//
// Use Default if you don't want to provide your own Logger.
func With(log Logger) Interface {
	return withLogger{log: log}
}

type withLogger struct {
	log Logger
}

func (p withLogger) Catcher() {
	if err := recover(); err != nil {
		p.LogWithContext(err, "")
	}
}

func (p withLogger) Reporter() {
	if err := recover(); err != nil {
		p.LogWithContext(err, "")
		panic(err)
	}
}

func (p withLogger) CatcherWithContext(contextData func() string, afterRecovery func()) {
	if err := recover(); err != nil {
		p.LogWithContext(err, strings.Trim(contextData(), "\n"))
		afterRecovery()
	}
}

func (p withLogger) LogWithContext(err interface{}, context string) {
	errString := fmt.Sprintf("%s", err)
	panicID := panicID()

	stackTrace := string(skipStack(5))

	// Logs
	errorMsg := stackTrace
	if context != "" {
		errorMsg += "\n" + context
	}
	p.log.Critical(addPanicIDScope(panicID, fmt.Sprintf("%s\n%s", errString, errorMsg)))

	if PanicReporter != nil {
		PanicReporter(panicID, stackTrace, errString, context)
	}
}

func panicID() string {
	s, err := uuid.NewRandom()
	if err != nil {
		return ""
	}
	return s.String()[:6] // 6 should be long enough.
}

// addPanicIDScope adds a "panic:<id>: " prefix to each line, for logs.
func addPanicIDScope(panicID, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("panic:%s: %s", panicID, line)
	}
	return strings.Join(lines, "\n")
}

func skipStack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return []byte("???")
	}
	return bytes.TrimSpace(lines[n])
}

// funcion gets function name without package and similar cruft
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return []byte("???")
	}
	name := []byte(fn.Name())
	if lastslash := bytes.LastIndex(name, []byte("/")); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, []byte(".")); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, []byte("·"), []byte("."), -1)
	return name
}
