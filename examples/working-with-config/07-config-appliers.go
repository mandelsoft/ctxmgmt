package main

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/config/extensions/data"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/goutils/errors"

	"github.com/mandelsoft/ctxmgmt/config/cpi"
)

// A typical config object bundles the configuration logic
// with particular configuration fields.
// This could be problematic if generic data stores
// should be used to store configuration data. Here, we
// have a third element, the technical access to the
// data repository. This would require to create
// config objects for every combination of
// repository, config fields and config logic.
// To circumvent this, the configuration management
// supports a decoupling of the data access from application
// logic by introducing ConfigAppliers.
// A ConfigApplier takes arbitrary config data
// and applies it according to its own logic.
// The storage access to retrieve the data is left to
// a storage technology specific implementation of
// a config object just reading the data and applying
// a named config applier.

// --- begin applier ---

const APPLIER = "config.acme.corp"

type ConfigApplier struct {
}

var _ cpi.ConfigApplier = (*ConfigApplier)(nil)

func (a *ConfigApplier) ApplyConfigTo(ctx config.Context, data, target any) error {
	if t, ok := target.(myconfig.ConfigTarget); ok {
		d, ok := data.(string)
		if ok {
			t.SetServiceAddress(d)
		} else {
			return fmt.Errorf("invalid config data: expceted a string value")
		}
	}
	return cpi.ErrNoContext(myconfig.TYPE)
}

// --- end applier ---

// --- begin init ---
func init() {
	cpi.RegisterConfigApplier(APPLIER, &ConfigApplier{})
}

// --- end init ---

func UsingConfigAppliers() error {
	// --- begin default context ---
	ctx := config.DefaultContext()
	// --- end default context ---

	// --- begin config data ---
	examplecfg, err := data.New(APPLIER, "localhost:8080")
	if err != nil {
		return errors.Wrapf(err, "invalid config data")
	}
	// --- end config data ---

	err = ctx.ApplyConfig(examplecfg, "generic acme config")
	if err != nil {
		errors.Wrapf(err, "apply config")
	}

	// --- begin apply config ---
	tgt := &ExampleTarget{}
	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "cannot configure target")
	}

	err = ctx.ApplyConfig(examplecfg, "special acme config")
	if err != nil {
		errors.Wrapf(err, "apply config")
	}
	// --- end apply config ---

	helper.Output("result", func() {
		fmt.Printf("using service address: %q\n", tgt.address)
	})

	return nil
}
