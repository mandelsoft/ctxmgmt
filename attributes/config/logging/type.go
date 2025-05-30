package logging

import (
	"github.com/mandelsoft/ctxmgmt/attributes"
	logdata "github.com/mandelsoft/ctxmgmt/utils/cobrautils/logopts/logging"
	"github.com/mandelsoft/logging"
	logcfg "github.com/mandelsoft/logging/config"

	"github.com/mandelsoft/ctxmgmt"
	"github.com/mandelsoft/ctxmgmt/config/cpi"
	local "github.com/mandelsoft/ctxmgmt/logging"
	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

const (
	ConfigType   = "logging" + cpi.CONFIG_TYPE_SUFFIX
	ConfigTypeV1 = ConfigType + runtime.VersionSeparator + "v1"
)

func init() {
	cpi.RegisterConfigType(cpi.NewConfigType[*Config](ConfigType, usage))
	cpi.RegisterConfigType(cpi.NewConfigType[*Config](ConfigTypeV1, usage))
}

// Config describes logging settings for a dedicated context type.
type Config struct {
	runtime.ObjectVersionedType `json:",inline"`

	// ContextType described the context type to apply the setting.
	// If not set, the settings will be applied to any logging context provider,
	// which are not derived contexts.
	ContextType string        `json:"contextType,omitempty"`
	Settings    logcfg.Config `json:"settings"`

	// ExtraId is used for the context type "default" or "global" to be able
	// to reapply the same config again using a different
	// identity given by the settings hash + the id.
	ExtraId string `json:"extraId,omitempty"`
}

// New creates a logging config specification.
func New(ctxtype string, deflvl int) *Config {
	return &Config{
		ObjectVersionedType: runtime.NewVersionedTypedObject(ConfigType),
		ContextType:         ctxtype,
		Settings: logcfg.Config{
			DefaultLevel: logging.LevelName(deflvl),
		},
	}
}

// NewWithConfig creates a logging config specification from a
// logging config object.
func NewWithConfig(ctxtype string, cfg *logcfg.Config) *Config {
	return &Config{
		ObjectVersionedType: runtime.NewVersionedTypedObject(ConfigType),
		ContextType:         ctxtype,
		Settings:            *cfg,
	}
}

func (c *Config) AddRuleSpec(r logcfg.Rule) error {
	c.Settings.Rules = append(c.Settings.Rules, r)
	return nil
}

func (c *Config) GetType() string {
	return ConfigType
}

func (c *Config) ApplyTo(ctx cpi.Context, target interface{}) error {
	// first: check for forward configuration
	// TODO: context type registry
	if lc, ok := target.(*logdata.LoggingConfiguration); ok {
		switch c.ContextType {
		case "default", "global", "slave":
			lc.LogConfig.DefaultLevel = c.Settings.DefaultLevel
			lc.LogConfig.Rules = append(lc.LogConfig.Rules, c.Settings.Rules...)
		}
		return nil
	}

	// TODO: add registry for context types

	// second: main use case is to configure vrious logging contexts
	switch c.ContextType {
	// configure local static logging context.
	// here, config is only applied once for every
	// setting hash.
	case "default":
		return local.Configure(&c.Settings, c.ExtraId)

	case "global":
		return local.ConfigureGlobal(&c.Settings, c.ExtraId)

	case "slave":
		return nil

	// configure logging context providers.
	case "":
		if _, ok := target.(attributes.AttributesContext); !ok {
			return cpi.ErrNoContext("attribute context")
		}

	// configure dedicated context types.
	default:
		dc, ok := target.(ctxmgmt.Context)
		if !ok {
			return cpi.ErrNoContext("data context")
		}
		if dc.GetType() != c.ContextType {
			return cpi.ErrNoContext(c.ContextType)
		}
	}
	lctx, ok := target.(logging.ContextProvider)
	if !ok {
		return cpi.ErrNoContext("logging context")
	}
	return logcfg.DefaultRegistry().Configure(lctx.LoggingContext(), &c.Settings)
}

const usage = `
The config type <code>` + ConfigType + `</code> can be used to configure the logging
aspect of a dedicated context type:

<pre>
    type: ` + ConfigType + `
    contextType: ` + attributes.CONTEXT_TYPE + `
    settings:
      defaultLevel: Info
      rules:
        - ...
</pre>

The context type ` + attributes.CONTEXT_TYPE + ` is the root context of a
context hierarchy.

If no context type is specified, the config will be applies to any target
acting as logging context provider, which is not a non-root context.
`
