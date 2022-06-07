package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Shutdown struct {
	ctx       context.Context
	cancel    context.CancelFunc
	shutdowns []func()
}

func NewShutdown(ctx context.Context) *Shutdown {
	ctx, cancel := context.WithCancel(ctx)
	shutdown := &Shutdown{
		ctx:    ctx,
		cancel: cancel,
	}
	return shutdown
}

// graceful shutdown context
func (s *Shutdown) Context() context.Context {
	return s.ctx
}

// append need graceful close handler
func (s *Shutdown) AppendGracefulClose(close func()) {
	s.shutdowns = append(s.shutdowns, close)
}

// cancel graceful shutdown with manually
func (s *Shutdown) Cancel() {
	s.cancel()
}

// running graceful shutdown service with blocked
func (s *Shutdown) Serve() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	done := s.ctx.Done()
	running := true
	for running {
		select {
		case <-done:
			running = false
		case <-sigint:
			running = false
		default:
			time.Sleep(time.Millisecond)
		}
	}
	// close modules
	for i := range s.shutdowns {
		s.shutdowns[i]()
	}
	s.cancel()
}
