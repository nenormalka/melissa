package main

import (
	"context"
	"log"
	"time"

	"github.com/nenormalka/melissa"
	"github.com/nenormalka/melissa/types"

	"go.uber.org/dig"
)

type (
	QueueConnectionImitator struct {
		ch      chan struct{}
		closeCh chan struct{}
	}

	SimpleService struct {
		ch <-chan struct{}
	}

	AdapterOut struct {
		dig.Out

		Service types.Runnable `group:"services"`
	}
)

var Module = melissa.Module{
	{CreateFunc: NewQueueConnectionImitator},
	{CreateFunc: NewSimpleService},
	{CreateFunc: Adapter},
}

func main() {
	melissa.NewEngine(MainFunc, Module).Run()
}

func MainFunc(
	ctx context.Context,
	qci *QueueConnectionImitator,
	app *melissa.App,
) {
	defer qci.Stop()

	if err := app.Run(ctx); err != nil {
		panic(err)
	}
}

func Adapter(s *SimpleService) AdapterOut {
	return AdapterOut{Service: s}
}

func NewSimpleService(qci *QueueConnectionImitator) *SimpleService {
	return &SimpleService{ch: qci.Channel()}
}

func (ss *SimpleService) Start(ctx context.Context) error {
	log.Println("SimpleService started")
	go func() {
		for {
			select {
			case <-ss.ch:
				log.Println("SimpleService: got message")
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (ss *SimpleService) Stop(_ context.Context) error {
	log.Println("SimpleService stopped")
	return nil
}

func NewQueueConnectionImitator() *QueueConnectionImitator {
	qci := &QueueConnectionImitator{
		ch:      make(chan struct{}),
		closeCh: make(chan struct{}),
	}

	qci.Start()

	return qci
}

func (qci *QueueConnectionImitator) Start() {
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()

		for {
			<-t.C

			select {
			case qci.ch <- struct{}{}:
			case <-qci.closeCh:
				return
			}
		}
	}()
}

func (qci *QueueConnectionImitator) Stop() {
	close(qci.closeCh)
	close(qci.ch)
}

func (qci *QueueConnectionImitator) Channel() <-chan struct{} {
	return qci.ch
}
