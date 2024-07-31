package types

import (
	"context"
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
