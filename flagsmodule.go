package orc

import (
	"github.com/spf13/pflag"
)

type flagsModule struct {
	setFlags func(*pflag.FlagSet)
}

func (f flagsModule) ModuleName() string { return "Flags" }

func FlagsModule(setFlags func(*pflag.FlagSet)) Module {
	return flagsModule{setFlags}
}

func (f flagsModule) OnRegister(hooks ModuleHooks) {
	// init -- this is only executed once.

	hooks.OnUse(func(ctx UseContext) {
		f.setFlags(ctx.Flags)
	})
}
