package myconfig

import (
	"github.com/mandelsoft/ctxmgmt/config"
	"github.com/mandelsoft/ctxmgmt/config/cpi"
	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/examples/service/identity"
	"github.com/mandelsoft/ctxmgmt/utils"
	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

// TYPE is the name of our new configuration object type.
// To be globally unique, it should always end with a
// DNS domain owned by the provider of the new type.
// --- begin type name ---
const TYPE = "application.config.acme.org"

// --- end type name ---

// Config is the new Go type for the config specification
// covering our example configuration.
// It just encapsulates our simple configuration structure
// used to configure the examples of our tour.
// --- begin config type ---
type Config struct {
	// ObjectVersionedType is the base type providing the type feature
	// for (config) specifications.
	runtime.ObjectVersionedType `json:",inline"`

	// ServiceAddress is the address of the service intended to be used by our
	// application.
	ServiceAddress string `json:"serviceAddress"`

	// Credentials are the credentials required to access the service
	// located at ServiceAddress.
	Credentials *Credentials `json:"credentials,omitempty"`
}

type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// A config object object must implememt the config.Config interface.
var _ config.Config = (*Config)(nil)

// --- end config type ---

// NewConfig provides a config object for our application configuration.
// --- begin constructor ---
func NewConfig(addr string) *Config {
	return &Config{
		ObjectVersionedType: runtime.NewVersionedTypedObject(TYPE),
		ServiceAddress:      addr,
	}
}

// --- end constructor ---

// additional setters can be used to configure the configuration object.
// Here, programmatic objects (like an ocm.RepositorySpec) are
// converted to a form storable in the configuration object.
// --- begin setters ---

// SetCredentials sets the credentials required by the application.
func (c *Config) SetCredentials(user, pass string) {
	c.Credentials = &Credentials{
		Username: user,
		Password: pass,
	}
}

// SetAddress sets address of used service.
func (c *Config) SetAddress(desc string) {
	c.ServiceAddress = desc
}

// --- end setters ---

// --- begin getters ---
func (c *Config) GetCredentials() credentials.Credentials {
	if c.Credentials != nil && c.Credentials.Username != "" && c.Credentials.Password != "" {
		return credentials.NewCredentials(utils.Properties{
			identity.ATTR_USERNAME: c.Credentials.Username,
			identity.ATTR_PASSWORD: c.Credentials.Password,
		})
	}
	return nil
}

// --- end getters ---

// --- begin target interface ---

// ConfigTarget consumes a repository name.
type ConfigTarget interface {
	SetServiceAddress(r string)
}

// --- end target interface ---

// ApplyTo is used to apply the provided configuration settings
// to a dedicated object, which wants to be configured.
// --- begin method apply ---.
func (c *Config) ApplyTo(_ cpi.Context, tgt interface{}) error {
	switch t := tgt.(type) {
	// if the target is a credentials context
	// configure the credentials to be used for the
	// described OCI repository.
	case credentials.Context:
		// determine the consumer id for our target repository.
		if c.Credentials != nil && c.Credentials.Username != "" && c.Credentials.Password != "" {
			id := identity.GetConsumerId(c.ServiceAddress)

			// create the credentials.
			if creds := c.GetCredentials(); creds != nil {
				// configure the targeted credential context with
				// the provided credentials (see previous examples).
				t.SetCredentialsForConsumer(id, creds)
			}
		}

	// if the target consumes an OCI repository, propagate
	// the provided OCI repository ref.
	case ConfigTarget:
		t.SetServiceAddress(c.ServiceAddress)

	// all other targets are ignored, we don't have
	// something to set at these objects.
	default:
		return cpi.ErrNoContext(TYPE)
	}
	return nil
}

// --- end method apply ---

// to enable automatic deserialization of our new config type,
// we have to tell the configuration management about our
// new type. This is done by a registration function,
// which gets called with a dedicated type object for
// the new config type.
// a type object describes the config type, its type name, how
// it is serialized and deserialized and some description.
// we use a standard type object, here, instead of implementing
// an own one. It is parameterized by the Go pointer type for
// our specification object.

// --- begin init ---.
func init() {
	// register the new config type, so that is can be used
	// by the config management to deserialize appropriately
	// typed specifications.
	cpi.RegisterConfigType(cpi.NewConfigType[*Config](TYPE, "this is a config object type based on the example config data."))
}

// --- end init ---.
