package main

import (
	"fmt"
	"sync"

	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/examples/helper"
	"github.com/mandelsoft/ctxmgmt/examples/myconfig"
	"github.com/mandelsoft/ctxmgmt/examples/service/identity"
	"github.com/mandelsoft/goutils/errors"

	"github.com/mandelsoft/ctxmgmt/config/cpi"
)

// we already have our new acme.corp config object type,
// now we want to provide an object, which configures
// itself when used.

// ServiceAddressProvider should be an object, which is
// able to provide an service address.
// It has a setter and a getter (the setter is
// provided by our ad-hoc ExampleTarget).
// --- begin type ---
type ServiceAddressProvider struct {
	lock sync.Mutex
	// cpi.Updater is a utility, which is able to
	// configure an object based on a managed configuration
	// watermark. It remembers which config objects from the
	// config queue are already applied, and replays
	// the config objects applied to the config context
	// after the last update.
	updater cpi.Updater
	ExampleTarget
}

// --- end type ---

// --- begin constructor ---
func NewServiceAddressProvider(ctx cpi.ContextProvider) *ServiceAddressProvider {
	p := &ServiceAddressProvider{}
	// To do its work, the updater needs a connection to
	// the config context to use and the object, which should be
	// configured.
	p.updater = cpi.NewUpdater(ctx.ConfigContext(), p)
	return p
}

// --- end constructor ---

// the magic now happens in the methods provided
// by our configurable object.
// the first step for methods of configurable objects
// dependent on potential configuration is always
// to update itself using the embedded updater.
//
// Please note, the config management reverses the
// request direction. Applying a config object to
// the config context does not configure dependent objects,
// it just manages a config queue, which is used by potential
// configuration targets to configure themselves.
// The actual configuration action is always initiated
// by the object, which wants to be configured.
// The reason for this is to avoid references from the
// management to managed objects. This would prohibit
// the garbage collection of all configurable objects.

// GetServiceAddress returns a repository ref.
// --- begin method ---
func (p *ServiceAddressProvider) GetServiceAddress() (string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	err := p.updater.Update()
	if err != nil {
		return "", err
	}
	// now, we can do our regular function, aka
	// providing a repository ref.
	return p.address, nil
}

// --- end method ---

func WriteConfigConsumer() error {
	// --- begin default context ---
	ctx := config.DefaultContext()
	// --- end default context ---

	// after defining or service address provider type
	// we can now use it.
	// --- begin object ---
	prov := NewServiceAddressProvider(ctx)
	// --- end object ---

	// If we ask now for a repository we will get the empty
	// answer.
	// --- begin initial query ---
	addr, err := prov.GetServiceAddress()
	if err != nil {
		errors.Wrapf(err, "get repo")
	}
	if addr != "" {
		return fmt.Errorf("Oops, found address %q", addr)
	}
	// --- end initial query ---

	// Now, we apply our config from the last example.
	// --- begin apply config ---
	examplecfg := myconfig.NewConfig("localhost:8080")
	examplecfg.SetCredentials("testuser", "testpass")
	err = ctx.ApplyConfig(examplecfg, "special acme config")
	if err != nil {
		errors.Wrapf(err, "apply config")
	}
	// --- end apply config ---

	// without any further action, asking for a service address now will return the
	// configured one.
	// --- begin query ---
	addr, err = prov.GetServiceAddress()
	if err != nil {
		errors.Wrapf(err, "get service address")
	}
	if addr == "" {
		return fmt.Errorf("no service address provided")
	}
	helper.Output("result", func() {
		fmt.Printf("using service address: %q\n", addr)
	})
	// --- end query ---

	// --- begin query credentials ---
	id := identity.GetConsumerId(examplecfg.ServiceAddress)
	fmt.Printf("required credentials: %s\n", id)

	// the returned credentials are provided via an interface, which might change its
	// content, if the underlying credential source changes.
	// The default credential context uses the default config context,
	// therefore it finds the expected settings.
	creds, err := credentials.CredentialsForConsumer(credentials.DefaultContext(), id)
	if err != nil {
		return errors.Wrapf(err, "credentials")
	}

	helper.Output("credentials", func() {
		fmt.Printf("credentials: %s\n", creds)
	})
	// --- end query credentials ---

	return nil
}
