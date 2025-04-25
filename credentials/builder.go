package credentials

import (
	"context"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/credentials/internal"
)

func WithContext(ctx context.Context) internal.Builder {
	return internal.Builder{}.WithContext(ctx)
}

func WithConfigs(ctx config.Context) internal.Builder {
	return internal.Builder{}.WithConfig(ctx)
}

func WithRepositoyTypeScheme(scheme RepositoryTypeScheme) internal.Builder {
	return internal.Builder{}.WithRepositoyTypeScheme(scheme)
}

func WithStandardConumerMatchers(matchers internal.IdentityMatcherRegistry) internal.Builder {
	return internal.Builder{}.WithStandardConumerMatchers(matchers)
}

func New(mode ...ctxmgmt.BuilderMode) Context {
	return internal.Builder{}.New(mode...)
}
