package melissa

import (
	"context"
	"os"

	log "github.com/nenormalka/melissa/logger"
	"github.com/nenormalka/melissa/types"

	"github.com/chapsuk/grace"
	"go.uber.org/dig"
)

type (
	Engine struct {
		container *dig.Container
		mainFunc  any
	}

	Provider struct {
		CreateFunc any
		Options    []dig.ProvideOption
	}

	Module []Provider

	ServiceAdapterIn struct {
		dig.In

		Services []types.Runnable `group:"services"`
		Servers  []types.Runnable `group:"servers"`
	}

	ServiceAdapterOut struct {
		dig.Out

		ServiceList types.ServiceList
		ServerList  types.ServerList
	}
)

func ServiceAdapter(in ServiceAdapterIn) ServiceAdapterOut {
	return ServiceAdapterOut{
		ServiceList: in.Services,
		ServerList:  in.Servers,
	}
}

var defaultModules = Module{
	{CreateFunc: ServiceAdapter},
	{CreateFunc: NewShutdownContext},
	{CreateFunc: NewApp},
	{CreateFunc: types.NewServerPool},
	{CreateFunc: types.NewServicePool},
}

func NewShutdownContext() context.Context {
	return grace.ShutdownContext(context.Background())
}

func NewEngine(mainFunc any, modules Module) *Engine {
	e := &Engine{
		container: dig.New(),
		mainFunc:  mainFunc,
	}

	e.provide(defaultModules.Append(modules))

	return e
}

func (e *Engine) Run() {
	if err := e.container.Invoke(e.mainFunc); err != nil {
		log.NewLogger().Error("invoke err %s", err.Error())
		os.Exit(1)
	}
}

func (e *Engine) provide(m Module) {
	for _, c := range m {
		if err := e.container.Provide(c.CreateFunc, c.Options...); err != nil {
			log.NewLogger().Error("provide err %s", err.Error())
			os.Exit(1)
		}
	}
}

func (m Module) Append(o Module) Module {
	return append(m, o...)
}
