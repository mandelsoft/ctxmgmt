package npm

import (
	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
)

type Cache struct {
	repos map[string]*Repository
}

func createCache(_ ctxmgmt.Context) interface{} {
	return &Cache{
		repos: map[string]*Repository{},
	}
}

func (r *Cache) GetRepository(ctx cpi.Context, name string, prop bool) (*Repository, error) {
	var (
		err  error = nil
		repo *Repository
	)
	if name != "" {
		repo = r.repos[name]
	}
	if repo == nil {
		repo, err = NewRepository(ctx, name, prop)
		if err == nil {
			r.repos[name] = repo
		}
	}
	return repo, err
}
