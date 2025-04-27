package main

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/ctxmgmt/examples/service/identity"
	"github.com/mandelsoft/goutils/errors"
	"sigs.k8s.io/yaml"

	configcfg "github.com/mandelsoft/ctxmgmt/config/extensions/config"
)

func WriteConfigType() error {
	// after preparing a new special config type
	// we can feed it into the config management.
	// because of the registration the config management
	// now knows about this new type.

	// A usual, we gain access to our required
	// contexts.
	// --- begin default context ---
	ctx := config.DefaultContext()
	// --- end default context ---

	// to setup our environment we create our new config based on the actual
	// settings and apply it to the config context.
	// --- begin apply ---
	examplecfg := myconfig.NewConfig("localhost:8080")
	examplecfg.SetCredentials("testuser", "testpassword")
	ctx.ApplyConfig(examplecfg, "special acme config")
	// --- end apply ---

	// If you omit the above call, no configuration
	// will be found later.

	// now we should be prepared to get configured

	// --- begin query credentials ---
	id := identity.GetConsumerId(examplecfg.ServiceAddress)
	fmt.Printf("required credentials: %s\n", id)

	// the returned credentials are provided via an interface, which might change its
	// content, if the underlying credential source changes.
	creds, err := credentials.CredentialsForConsumer(credentials.DefaultContext(), id)
	if err != nil {
		return errors.Wrapf(err, "credentials")
	}

	helper.Output("credentials", func() {
		fmt.Printf("credentials: %s\n", creds)
	})
	// --- end query credentials ---

	// the config object supports a special interface for potential
	// configuration targets, the myconfig.ConfigTarget interface.
	// Just by providing an implementation for this interface, we can
	// configure such an object using the config management.

	// --- begin request config ---
	tgt := &ExampleTarget{}
	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "cannot configure")
	}
	helper.Output("result", func() {
		fmt.Printf("configured address is %q\n", tgt.address)
	})
	// --- end request config ---

	// Because of the new config type, such a specification can
	// now be added to the ocm config, also.
	// So, we could use our special tour config file content
	// directly as part of the ocm config.

	// --- begin in central config ---
	cfg := configcfg.New()
	err = cfg.AddConfig(examplecfg)

	spec, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrapf(err, "marshal central config")
	}

	// the result is a minimal configuration file
	// just providing our new example configuration.

	fmt.Printf("this a typical config file:\n")
	helper.Output("config file", func() {
		fmt.Printf("%s\n", string(spec))
	})
	// --- end in central config ---

	// This way any specialized configuration object can be added
	// by a user of the library. It can be used to configure
	// existing objects or even new object types, even in combination.
	//
	// What is still required is a way
	// to implement new config targets, objects, which want
	// to be configured and which autoconfigure themselves when
	// used. Our simple repository target is just an example
	// for some kind of ad-hoc configuration.
	// a complete scenario is shown in the next example.
	return nil
}
