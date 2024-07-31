package types

import (
	"context"
	"fmt"
	"sync"

	log "github.com/nenormalka/melissa/logger"
)

type (
	Runnable interface {
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
	}

	runnableType string
)

const (
	runnableServer  runnableType = "server"
	runnableService runnableType = "service"
)

func start(ctx context.Context, pool []Runnable, name runnableType) error {
	if len(pool) == 0 {
		return nil
	}

	logger := log.NewLogger()

	logger.Info(fmt.Sprintf("Graceful: Run %s...", name))

	for _, s := range pool {
		if err := s.Start(ctx); err != nil {
			return fmt.Errorf("run %s... err: %w", name, err)
		}

		logger.Info(fmt.Sprintf("Graceful: %s `%T` started", name, s))
	}

	logger.Info(fmt.Sprintf("Graceful: All %s run", name))

	return nil
}

func stop(ctx context.Context, pool []Runnable, name runnableType) {
	if len(pool) == 0 {
		return
	}

	logger := log.NewLogger()

	logger.Info(fmt.Sprintf("Graceful: Stopping %s...", name))

	var wg sync.WaitGroup

	wg.Add(len(pool))

	for _, r := range pool {
		go func(r Runnable) {
			defer wg.Done()

			if err := r.Stop(ctx); err != nil {
				logger.Error(fmt.Sprintf("stop %s... err %s", name, err.Error()))

				return
			}

			logger.Info(fmt.Sprintf("Graceful: %s `%T` stoped", name, r))
		}(r)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info(fmt.Sprintf("Graceful: All %s stopped", name))
	case <-ctx.Done():
		logger.Error("Graceful: Stop aborted by context done")
	}
}
