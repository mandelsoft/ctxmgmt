package cpi

// This is the Context Provider Interface for credential providers

import (
	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/credentials/internal"
	"github.com/mandelsoft/ctxmgmt/utils"
)

const (
	KIND_CREDENTIALS = internal.KIND_CREDENTIALS
	KIND_REPOSITORY  = internal.KIND_REPOSITORY
)

const CONTEXT_TYPE = internal.CONTEXT_TYPE

type (
	Context                = internal.Context
	ContextProvider        = internal.ContextProvider
	Repository             = internal.Repository
	RepositoryType         = internal.RepositoryType
	RepositoryTypeProvider = internal.RepositoryTypeProvider
	RepositoryTypeScheme   = internal.RepositoryTypeScheme
	Credentials            = internal.Credentials
	CredentialsSource      = internal.CredentialsSource
	CredentialsChain       = internal.CredentialsChain
	CredentialsSpec        = internal.CredentialsSpec
	RepositorySpec         = internal.RepositorySpec
	GenericRepositorySpec  = internal.GenericRepositorySpec
	GenericCredentialsSpec = internal.GenericCredentialsSpec
	DirectCredentials      = internal.DirectCredentials
	EvaluationContext      = internal.EvaluationContext
)

type (
	ConsumerIdentity         = internal.ConsumerIdentity
	ConsumerIdentityProvider = internal.ConsumerIdentityProvider
	ProviderIdentity         = internal.ProviderIdentity
	ConsumerProvider         = internal.ConsumerProvider
	UsageContext             = internal.UsageContext
	StringUsageContext       = internal.StringUsageContext
	IdentityMatcher          = internal.IdentityMatcher
	IdentityMatcherInfo      = internal.IdentityMatcherInfo
	IdentityMatcherRegistry  = internal.IdentityMatcherRegistry
)

var DefaultContext = internal.DefaultContext

func FromProvider(p ContextProvider) Context {
	return internal.FromProvider(p)
}

func New(m ...ctxmgmt.BuilderMode) Context {
	return internal.Builder{}.New(m...)
}

func NewConsumerIdentity(typ string, attrs ...string) ConsumerIdentity {
	return internal.NewConsumerIdentity(typ, attrs...)
}

func NewGenericCredentialsSpec(name string, repospec *GenericRepositorySpec) *GenericCredentialsSpec {
	return internal.NewGenericCredentialsSpec(name, repospec)
}

func NewCredentialsSpec(name string, repospec RepositorySpec) CredentialsSpec {
	return internal.NewCredentialsSpec(name, repospec)
}

func ToGenericCredentialsSpec(spec CredentialsSpec) (*GenericCredentialsSpec, error) {
	return internal.ToGenericCredentialsSpec(spec)
}

func ToGenericRepositorySpec(spec RepositorySpec) (*GenericRepositorySpec, error) {
	return internal.ToGenericRepositorySpec(spec)
}

func RegisterStandardIdentityMatcher(typ string, matcher IdentityMatcher, desc string) {
	internal.StandardIdentityMatchers.Register(typ, matcher, desc)
}

func RegisterStandardIdentity(typ string, matcher IdentityMatcher, desc string, attrs string) {
	internal.StandardIdentityMatchers.Register(typ, matcher, desc, attrs)
}

func NewCredentials(props utils.Properties) Credentials {
	return internal.NewCredentials(props)
}

func ErrUnknownCredentials(name string) error {
	return internal.ErrUnknownCredentials(name)
}

func ErrUnknownRepository(kind, name string) error {
	return internal.ErrUnknownRepository(kind, name)
}

func CredentialsForConsumer(ctx ContextProvider, id ConsumerIdentity, matchers ...IdentityMatcher) (Credentials, error) {
	return internal.CredentialsForConsumer(ctx, id, false, matchers...)
}

func RequiredCredentialsForConsumer(ctx ContextProvider, id ConsumerIdentity, matchers ...IdentityMatcher) (Credentials, error) {
	return internal.CredentialsForConsumer(ctx, id, true, matchers...)
}

func GetCredentialsForConsumer(ctx Context, ectx EvaluationContext, identity ConsumerIdentity, matchers ...IdentityMatcher) (CredentialsSource, error) {
	return internal.GetCredentialsForConsumer(ctx, ectx, identity, matchers...)
}

func GetEvaluationContextFor[T any](ectx EvaluationContext) T {
	return internal.GetEvaluationContextFor[T](ectx)
}

func SetEvaluationContextFor(ectx EvaluationContext, e any) {
	internal.SetEvaluationContextFor(ectx, e)
}

func SimpleCredentials(user, passwd string) Credentials {
	return internal.SimpleCredentials(user, passwd)
}

var (
	CompleteMatch = internal.CompleteMatch
	NoMatch       = internal.NoMatch
	PartialMatch  = internal.PartialMatch
)

// provide context interface for other files to avoid diffs in imports.
var (
	newStrictRepositoryTypeScheme = internal.NewStrictRepositoryTypeScheme
	defaultRepositoryTypeScheme   = internal.DefaultRepositoryTypeScheme
)
