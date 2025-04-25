package dockerconfig

import (
	"sync"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
)

const ATTR_REPOS = "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/dockerconfig"

type Repositories struct {
	lock  sync.Mutex
	repos map[string]*Repository
}

func newRepositories(ctxmgmt.Context) interface{} {
	return &Repositories{
		repos: map[string]*Repository{},
	}
}

func (r *Repositories) GetRepository(ctx cpi.Context, name string, data []byte, propagate bool) (*Repository, error) {
	r.lock.Lock()
	defer r.lock.Unlock()
	var (
		err  error = nil
		repo *Repository
	)
	if name != "" {
		repo = r.repos[name]
	}
	if repo == nil {
		repo, err = NewRepository(ctx, name, data, propagate)
		if err == nil {
			r.repos[name] = repo
		}
	}
	return repo, err
}
