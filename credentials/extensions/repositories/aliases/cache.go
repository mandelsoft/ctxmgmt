package aliases

import (
	"sync"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
)

const ATTR_REPOS = "github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/aliases"

type Repositories struct {
	sync.RWMutex
	repos map[string]*Repository
}

func newRepositories(ctxmgmt.Context) interface{} {
	return &Repositories{
		repos: map[string]*Repository{},
	}
}

func (c *Repositories) GetRepository(name string) *Repository {
	c.RLock()
	defer c.RUnlock()
	return c.repos[name]
}

func (c *Repositories) Set(name string, spec cpi.RepositorySpec, creds cpi.CredentialsSource) {
	c.Lock()
	defer c.Unlock()
	c.repos[name] = &Repository{
		name:  name,
		spec:  spec,
		creds: creds,
	}
}
