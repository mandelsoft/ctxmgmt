package config

import (
	"context"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/attributes"
	"github.com/mandelsoft/ctxmgmt/config/internal"
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

func New(mode ...ctxmgmt.BuilderMode) Context {
	return internal.Builder{}.New(mode...)
}
