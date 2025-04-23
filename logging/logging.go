package logging

import (
	"encoding/json"
	"sync"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/logging"
	logcfg "github.com/mandelsoft/logging/config"
	"github.com/opencontainers/go-digest"
)

// REALM is used to tag all logging done by this library with the datacontext tag.
// This is also used as message context to configure settings for all
// log output provided by this library.
var REALM = logging.DefineRealm("contexts", "general realm used for the datacontext library.")

type StaticContext struct {
	logging.Context
	applied map[string]struct{}
	lock    sync.Mutex
}

func NewContext(ctx logging.Context, global ...bool) *StaticContext {
	if ctx == nil {
		ctx = logging.DefaultContext()
	}
	if !general.Optional(global...) {
		ctx = ctx.WithContext(REALM)
	}
	return &StaticContext{
		Context: ctx,
		applied: map[string]struct{}{},
	}
}

// Configure applies a configuration once.
// Every config identified by its hash is applied
// only once.
func (s *StaticContext) Configure(config *logcfg.Config, extra ...string) error {
	add := ""
	for _, e := range extra {
		if e != "" {
			add += "/" + e
		}
	}
	data, err := json.Marshal(config)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal log config")
	}
	d := digest.FromBytes(data).String() + add

	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.applied[d]; ok {
		return nil
	}
	s.applied[d] = struct{}{}
	return logcfg.Configure(s.Context, config)
}

// global is a wrapper for the default global log content.
var global = NewContext(nil, true)

// ctx is a wrapper for the default ctx log content.
var ctx = NewContext(nil)

// logContext is the ctx log context.
// It can be replaced by SetContext.
var logContext = ctx

// SetContext sets a new preconfigured context.
// This function should be called prior to any configuration
// to avoid losing them.
func SetContext(ctx logging.Context) {
	logContext = NewContext(ctx)
}

// Context returns the default logging configuration used for this library.
func Context() *StaticContext {
	return logContext
}

// Logger determines a default logger for this given message context
// based on the rule settings for this library.
func Logger(messageContext ...logging.MessageContext) logging.Logger {
	return logContext.Logger(messageContext...)
}

func LogContext(ctx logging.Context, provider ...logging.ContextProvider) logging.Context {
	if ctx != nil {
		return ctx
	}

	for _, p := range provider {
		if p != nil {
			return p.LoggingContext()
		}
	}
	return Context()
}

// Configure applies configuration for the default log context
// provided by this package.
func Configure(config *logcfg.Config, extra ...string) error {
	return logContext.Configure(config, extra...)
}

// ConfigureGlobal applies configuration for the default global log context
// provided by this package.
func ConfigureGlobal(config *logcfg.Config, extra ...string) error {
	return global.Configure(config, extra...)
}

// DynamicLogger gets an unbound logger based on the default library logging context.
func DynamicLogger(messageContext ...logging.MessageContext) logging.UnboundLogger {
	return logging.DynamicLogger(Context(), messageContext...)
}

var (
	contexts []*StaticContext
	lock     sync.Mutex
)

func PushContext(ctx logging.Context) {
	lock.Lock()
	defer lock.Unlock()
	contexts = append(contexts, logContext)
	SetContext(ctx)
}

func PopContext() {
	lock.Lock()
	defer lock.Unlock()
	logContext = contexts[len(contexts)-1]
	contexts = contexts[:len(contexts)-1]
}
