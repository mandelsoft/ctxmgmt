package attributes

import (
	"context"

	"github.com/mandelsoft/datacontext"
	"github.com/mandelsoft/datacontext/action/api"
	"github.com/mandelsoft/datacontext/action/handlers"
)

type Builder struct {
	ctx        context.Context
	attributes datacontext.Attributes
	actions    handlers.Registry
}

func (b *Builder) getContext() context.Context {
	if b.ctx == nil {
		return context.Background()
	}
	return b.ctx
}

func (b Builder) WithContext(ctx context.Context) Builder {
	b.ctx = ctx
	return b
}

func (b Builder) WithAttributes(parentAttr datacontext.Attributes) Builder {
	b.attributes = parentAttr
	return b
}

func (b Builder) WithActionHandlers(hdlrs handlers.Registry) Builder {
	b.actions = hdlrs
	return b
}

func (b Builder) Bound() (datacontext.Context, context.Context) {
	c := b.New()
	return c, context.WithValue(b.getContext(), key, c)
}

func (b Builder) New(m ...datacontext.BuilderMode) datacontext.Context {
	mode := datacontext.Mode(m...)

	if b.actions == nil {
		switch mode {
		case datacontext.MODE_INITIAL:
			b.actions = handlers.NewRegistry(api.NewActionTypeRegistry())
		case datacontext.MODE_CONFIGURED:
			b.actions = handlers.NewRegistry(api.DefaultRegistry().Copy())
			handlers.DefaultRegistry().AddTo(b.actions)
		case datacontext.MODE_EXTENDED:
			b.actions = handlers.NewRegistry(api.DefaultRegistry(), handlers.DefaultRegistry())
		case datacontext.MODE_DEFAULTED:
			fallthrough
		case datacontext.MODE_SHARED:
			b.actions = handlers.DefaultRegistry()
		}
	}

	return newWithActions(mode, b.attributes, b.actions)
}
