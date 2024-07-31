package types

import (
	"context"
)

type (
	ServiceList []Runnable

	ServicePool struct {
		p []Runnable
	}
)

func NewServicePool(sl ServiceList) *ServicePool {
	return &ServicePool{
		p: sl,
	}
}

func (p *ServicePool) Start(ctx context.Context) error {
	return start(ctx, p.p, runnableService)
}

func (p *ServicePool) Stop(ctx context.Context) {
	stop(ctx, p.p, runnableService)
}
