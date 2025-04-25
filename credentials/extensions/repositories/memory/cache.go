package memory

import (
	"sync"

	"github.com/mandelsoft/ctxmgmt"
)

const ATTR_REPOS = "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/memory"

type Repositories struct {
	lock  sync.Mutex
	repos map[string]*Repository
}

func newRepositories(ctxmgmt.Context) interface{} {
	return &Repositories{
		repos: map[string]*Repository{},
	}
}

func (r *Repositories) GetRepository(name string) *Repository {
	r.lock.Lock()
	defer r.lock.Unlock()
	repo := r.repos[name]
	if repo == nil {
		repo = NewRepository(name)
		r.repos[name] = repo
	}
	return repo
}
