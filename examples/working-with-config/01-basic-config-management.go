package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-test/deep"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/goutils/errors"

	"github.com/mandelsoft/ctxmgmt/config"
)

////////////////////////////////////////////////////////////////////////////////

func BasicConfigurationHandling() error {
	// configuration is handled by the configuration context.
	// --- begin default context ---
	ctx := config.DefaultContext()
	// --- end default context ---

	// the configuration context handles configuration objects.
	// a configuration object is any object implementing
	// the config.Config interface. The task of a config object
	// is to apply configuration to some target object.

	// we created a special configuration object type for our
	// application in examples/myconfig.

	// now we create such an object to configure our application
	// --- begin my config ---
	cfg := myconfig.NewConfig("service.provider.com")
	cfg.Credentials = &myconfig.Credentials{
		Username: "appuser",
		Password: "apppass",
	}
	// --- end my config ---

	// configuration objects are typically serializable and deserializable.

	// --- begin marshal ---
	spec, err := json.MarshalIndent(cfg, "  ", "  ")
	if err != nil {
		return errors.Wrapf(err, "marshal config")
	}

	fmt.Printf("this is our configuration object:\n")
	helper.Output("format", func() {
		fmt.Printf("  %s\n", string(spec))
	})

	// --- end marshal ---

	// like all the other manifest based descriptions this format always includes
	// a type field, which can be used to deserialize a specification into
	// the appropriate object.
	// This can be done by the config context. It accepts YAML or JSON.

	// --- begin unmarshal ---
	o, err := ctx.GetConfigForData(spec, nil)
	if err != nil {
		return errors.Wrapf(err, "deserialize config")
	}

	if diff := deep.Equal(o, cfg); len(diff) != 0 {
		fmt.Printf("diff:\n%v\n", diff)
		return fmt.Errorf("invalid des/erialization")
	}
	// --- end unmarshal ---

	// regardless what variant is used (direct object or descriptor)
	// the config object can be added to a config context.
	// --- begin apply config ---
	err = ctx.ApplyConfig(cfg, "explicit setting")
	if err != nil {
		return errors.Wrapf(err, "cannot apply config")
	}
	// --- end apply config ---

	// Every config object implements the
	// ApplyTo(ctx config.Context, target interface{}) error method.
	// It takes an object, which wants to be configured.
	// The config object then decides, whether it provides
	// settings for the given object and calls the appropriate
	// methods on this object (after a type cast).
	//
	// This way the config mechanism reverts the configuration
	// request, it does not actively configure something, instead
	// an object, which wants to be configured calls the config
	// context to apply pending configs.
	// The config context manages a queue of config objects
	// and applies them to an object to be configured.
	// This way a config object may configure any (and even multiple)
	// objects requesting a configuration.

	// Our example config object configures the credential context
	// with credentials, and a myconfig.ConfigTarget  with a service address.

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

// --- begin target ---
type ExampleTarget struct {
	address string
}

var _ myconfig.ConfigTarget = (*ExampleTarget)(nil)

func (e *ExampleTarget) SetServiceAddress(r string) {
	e.address = r
}

// --- end target ---
