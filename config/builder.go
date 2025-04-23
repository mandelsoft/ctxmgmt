package config

import (
	"context"

	"github.com/mandelsoft/datacontext"
	"github.com/mandelsoft/datacontext/attributes"
	"github.com/mandelsoft/datacontext/config/internal"
)

func WithContext(ctx context.Context) internal.Builder {
	return internal.Builder{}.WithContext(ctx)
}

func WithSharedAttributes(ctx attributes.AttributesContext) internal.Builder {
	return internal.Builder{}.WithSharedAttributes(ctx)
}

func WithConfigTypeScheme(scheme ConfigTypeScheme) internal.Builder {
	return internal.Builder{}.WithConfigTypeScheme(scheme)
}

func New(mode ...datacontext.BuilderMode) Context {
	return internal.Builder{}.New(mode...)
}
