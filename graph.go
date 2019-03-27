package orc

type Edge struct {
	Dependent  Module
	Dependency Module
}

type Graph struct {
	Edges []Edge
}
