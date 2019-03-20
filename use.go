package orc

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type RunContext struct {
	orderedModules []*moduleContext
}

type UseContext struct {
	Flags *pflag.FlagSet

	runCtx       *RunContext
	modCtx       *moduleContext
	originalOpts []Option
}

func (u *UseContext) Use(mod Module, newOptions ...Option) {
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

func (r *Registry) use(runCtx *RunContext, module Module, options ...Option) {
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
	}
	parsed.intoUseContext(&useCtx)
	modCtx.runUseHooks(useCtx)

	runCtx.orderedModules = append(runCtx.orderedModules, modCtx)
}

func (m *moduleContext) use(runCtx *RunContext, mod Module, options ...Option) {
	// TODO check cycles. may be recursive, but not cyclic.
	globalRegistry.use(runCtx, mod, options...)
}

func Use(module Module, options ...Option) *RunContext {
	// TODO check this is a top-level call (never recursive)
	runCtx := &RunContext{}
	globalRegistry.use(runCtx, module, options...)
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
