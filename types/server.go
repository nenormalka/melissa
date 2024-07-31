package types

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	waitingTime = 300 * time.Millisecond
)

type (
	ServerList []Runnable

	ServerPool struct {
		p []Runnable
	}
)

func NewServerPool(sl ServerList) *ServerPool {
	return &ServerPool{
		p: sl,
	}
}

func (p *ServerPool) Start(ctx context.Context) error {
	return start(ctx, p.p, runnableServer)
}

func (p *ServerPool) Stop(ctx context.Context) {
	stop(ctx, p.p, runnableServer)
}

func StartServerWithWaiting(
	ctx context.Context,
	f func(errCh chan error),
) error {
	errCh := make(chan error)
	ctxT, cancel := context.WithTimeout(ctx, waitingTime)
	defer cancel()
	defer func() {
		go func() {
			for err := range errCh {
				if err != nil {
					slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
						Level:     slog.LevelError,
						AddSource: true,
					})).Error("server err: %v", err)
				}
			}
		}()
	}()

	go f(errCh)

	select {
	case <-ctxT.Done():
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

func CheckAddr(addr string) string {
	if addr == "" || strings.Contains(addr, ":") {
		return addr
	}

	return ":" + addr
}
