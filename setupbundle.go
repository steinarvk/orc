package orc

type ModuleBundleWithSetup struct {
	setup   func()
	Modules []Module
}

func (m *ModuleBundleWithSetup) ModuleName() string {
	return "ModulesWithSetup(...)"
}

func ModulesWithSetup(setup func(), ms ...Module) *ModuleBundleWithSetup {
	return &ModuleBundleWithSetup{setup, ms}
}

func (m *ModuleBundleWithSetup) OnRegister(h ModuleHooks) {
	h.OnUse(func(u UseContext) {
		m.setup()
		for _, mod := range m.Modules {
			u.Use(mod)
		}
	})
}
