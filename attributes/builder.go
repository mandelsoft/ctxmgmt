package attributes

import (
	"context"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/action/api"
	"github.com/mandelsoft/ctxmgmt/action/handlers"
)

type Builder struct {
	ctx        context.Context
	attributes ctxmgmt.Attributes
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

func (b Builder) WithAttributes(parentAttr ctxmgmt.Attributes) Builder {
	b.attributes = parentAttr
	return b
}

func (b Builder) WithActionHandlers(hdlrs handlers.Registry) Builder {
	b.actions = hdlrs
	return b
}

func (b Builder) Bound() (ctxmgmt.Context, context.Context) {
	c := b.New()
	return c, context.WithValue(b.getContext(), key, c)
}

func (b Builder) New(m ...ctxmgmt.BuilderMode) ctxmgmt.Context {
	mode := ctxmgmt.Mode(m...)

	if b.actions == nil {
		switch mode {
		case ctxmgmt.MODE_INITIAL:
			b.actions = handlers.NewRegistry(api.NewActionTypeRegistry())
		case ctxmgmt.MODE_CONFIGURED:
			b.actions = handlers.NewRegistry(api.DefaultRegistry().Copy())
			handlers.DefaultRegistry().AddTo(b.actions)
		case ctxmgmt.MODE_EXTENDED:
			b.actions = handlers.NewRegistry(api.DefaultRegistry(), handlers.DefaultRegistry())
		case ctxmgmt.MODE_DEFAULTED:
			fallthrough
		case ctxmgmt.MODE_SHARED:
			b.actions = handlers.DefaultRegistry()
		}
	}

	return newWithActions(mode, b.attributes, b.actions)
}
