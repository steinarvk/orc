package orc

import "github.com/sirupsen/logrus"

type Module interface {
	OnRegister(ModuleHooks)
}

type ModuleBundle struct {
	Modules []Module
}

func Modules(ms ...Module) *ModuleBundle {
	return &ModuleBundle{ms}
}

func (m *ModuleBundle) OnRegister(h ModuleHooks) {
	for _, m := range m.Modules {
		m.OnRegister(h)
	}
}

type Runnable interface {
	Run(func() error)
}

func Fail(err error) {
	logrus.Fatalf("fatal: internal error: orc setup error: %v", err)
}
