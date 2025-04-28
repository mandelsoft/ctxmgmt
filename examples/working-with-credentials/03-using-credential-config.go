package main

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/config/cfgutils"
	configcfg "github.com/mandelsoft/ctxmgmt/config/extensions/config"
	"github.com/mandelsoft/ctxmgmt/credentials"
	credcfg "github.com/mandelsoft/ctxmgmt/credentials/config"
	"github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/dockerconfig"
	"github.com/mandelsoft/ctxmgmt/credentials/identity/oci"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/goutils/errors"
	"sigs.k8s.io/yaml"
)

func UsingCredentialConfig() error {

	// --- begin default context ---
	ctx := credentials.DefaultContext()
	// --- end default context ---

	// the credential management uses the config management.
	// Therefore, a credential context incorporates a config context.

	// --- begin retrieve config context ---
	_ = ctx.ConfigContext()
	// --- end retrieve config context ---

	// we use this context now to configure
	// our credential repository.

	// --- begin repository spec ---
	repospec := dockerconfig.NewRepositorySpec("~/.docker/config.json", true)
	// --- end repository spec ---

	// a repository specification is serializable and therefore embeddable
	// into a configuration object. The credential management provides
	// an own configuration object type, which can be configured
	// with arbitrary repository specifications.

	// add this repository specification to a credential configuration.
	// --- begin docker config ---
	ccfg := credcfg.New()
	err := ccfg.AddRepository(repospec)
	if err != nil {
		return errors.Wrapf(err, "invalid credential config")
	}
	// --- end docker config ---

	// By adding the default location for the standard Docker config
	// file, all credentials provided by the `docker login` command
	// are available in our context management environment, also.
	// For testing, we just use our local version.

	// A typical minimal <code>.acmeconfig</code> file can be composed as follows.
	// We add this config object to an empty generic configuration object
	// and print the serialized form. The result can be used as
	// default initial configuration file.

	// --- begin default config ---
	cfg := configcfg.New()
	err = cfg.AddConfig(ccfg)

	spec, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrapf(err, "marshal central config")
	}
	// --- end default config ---

	// the result is a typical minimal configuration file
	// just providing the credentials configured with
	// <code>docker login</code>.
	fmt.Printf("this a typical application config file:\n")
	helper.Output("config file", func() {
		fmt.Printf("%s\n", string(spec))
	})

	// we can use this config file now to configure
	// the configuration context of our credential context.
	// The credential context acts as config context provider
	// and can directly be passed to the configuration function.

	// --- begin reading config ---
	err = cfgutils.ConfigureByData(ctx, spec, "central config")
	if err != nil {
		return errors.Wrapf(err, "invalid central config")
	}
	// --- end reading config ---

	// now, we can instantly request our credentials.
	// the configuration context configures itself on-the fly
	// from its configuration context.

	// We create the consumer id for the repository.
	// --- begin get consumer id ---
	cid := oci.GetConsumerId("ghcr.io", "acme.org/service")
	// --- end get consumer id ---

	// and then get the credentials from the credentials context
	// like in the previous example.
	// --- begin get credentials ---
	creds, err := credentials.CredentialsForConsumer(ctx, cid)
	if err != nil {
		return errors.Wrapf(err, "no credentials")
	}
	// an error is only provided if something went wrong while determining
	// the credentials. Delivering NO credentials is a valid result.
	if creds == nil {
		return fmt.Errorf("no credentials found")
	}
	// --- end get credentials ---

	helper.Output("credentials", func() {
		fmt.Printf("credentials: %s\n", creds)
	})
	return nil
}
