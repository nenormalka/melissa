package melissa

import (
	"log"

	"go.uber.org/dig"
)

type (
	Engine struct {
		container *dig.Container
		mainFunc  any
	}
)

func NewEngine(mainFunc any, modules Module) *Engine {
	e := &Engine{
		container: dig.New(),
		mainFunc:  mainFunc,
	}

	e.provide(modules)

	return e
}

func (e *Engine) provide(m Module) {
	for _, c := range m {
		if err := e.container.Provide(c.CreateFunc, c.Options...); err != nil {
			log.Fatalf("provide err %s", err.Error())
		}
	}
}

func (e *Engine) Run() {
	if err := e.container.Invoke(e.mainFunc); err != nil {
		log.Fatalf("invoke err %s", err.Error())
	}
}
