package vault

import (
	"github.com/mandelsoft/ctxmgmt/credentials/identity/vault"
	"github.com/mandelsoft/goutils/errors"

	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	"github.com/mandelsoft/ctxmgmt/credentials/internal"
)

type Repository struct {
	ctx      cpi.Context
	spec     *RepositorySpec
	id       cpi.ConsumerIdentity
	provider *ConsumerProvider
}

var (
	_ cpi.Repository               = (*Repository)(nil)
	_ cpi.ConsumerIdentityProvider = (*Repository)(nil)
)

func NewRepository(ctx cpi.Context, spec *RepositorySpec) (*Repository, error) {
	id, err := vault.GetConsumerId(spec.ServerURL, spec.Namespace, spec.MountPath, spec.Path)
	if err != nil {
		return nil, err
	}
	r := &Repository{
		ctx:  ctx,
		spec: spec,
		id:   id,
	}
	if spec.ServerURL == "" {
		return nil, errors.ErrInvalid("server url")
	}
	r.provider, err = NewConsumerProvider(r)
	if err != nil {
		return nil, err
	}
	if spec.PropgateConsumerIdentity {
		ctx.RegisterConsumerProvider(spec.GetKey(), r.provider)
	}
	return r, err
}

var _ cpi.Repository = &Repository{}

func (r *Repository) ExistsCredentials(name string) (bool, error) {
	return r.provider.ExistsCredentials(name)
}

func (r *Repository) LookupCredentials(name string) (cpi.Credentials, error) {
	return r.provider.LookupCredentials(name)
}

func (r *Repository) WriteCredentials(name string, creds cpi.Credentials) (cpi.Credentials, error) {
	return nil, errors.ErrNotSupported("write", "credentials", Type)
}

func (r *Repository) GetConsumerId(uctx ...internal.UsageContext) internal.ConsumerIdentity {
	return r.id
}

func (r *Repository) GetIdentityMatcher() string {
	return vault.CONSUMER_TYPE
}
