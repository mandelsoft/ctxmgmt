package internal

import (
	"context"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/attributes"
)

type Builder struct {
	ctx        context.Context
	shared     attributes.AttributesContext
	reposcheme ConfigTypeScheme
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

func (b Builder) WithSharedAttributes(ctx attributes.AttributesContext) Builder {
	b.shared = ctx
	return b
}

func (b Builder) WithConfigTypeScheme(scheme ConfigTypeScheme) Builder {
	b.reposcheme = scheme
	return b
}

func (b Builder) Bound() (Context, context.Context) {
	c := b.New()
	return c, context.WithValue(b.getContext(), key, c)
}

func (b Builder) New(m ...ctxmgmt.BuilderMode) Context {
	mode := ctxmgmt.Mode(m...)
	ctx := b.getContext()

	if b.shared == nil {
		if mode == ctxmgmt.MODE_SHARED {
			b.shared = attributes.ForContext(ctx)
		} else {
			b.shared = attributes.New(nil)
		}
	}
	if b.reposcheme == nil {
		switch mode {
		case ctxmgmt.MODE_INITIAL:
			b.reposcheme = NewConfigTypeScheme(nil)
		case ctxmgmt.MODE_CONFIGURED:
			b.reposcheme = NewConfigTypeScheme(nil)
			b.reposcheme.AddKnownTypes(DefaultConfigTypeScheme)
		case ctxmgmt.MODE_EXTENDED:
			b.reposcheme = NewConfigTypeScheme(nil, DefaultConfigTypeScheme)
		case ctxmgmt.MODE_DEFAULTED:
			fallthrough
		case ctxmgmt.MODE_SHARED:
			b.reposcheme = DefaultConfigTypeScheme
		}
	}
	return ctxmgmt.SetupContext(mode, newContext(b.shared, b.reposcheme, b.shared))
}
