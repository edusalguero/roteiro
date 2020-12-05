package shutdown

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/edusalguero/roteiro.git/internal/utils/gox"
)

// Stopper defines a dependency that can perform graceful shutdown
type Stopper interface {
	Stop(ctx context.Context)
}

// StartStopper is a Stopper that can be started
type StartStopper interface {
	Stopper
	Start() error
}

const startTimeout = time.Minute

type phase struct {
	p Phase
	sync.Mutex
	stoppers []Stopper
}

func (p *phase) AfterStarting(ss StartStopper) {
	ctx, cancel := context.WithTimeout(context.Background(), startTimeout)
	defer cancel()

	err := gox.CallWithContext(ctx, ss.Start)
	if err != nil {
		if err == context.DeadlineExceeded {
			err = fmt.Errorf("component %T did not start in the allocated time (%s)", ss, startTimeout)
		}
		log.Panic(err)
	}

	p.Register(ss)
}

func (p *phase) Register(s Stopper) {
	p.Lock()
	p.stoppers = append(p.stoppers, s)
	p.Unlock()
}

func (p *phase) RegisterFunc(fn func(context.Context)) {
	p.Lock()
	p.stoppers = append(p.stoppers, stopFunc(fn))
	p.Unlock()
}

type stopFunc func(ctx context.Context)

func (stop stopFunc) Stop(ctx context.Context) {
	stop(ctx)
}
