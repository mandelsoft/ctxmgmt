package dockerconfig

import (
	"github.com/docker/cli/cli/config/configfile"
	dockercred "github.com/docker/cli/cli/config/credentials"
	"github.com/docker/cli/cli/config/types"
	"github.com/mandelsoft/goutils/set"

	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	"github.com/mandelsoft/ctxmgmt/utils"
)

type Credentials struct {
	config *configfile.ConfigFile
	name   string
	store  dockercred.Store
}

var _ cpi.Credentials = (*Credentials)(nil)

// NewCredentials describes a default getter method for a authentication method.
func NewCredentials(cfg *configfile.ConfigFile, name string, store dockercred.Store) cpi.Credentials {
	return &Credentials{
		config: cfg,
		name:   name,
		store:  store,
	}
}

func (c *Credentials) get() utils.Properties {
	auth, err := c.config.GetAuthConfig(c.name)
	if err != nil {
		return utils.Properties{}
	}
	return newCredentials(auth).Properties()
}

func (c *Credentials) Credentials(context cpi.Context, source ...cpi.CredentialsSource) (cpi.Credentials, error) {
	var auth types.AuthConfig
	var err error
	if c.store == nil {
		auth, err = c.config.GetAuthConfig(c.name)
	} else {
		auth, err = c.store.Get(c.name)
	}
	if err != nil {
		return nil, err
	}
	return newCredentials(auth), nil
}

func (c *Credentials) ExistsProperty(name string) bool {
	_, ok := c.get()[name]
	return ok
}

func (c *Credentials) GetProperty(name string) string {
	return c.get()[name]
}

func (c *Credentials) PropertyNames() set.Set[string] {
	return c.get().Names()
}

func (c *Credentials) Properties() utils.Properties {
	return c.get()
}
