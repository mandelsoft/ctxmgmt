package main

import (
	"encoding/json"
	"fmt"

	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/goutils/errors"

	"github.com/mandelsoft/ctxmgmt/config"
	configcfg "github.com/mandelsoft/ctxmgmt/config/extensions/config"
)

func HandleArbitraryConfiguration() error {
	// The configuration management provides a configuration object
	// of its own. It can be used to aggregate arbitrary configuration
	// objects.

	// --- begin config config ---
	generic := configcfg.New()
	// --- end config config ---

	// the generic config holds a list of any other config objects,
	// or their specification formats.

	// we recycle our application config from the last example.
	// --- begin sub config ---
	subcfg := myconfig.NewConfig("localhost:8080")
	// --- end sub config ---

	// now, we can add this credential config object to
	// our generic config list.
	// --- begin add config ---
	err := generic.AddConfig(subcfg)
	if err != nil {
		return errors.Wrapf(err, "adding config")
	}
	// --- end add config ---

	// as we have seen in the previous example, config objects are typically
	// serializable and deserializable.
	// this also holds for the generic config object of the config context.

	// --- begin serialized ---
	spec, err := json.MarshalIndent(generic, "  ", "  ")
	if err != nil {
		return errors.Wrapf(err, "marshal aggregated config")
	}

	fmt.Printf("this is a generic configuration object:\n")
	helper.Output("format", func() {
		fmt.Printf("%s\n", string(spec))
	})

	// --- end serialized ---
	// the result is a config object hosting a list (with 1 entry)
	// of other config object specifications.

	// The generic config object can be added to a config context, again, like
	// any other config object. If it is asked to configure a configuration
	// context it uses the methods of the configuration context to apply the
	// contained list of config objects.
	// Therefore, all config objects applied to a configuration context are
	// asked to configure the configuration context itself when queued to the
	// list of applied configuration objects.

	// If we now ask the context to configure our application, the nested
	// config applies its value.

	// --- begin apply config ---
	ctx := config.DefaultContext()
	err = ctx.ApplyConfig(generic, "generic setting")
	if err != nil {
		return errors.Wrapf(err, "cannot apply config")
	}
	// --- end apply config ---

	// --- begin configure ---
	tgt := &ExampleTarget{}

	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "request configuration")
	}
	// --- end configure ---

	// now, the address should be configured in our config target

	helper.Output("result", func() {
		fmt.Printf("configured address is %q\n", tgt.address)
	})
	return nil
}
