package npm

import (
	"fmt"
	"os"

	"github.com/mandelsoft/filepath/pkg/filepath"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/config/defaultconfigregistry"
	credcfg "github.com/mandelsoft/ctxmgmt/credentials/config"
)

const (
	ConfigFileName = ".npmrc"
)

func init() {
	defaultconfigregistry.RegisterDefaultConfigHandler(DefaultConfigHandler, desc)
}

func DefaultConfig() (string, error) {
	d, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, ConfigFileName), nil
}

func DefaultConfigHandler(cfg config.Context) (string, config.Config, error) {
	// use docker config as default config
	d, err := DefaultConfig()
	if err != nil {
		return "", nil, nil
	}
	if ok, err := vfs.FileExists(osfs.OsFs, d); ok && err == nil {
		ccfg := credcfg.New()
		ccfg.AddRepository(NewRepositorySpec(d, true))
		return d, ccfg, nil
	}
	return "", nil, nil
}

var desc = fmt.Sprintf(`
The npm configuration file at <code>~/%s</code> is
read to feed in the configured credentials for NPM registries.
`, ConfigFileName)
