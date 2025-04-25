package cfgutils

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/ioutils"
	"github.com/mandelsoft/goutils/pkgutils"
	"github.com/mandelsoft/spiff/features"
	"github.com/mandelsoft/spiff/spiffing"
	"github.com/mandelsoft/vfs/pkg/osfs"
	"github.com/mandelsoft/vfs/pkg/vfs"

	"github.com/mandelsoft/ctxmgmt/config/defaultconfigregistry"
	configcfg "github.com/mandelsoft/ctxmgmt/config/extensions/config"
)

// Configure configures a config context from some config file.
// It handles the ~/ prefix for the home directory
// and preprocesses the read config data with [github.com/mandelsoft/spiff/spiffing.Spiff].
func Configure(ctxp config.ContextProvider, path string, fss ...vfs.FileSystem) error {
	_, err := Configure2(ctxp, path, fss...)
	return err
}

func Configure2(ctx config.ContextProvider, path string, fss ...vfs.FileSystem) (config.Config, error) {
	cfg, err := configcfg.NewAggregator(false)
	if err != nil {
		return nil, err
	}
	fs := general.OptionalDefaulted[vfs.FileSystem](osfs.OsFs, fss...)
	if ctx == nil {
		ctx = config.DefaultContext()
	}
	path, err = ioutils.ResolvePath(path)
	if err != nil {
		return nil, err
	}
	if path != "" && path != "None" {
		data, err := vfs.ReadFile(fs, path)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot read ocm config file %q", path)
		}

		if eff, err := ConfigureByData2(ctx, data, path); err != nil {
			return nil, err
		} else {
			err = cfg.AddConfig(eff)
			if err != nil {
				return nil, err
			}
		}
	} else {
		for _, h := range defaultconfigregistry.Get() {
			desc, def, err := h(ctx.ConfigContext())
			if err != nil {
				return nil, err
			}
			if def != nil {
				name, err := pkgutils.GetPackageName(h)
				if err != nil {
					name = "unknown handler"
				}
				err = ctx.ConfigContext().ApplyConfig(def, fmt.Sprintf("%s: %s", name, desc))
				if err != nil {
					return nil, errors.Wrapf(err, "cannot apply default config from %s(%s)", name, desc)
				}
				err = cfg.AddConfig(def)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return cfg.Get(), nil
}

func ConfigureByData(ctx config.ContextProvider, data []byte, info string) error {
	_, err := ConfigureByData2(ctx, data, info)
	return err
}

func ConfigureByData2(ctx config.ContextProvider, data []byte, info string) (config.Config, error) {
	var err error

	sctx := spiffing.New().WithFeatures(features.INTERPOLATION, features.CONTROL)
	data, err = spiffing.Process(sctx, spiffing.NewSourceData(info, data))
	if err != nil {
		return nil, errors.Wrapf(err, "processing ocm config %q", info)
	}
	cfg, err := ctx.ConfigContext().GetConfigForData(data, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid ocm config file %q", info)
	}
	err = ctx.ConfigContext().ApplyConfig(cfg, info)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot apply ocm config %q", info)
	}
	return cfg, nil
}
