package data

import (
	"encoding/json"

	"github.com/mandelsoft/ctxmgmt/config/cpi"
	"github.com/mandelsoft/ctxmgmt/utils/runtime"
	"github.com/mandelsoft/goutils/errors"
)

const (
	ConfigType   = "data" + cpi.CONFIG_TYPE_SUFFIX
	ConfigTypeV1 = ConfigType + runtime.VersionSeparator + "v1"
)

func init() {
	cpi.RegisterConfigType(cpi.NewConfigType[*Config](ConfigType, usage))
	cpi.RegisterConfigType(cpi.NewConfigType[*Config](ConfigTypeV1, usage))
}

// Config describes arbitrary configuration data
// // passed to a named ConfigApplier.
type Config struct {
	runtime.ObjectVersionedType `json:",inline"`
	Data                        json.RawMessage `json:"data,omitempty"`
	Applier                     string          `json:"applier"`
}

// New creates a new memory ConfigSpec.
func New(applier string, data interface{}) (*Config, error) {
	var rendered json.RawMessage
	var err error

	if data != nil {
		switch d := data.(type) {
		case []byte:
			var m interface{}
			err = runtime.DefaultYAMLEncoding.Unmarshal(d, &m)
			if err != nil {
				return nil, err
			}
			rendered, _ = json.Marshal(m)
		default:
			rendered, err = json.Marshal(d)
			if err != nil {
				return nil, err
			}
		}
	}
	return &Config{
		ObjectVersionedType: runtime.NewVersionedTypedObject(ConfigType),
		Applier:             applier,
		Data:                rendered,
	}, nil
}

func (c *Config) GetType() string {
	return ConfigType
}

func (c *Config) ApplyTo(ctx cpi.Context, target interface{}) error {
	a := ctx.ConfigAppliers().Get(c.Applier)
	if a == nil {
		return errors.ErrUnknown(cpi.KIND_CONFIGAPPLIER)
	}
	var data any
	err := json.Unmarshal(c.Data, &data)
	if err != nil {
		return err
	}
	return a.ApplyConfigTo(ctx, data, target)
}

const usage = `
The config type <code>` + ConfigType + `</code> can be used to pass arbitrary 
configuration data to a named config applier known to the config context:

<pre>
    type: ` + ConfigType + `
    data: ...
    applier: <applier name>
</pre>
`
