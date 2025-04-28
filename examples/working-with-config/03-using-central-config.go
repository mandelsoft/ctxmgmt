package main

import (
	"fmt"

	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/service/identity"
	"github.com/mandelsoft/goutils/errors"
	"sigs.k8s.io/yaml"

	"github.com/mandelsoft/ctxmgmt/config/cfgutils"
	configcfg "github.com/mandelsoft/ctxmgmt/config/extensions/config"
	"github.com/mandelsoft/ctxmgmt/credentials"
	credcfg "github.com/mandelsoft/ctxmgmt/credentials/config"
	"github.com/mandelsoft/ctxmgmt/credentials/extensions/repositories/dockerconfig"
)

func HandleCentralConfiguration() error {
	// Although the configuration of a context can
	// be done by a sequence of explicit calls according to the mechanism
	// shown in the examples before, it provides a simple
	// library function, which can be used to configure a
	// context and all related other contexts with a single call
	// based on arbitrary central configuration files.

	// --- begin central config ---
	ctx := config.DefaultContext()

	err := cfgutils.Configure(ctx, ".appconfig")
	if err != nil {
		return errors.Wrapf(err, "configuration")
	}
	// --- end central config ---

	// This file typically contains the serialization of such a generic
	// configuration specification (or any other serialized configuration object),
	// enriched with other configuration object serializations.
	// Here, we just use our application config object.

	// --- begin configure ---
	tgt := &ExampleTarget{}

	err = ctx.ApplyAllTo(tgt)
	if err != nil {
		return errors.Wrapf(err, "request configuration")
	}
	// --- end configure ---

	// now, the address should be configured in out config target

	helper.Output("result", func() {
		fmt.Printf("configured address is %q\n", tgt.address)
	})

	// Most important are here credentials.
	// Because a program or library potentially embraces lots of storage
	// technologies as well as used network based remote services
	// there are typically multiple technology specific ways
	// to configure credentials for command line tools.
	// Using the credentials settings shown in the next tour,
	// it is possible to specify credentials for all
	// required purposes, and the configuration management provides
	// an extensible way to embed native technology specific ways
	// to provide credentials just by adding an appropriate type
	// of credential repository, which reads the specialized storage and
	// feeds it into the credential context. Those specifications
	// can be added via the credential configuration object to
	// the central configuration.
	//
	// One such repository type is the Docker config type. It
	// reads a `dockerconfig.json` file and feeds in the credentials.
	// Because it is used for a dedicated purpose (credentials for
	// OCI registries), it not only can feed the credentials, but
	// also their mapping to consumer ids.

	// We first create the specification for a new credential repository of
	// type `dockerconfig` describing the default location
	// of the standard Docker config file.

	// --- begin docker config ---
	credspec := dockerconfig.NewRepositorySpec("~/.docker/config.json", true)

	// add this repository specification to a credential configuration.
	ccfg := credcfg.New()
	err = ccfg.AddRepository(credspec)
	if err != nil {
		return errors.Wrapf(err, "invalid credential config")
	}
	// --- end docker config ---

	// By adding the default location for the standard Docker config
	// file, all credentials provided by the `docker login` command
	// are available in our context management environment, also.

	// A typical minimal <code>.acmeconfig</code> file can be composed as follows.
	// We add this config object to an empty generic configuration object
	// and print the serialized form. The result can be used as
	// default initial configuration file.

	// --- begin default config ---
	cfg := configcfg.New()
	err = cfg.AddConfig(ccfg)

	spec, err := yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrapf(err, "marshal config")
	}
	// --- end default config ---

	// the result is a typical minimal configuration file
	// for applications or libraries working with OCI
	// just providing the credentials configured with
	// <code>docker login</code>.
	fmt.Printf("this a typical application config file:\n")
	helper.Output("config file", func() {
		fmt.Printf("%s\n", string(spec))
	})

	// Besides from a file, such a config can be provided as data, also,
	// taken from any other source, for example from a Kubernetes secret

	// --- begin by data ---
	err = cfgutils.ConfigureByData(ctx, spec, "from data")
	if err != nil {
		return errors.Wrapf(err, "configuration")
	}
	// --- end by data ---

	// the configuration library function does not only read the
	// config file, it also applies [*spiff*](github.com/mandelsoft/spiff)
	// processing to the provided YAML/JSON content. *Spiff* is an
	// in-domain yaml-based templating engine. Therefore, you can use
	// any spiff dynaml expression to define values or even complete
	// sub structures.

	// --- begin spiff ---
	cfg = configcfg.New()
	ccfg = credcfg.New()
	cspec := credentials.CredentialsSpecFromList("clientCert", `(( read("~/acme/keys/myClientCert.pem") ))`)
	id := credentials.NewConsumerIdentity(identity.CONSUMER_TYPE, "hostname", "service.acme.corp")
	ccfg.AddConsumer(id, cspec)
	cfg.AddConfig(ccfg)
	// --- end spiff ---

	spec, err = yaml.Marshal(cfg)
	if err != nil {
		return errors.Wrapf(err, "marshal ocm config")
	}
	fmt.Printf("this a typical ocm config file using spiff file operations:\n")
	helper.Output("spiffconfig", func() {
		fmt.Printf("%s\n", string(spec))
	})

	// this config object is not directly usable, because the cert value is not
	// a valid certificate. We use it here just to generate the serialized form.
	// if this is used with the above library functions, the finally generated
	// config object will contain the read file content, which is hopefully a
	// valid certificate.

	return nil
}
