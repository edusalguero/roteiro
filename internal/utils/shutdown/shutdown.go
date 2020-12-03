package shutdown

import (
	"context"
	"math"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/edusalguero/roteiro.git/internal/utils/paniccatcher"
	log "github.com/sirupsen/logrus"
)

// Phase defines when shutdown components are stopped, they will be stopped in the natural integer order
// Phases are intentionally undefined here: we'd like to define:
// PhaseConsumers = 1
// PhaseTearDown = 2
// PhaseConnections = 3
// And they will fit our needs as this code is being written, but they'll get quickly outdated, just like the initial
// approach with First/Last did: someone will need a phase between Consumers & TearDown.
// We leave the definition of the phases to the user, except for the legacy First() & Last() which still apply,
// Being explicit on the order is the most idiomatic approach for golang, and it takes no time to add 3 constants in your main()
type Phase int64

const (
	// phaseFirst is used for retro-compatibility
	phaseFirst Phase = math.MinInt64

	// phaseLast is used for retro-compatibility
	phaseLast Phase = math.MaxInt64
)

var (
	lock   sync.Mutex
	phases = make(map[Phase]*phase)
)

// Config configures the graceful shutdown params
type Config struct {
	Timeout time.Duration `default:"8s"`
}

// Gracefully will shutdown gracefully the registered phase,
// probably you want to defer shutdown.Gracefully() in your main.go
func Gracefully(cfg Config) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	for _, p := range sortedPhases() {
		select {
		case <-ctx.Done():
			log.Errorf("There's no time to shut down more dependencies, context is already done: %s", ctx.Err())
		default:
			p.Lock()
			shutdown(ctx, p.stoppers)
			p.Unlock()
		}
	}

	log.Infof("Shutdown performed, bye!")
}

// Hook provides handy methods to hook your dep into shutdown process
// You can register a Stopper or a Stopper function to be executed in the shutdown process,
// or you can wrap your starter (or critical starter) so it will be automatically registered once started
type Hook interface {
	Register(Stopper)
	RegisterFunc(func(context.Context))
	AfterStarting(StartStopper)
}

// First provides hooks for dependencies that should be shut down in the first place, this includes
// all kind of event-originators: nsq consumer, http server, etc.
func First() Hook {
	return On(phaseFirst)
}

// Last provides hooks for dependencies that should be stopped in the last place, like db connections, etc.
// This are the kind of dependencies that never originate any request, so they don't care if other dependencies are
// already stopped
func Last() Hook {
	return On(phaseLast)
}

// On returns a Hook for the required Phase, it will register the deps to be stopped on that Phase
func On(p Phase) Hook {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := phases[p]; !ok {
		phases[p] = &phase{p: p}
	}
	return phases[p]
}

// WaitForStopSignal blocks until the process receives a signal to terminate
func WaitForStopSignal() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	log.Infof("Waiting for stop signal...")
	sig := <-sigchan
	log.Infof("Received stop signal: %s. Stopping. Have a nice day!", sig.String())
}

func shutdown(ctx context.Context, stoppers []Stopper) {
	wg := sync.WaitGroup{}
	for _, s := range stoppers {
		wg.Add(1)
		go func(s Stopper) {
			defer paniccatcher.Catcher()
			log.Debugf("Stopping %T", s)
			s.Stop(ctx)
			wg.Done()
		}(s)
	}

	// Close done when all deps have shutdown
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		// Triggered timeout, took too long to finish
		log.Errorf("Did not shutdown gracefully, ctx is Done: %s", ctx.Err())
		return
	case <-done:
		return
	}
}

func sortedPhases() []*phase {
	var phs sorted
	lock.Lock()
	for _, p := range phases {
		phs = append(phs, p)
	}
	lock.Unlock()

	sort.Sort(phs)
	return phs
}

type sorted []*phase

func (s sorted) Len() int {
	return len(s)
}

func (s sorted) Less(i, j int) bool {
	return s[i].p < s[j].p
}

func (s sorted) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
