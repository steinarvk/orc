package orc

import (
	"github.com/sirupsen/logrus"
)

type Module interface {
	OnRegister(ModuleHooks)
}

type Runnable interface {
	Run(func() error)
}

func Fail(err error) {
	logrus.Fatalf("fatal: internal error: orc setup error: %v", err)
}
