package orc

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type RunContext struct {
	modulesBeganUse []*moduleContext
	orderedModules  []*moduleContext
	graph           Graph
}

func (r *RunContext) registerEdge(dependent, dependency Module) {
	logrus.Infof("[%p] Dependency: %s --> %s", r, Describe(dependent), Describe(dependency))
	r.graph.Edges = append(r.graph.Edges, Edge{
		Dependent:  dependent,
		Dependency: dependency,
	})
}

type UseContext struct {
	Flags *pflag.FlagSet

	runCtx       *RunContext
	modCtx       *moduleContext
	originalOpts []Option
	alreadyUsed  map[Module]bool
}

func (u UseContext) Use(mod Module, newOptions ...Option) {
	u.runCtx.registerEdge(u.modCtx.module, mod)
	u.modCtx.use(u.runCtx, mod, append(u.originalOpts, newOptions...)...)
}

type ModuleHooks interface {
	OnUse(func(UseContext))

	OnValidate(func() error)
	OnStart(func() error)
	OnStop(func() error)
}

type hookId string

const (
	hookValidate = hookId("validate")
	hookStart    = hookId("start")
	hookStop     = hookId("stop")
)

type labelledHook struct {
	hookType hookId
	hook     func() error
}

type moduleContext struct {
	reg      *Registry
	module   Module
	usehooks []func(UseContext)
	hooks    []labelledHook
}

func (m *moduleContext) addHook(stage hookId, f func() error) {
	m.hooks = append(m.hooks, labelledHook{stage, f})
}

func (m *moduleContext) OnUse(f func(UseContext)) {
	m.usehooks = append(m.usehooks, f)
}

func (m *moduleContext) OnValidate(f func() error) {
	m.addHook(hookValidate, f)
}

func (m *moduleContext) OnStart(f func() error) {
	m.addHook(hookStart, f)
}

func (m *moduleContext) OnStop(f func() error) {
	m.addHook(hookStop, f)
}

func (m *moduleContext) runUseHooks(ctx UseContext) {
	if ctx.alreadyUsed[m.module] {
		return
	}
	for _, f := range m.usehooks {
		f(ctx)
	}
}

type Registry struct {
	mu      sync.Mutex
	modules map[Module]*moduleContext
}

var (
	globalRegistry = &Registry{
		modules: map[Module]*moduleContext{},
	}
)

func (r *Registry) getContext(module Module, options ...Option) (*moduleContext, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx, ok := r.modules[module]
	if !ok {
		ctx = &moduleContext{
			reg:    r,
			module: module,
		}
		r.modules[module] = ctx
	}

	return ctx, ok
}

func (r *Registry) use(runCtx *RunContext, parent Module, module Module, options ...Option) {
	// Is the module registered? If not register it, otherwise get its context.
	modCtx, wasOld := r.getContext(module)
	if !wasOld {
		module.OnRegister(modCtx)
	}

	parsed, err := parseOptions(options)
	if err != nil {
		Fail(err)
	}

	useCtx := UseContext{
		runCtx:       runCtx,
		modCtx:       modCtx,
		originalOpts: options,
		alreadyUsed:  map[Module]bool{},
	}
	for _, mod := range runCtx.modulesBeganUse {
		useCtx.alreadyUsed[mod.module] = true
	}
	// Prevent use hooks from running twice when used in the same RunContext.
	runCtx.modulesBeganUse = append(runCtx.modulesBeganUse, modCtx)

	parsed.intoUseContext(&useCtx)
	modCtx.runUseHooks(useCtx)

	runCtx.orderedModules = append(runCtx.orderedModules, modCtx)
}

func (m *moduleContext) use(runCtx *RunContext, mod Module, options ...Option) {
	// TODO check cycles. may be recursive, but not cyclic.
	globalRegistry.use(runCtx, m.module, mod, options...)
}

func Use(module Module, options ...Option) *RunContext {
	// TODO check this is a top-level call (never recursive)
	runCtx := &RunContext{}
	logrus.Infof("[%p] Top-level use: %s", runCtx, Describe(module))
	globalRegistry.use(runCtx, nil, module, options...)
	return runCtx
}

func (r *RunContext) Run(body func() error) {
	if err := globalRegistry.tryRun(r, body); err != nil {
		logrus.Fatalf("fatal: %v", err)
	}
}

func (r *Registry) tryRun(runCtx *RunContext, body func() error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Visit all the modules in order, run their validate hooks.
	for _, modCtx := range runCtx.orderedModules {
		for _, labelledHook := range modCtx.hooks {
			if labelledHook.hookType == hookValidate {
				if err := labelledHook.hook(); err != nil {
					return err
				}
			}
		}
	}

	// Visit all the modules in order, run start and defer stop.
	for _, modCtx := range runCtx.orderedModules {
		for _, labelledHook := range modCtx.hooks {
			if labelledHook.hookType == hookStart {
				if err := labelledHook.hook(); err != nil {
					return err
				}
			}
			if labelledHook.hookType == hookStop {
				labelledHook := labelledHook
				defer func() {
					err := labelledHook.hook()
					if err != nil {
						logrus.Warningf("error: %v", err)
					}
				}()
			}
		}
	}

	return body()
}
