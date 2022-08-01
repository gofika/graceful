package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type VoidCloser func()

func NewShutdown(c context.Context) (ctx context.Context, cancel context.CancelFunc, gracefulClose func(VoidCloser)) {
	ctx, cancel = context.WithCancel(c)
	var shutdowns []any
	gracefulClose = func(fn VoidCloser) {
		shutdowns = append(shutdowns, fn)
	}
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		done := ctx.Done()
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
		for _, shutdown := range shutdowns {
			if closer, ok := shutdown.(VoidCloser); ok {
				closer()
			}
		}
		cancel()
	}()
	return
}
