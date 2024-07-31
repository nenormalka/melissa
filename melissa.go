package melissa

import (
	"context"
	"log"

	"github.com/nenormalka/melissa/config"
	"github.com/nenormalka/melissa/logger"
	"github.com/nenormalka/melissa/runnable"

	"github.com/chapsuk/grace"
	"github.com/joho/godotenv"
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

		Services []runnable.Runnable `group:"services"`
		Servers  []runnable.Runnable `group:"servers"`
	}

	ServiceAdapterOut struct {
		dig.Out

		ServiceList runnable.ServiceList
		ServerList  runnable.ServerList
	}
)

func ServiceAdapter(in ServiceAdapterIn) ServiceAdapterOut {
	return ServiceAdapterOut{
		ServiceList: in.Services,
		ServerList:  in.Servers,
	}
}

var defaultModules = Module{
	{CreateFunc: config.NewConfig},
	{CreateFunc: logger.NewLogger},
	{CreateFunc: ServiceAdapter},
	{CreateFunc: NewShutdownContext},
	{CreateFunc: runnable.NewServerPool},
	{CreateFunc: runnable.NewServicePool},
	{CreateFunc: NewApp},
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
	godotenv.Overload()

	if err := e.container.Invoke(e.mainFunc); err != nil {
		log.Fatalf("invoke err %s", err.Error())
	}
}

func (e *Engine) provide(m Module) {
	for _, c := range m {
		if err := e.container.Provide(c.CreateFunc, c.Options...); err != nil {
			log.Fatalf("provide err %s", err.Error())
		}
	}
}

func (m Module) Append(o Module) Module {
	return append(m, o...)
}
