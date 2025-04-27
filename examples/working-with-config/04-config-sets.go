package main

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/config"
	configcfg "github.com/mandelsoft/ctxmgmt/config/extensions/config"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/goutils/errors"
)

func WorkingWithConfigSets() error {
	ctx := config.New(ctxmgmt.MODE_EXTENDED)

	// a config set is a list of config objects.
	// --- begin config set ---
	set := config.NewConfigSet("application")
	set.AddConfig(myconfig.NewConfig("localhost:8080"))
	// --- end config set ---

	// config sets are used to preconfigure sets of config
	// objects at a config context.
	// --- begin add config set ---
	ctx.AddConfigSet("application", set)
	// --- end add config set ---

	// By default, the config objects of a config set are
	// not active.
	// --- begin not configured ---
	tgt := &ExampleTarget{}
	err := ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "cannot configure target")
	}

	if tgt.address != "" {
		return fmt.Errorf("address should not be configured")
	}
	// --- end not configured ---

	// any preconfigured config set can be activated
	// via API call using its name.
	// --- begin activate ---
	ctx.ApplyConfigSet("application")
	// --- end activate ---

	// now, the included config objects are applied
	// to the configuration context, and
	// our application can be configured.

	// --- begin configure ---
	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "request configuration")
	}
	// --- end configure ---

	// now, the address should be configured in our config target

	helper.Output("result", func() {
		fmt.Printf("configured address is %q\n", tgt.address)
	})

	// the config object type for the config management
	// also supports configuring config sets.
	// we add a new set with a modified service address
	// and apply this config.

	// --- begin configuring sets ---
	cfg := configcfg.New()
	set = config.NewConfigSet("application")
	set.AddConfig(myconfig.NewConfig("service.acme.corp:443"))
	cfg.AddConfigSet("modified", set)
	// --- end configuring sets ---

	err = ctx.ApplyConfig(cfg, "explicit set")
	if err != nil {
		return errors.Wrapf(err, "cannot configure explicit set")
	}

	// if we reapply the config to our target, nothing should change
	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "request configuration")
	}

	if tgt.address != "localhost:8080" {
		return fmt.Errorf("configured address should bor be changed")
	}

	// after activating the modified set the configured address should
	// be adapted.
	// --- begin activate modified ---
	ctx.ApplyConfigSet("modified")
	// --- end activate modified ---

	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "request configuration")
	}

	// now, the address should be configured in our config target

	helper.Output("modified result", func() {
		fmt.Printf("modified address is %q\n", tgt.address)
	})

	return nil
}
