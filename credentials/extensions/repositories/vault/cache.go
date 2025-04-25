package vault

import (
	"sync"

	"github.com/mandelsoft/goutils/errors"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
)

const ATTR_REPOS = "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/vault"

type Repositories struct {
	lock  sync.Mutex
	repos map[cpi.ProviderIdentity]*Repository
}

func newRepositories(ctxmgmt.Context) interface{} {
	return &Repositories{
		repos: map[cpi.ProviderIdentity]*Repository{},
	}
}

func (r *Repositories) GetRepository(ctx cpi.Context, spec *RepositorySpec) (*Repository, error) {
	var repo *Repository

	if spec.ServerURL == "" {
		return nil, errors.ErrInvalid("server url")
	}
	r.lock.Lock()
	defer r.lock.Unlock()

	var err error
	key := spec.GetKey()
	repo = r.repos[key]
	if repo == nil {
		repo, err = NewRepository(ctx, spec)
		if err == nil {
			r.repos[key] = repo
		}
	}
	return repo, err
}
