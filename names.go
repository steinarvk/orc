package orc

import "fmt"

type NamedModule interface {
	ModuleName() string
}

func Describe(m Module) string {
	if m == nil {
		return "<root>"
	}
	nm, ok := m.(NamedModule)
	if ok {
		return fmt.Sprintf("%s<%p>", nm.ModuleName(), m)
	}
	return fmt.Sprintf("<%p>", m)
}
