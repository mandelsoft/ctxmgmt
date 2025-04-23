package attributes

import (
	"context"
	"reflect"

	"github.com/mandelsoft/datacontext"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/logging"

	"github.com/mandelsoft/datacontext/action/handlers"
	ctxlog "github.com/mandelsoft/datacontext/logging"
)

// CONTEXT_TYPE is the global type for an attribute context.
const CONTEXT_TYPE = "attributes" + datacontext.CONTEXT_SUFFIX

type AttributesContext = datacontext.AttributesContext

// DefaultContext is the default context initialized by init functions.
var DefaultContext = NewWithActions(nil, handlers.DefaultRegistry())

// ForContext returns the Context to use for context.Context.
// This is either an explicit context or the default context.
func ForContext(ctx context.Context) AttributesContext {
	c, _ := datacontext.ForContextByKey(ctx, key, DefaultContext)
	if c == nil {
		return nil
	}
	return c.(AttributesContext)
}

// WithContext create a new Context bound to a context.Context.
func WithContext(ctx context.Context, parentAttrs datacontext.Attributes) (datacontext.Context, context.Context) {
	c := New(parentAttrs)
	return c, c.BindTo(ctx)
}

////////////////////////////////////////////////////////////////////////////////

var key = reflect.TypeOf(_context{})

type _InternalContext = datacontext.InternalContext

type _context struct {
	_InternalContext
	updater datacontext.Updater
}

var (
	_ datacontext.Context                        = (*_context)(nil)
	_ datacontext.ViewCreator[AttributesContext] = (*_context)(nil)
)

// New provides a root attribute context.
func New(parentAttrs ...datacontext.Attributes) AttributesContext {
	return NewWithActions(general.Optional(parentAttrs...), handlers.NewRegistry(nil, handlers.DefaultRegistry()))
}

func NewWithActions(parentAttrs datacontext.Attributes, actions handlers.Registry) AttributesContext {
	return newWithActions(datacontext.MODE_DEFAULTED, parentAttrs, actions)
}

func newWithActions(mode datacontext.BuilderMode, parentAttrs datacontext.Attributes, actions handlers.Registry) AttributesContext {
	c := &_context{}

	c._InternalContext = datacontext.NewContextBase(c, CONTEXT_TYPE, key, parentAttrs, datacontext.ComposeDelegates(ctxlog.Context(), actions))
	/*
	       c.internalContext = newContextBase(c, CONTEXT_TYPE, key, parentAttrs, &c.updater,
	   		ComposeDelegates(logging.NewWithBase(ctxlog.Context()), handlers.NewRegistry(nil, actions)),
	   	)
	*/
	return datacontext.SetupContext(mode, c.CreateView()) // see above
}

func (c *_context) Update() error {
	return c.updater.Update()
}

func (c *_context) CreateView() AttributesContext {
	return newView(c, true)
}

func (c *_context) AttributesContext() AttributesContext {
	if c.updater != nil {
		c.updater.Update()
	}
	return newView(c)
}

func (c *_context) IsAttributesContext() bool {
	return true
}

func (c *_context) Actions() handlers.Registry {
	if c.updater != nil {
		c.updater.Update()
	}
	return c._InternalContext.GetActions()
}

func (c *_context) LoggingContext() logging.Context {
	if c.updater != nil {
		c.updater.Update()
	}
	return c._InternalContext.LoggingContext()
}

func (c *_context) Logger(messageContext ...logging.MessageContext) logging.Logger {
	if c.updater != nil {
		c.updater.Update()
	}
	return c._InternalContext.Logger(messageContext...)
}

////////////////////////////////////////////////////////////////////////////////

// gcWrapper is used as garbage collectable
// wrapper for a context implementation
// to establish a runtime finalizer.
type gcWrapper struct {
	datacontext.GCWrapper
	*_context
}

func newView(c *_context, ref ...bool) AttributesContext {
	if general.Optional(ref...) {
		return datacontext.FinalizedContext[gcWrapper](c)
	}
	return c
}

func (w *gcWrapper) SetContext(c *_context) {
	w._context = c
}

// AssureUpdater is used to assure the existence of an updater in
// a root context if a config context is down the context hierarchy.
// This method SHOULD only be called by a config context.
func AssureUpdater(attrs datacontext.Context, u datacontext.Updater) {
	c, ok := attrs.(*gcWrapper)
	if !ok {
		return
	}
	if c.updater == nil {
		c.updater = u
	}
}
