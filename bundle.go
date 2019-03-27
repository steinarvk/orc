package orc

type ModuleBundle struct {
	Modules []Module
}

func (m *ModuleBundle) ModuleName() string {
	return "Modules(...)"
}

func Modules(ms ...Module) *ModuleBundle {
	return &ModuleBundle{ms}
}

func (m *ModuleBundle) OnRegister(h ModuleHooks) {
	h.OnUse(func(u UseContext) {
		for _, mod := range m.Modules {
			u.Use(mod)
		}
	})
}
