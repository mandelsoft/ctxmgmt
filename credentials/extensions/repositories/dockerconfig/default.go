package dockerconfig

import (
	dockercli "github.com/docker/cli/cli/config"
	"github.com/mandelsoft/filepath/pkg/filepath"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/mandelsoft/datacontext/config"
	"github.com/mandelsoft/datacontext/config/defaultconfigregistry"
	credcfg "github.com/mandelsoft/datacontext/credentials/config"
)

func init() {
	defaultconfigregistry.RegisterDefaultConfigHandler(DefaultConfigHandler, desc)
}

func DefaultConfigHandler(cfg config.Context) (string, config.Config, error) {
	// use docker config as default config
	d := filepath.Join(dockercli.Dir(), dockercli.ConfigFileName)
	if ok, err := vfs.FileExists(osfs.New(), d); ok && err == nil {
		ccfg := credcfg.New()
		ccfg.AddRepository(NewRepositorySpec(d, true))
		return d, ccfg, nil
	}
	return "", nil, nil
}

var desc = `
The docker configuration file at <code>~/.docker/config.json</code> is
read to feed in the configured credentials for OCI registries.
`
