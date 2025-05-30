package internal

import (
	"fmt"
	"strings"

	"github.com/mandelsoft/goutils/errors"
	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/stringutils"
	"github.com/modern-go/reflect2"

	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

type ConfigType interface {
	runtime.VersionedTypedObjectType[Config]
	Usage() string
}

var _ runtime.VersionedTypedObjectType[Config] = (ConfigType)(nil)

type (
	ConfigDecoder      = runtime.TypedObjectDecoder[Config]
	ConfigTypeProvider = runtime.KnownTypesProvider[Config, ConfigType]
)

type ConfigTypeScheme interface {
	runtime.TypeScheme[Config, ConfigType]

	Usage() string
}

type _Scheme = runtime.TypeScheme[Config, ConfigType]

type configTypeScheme struct {
	_Scheme
}

func NewConfigTypeScheme(defaultDecoder ConfigDecoder, base ...ConfigTypeScheme) ConfigTypeScheme {
	scheme := runtime.MustNewDefaultTypeScheme[Config, ConfigType](&GenericConfig{}, true, defaultDecoder, general.Optional(base...))
	return &configTypeScheme{scheme}
}

// KnownTypes required just for Goland.
func (s *configTypeScheme) KnownTypes() runtime.KnownTypes[Config, ConfigType] {
	return s._Scheme.KnownTypes()
}

func (t *configTypeScheme) DecodeConfig(data []byte, unmarshaler runtime.Unmarshaler) (Config, error) {
	return t._Scheme.Decode(data, unmarshaler) // Goland
}

type versionRegistry struct {
	_Scheme
}

func NewStrictConfigTypeScheme(base ...ConfigTypeScheme) runtime.VersionedTypeRegistry[Config, ConfigType] {
	scheme := runtime.MustNewDefaultTypeScheme[Config, ConfigType](nil, false, nil, general.Optional(base...))
	return &versionRegistry{scheme}
}

func (s *versionRegistry) KnownTypes() runtime.KnownTypes[Config, ConfigType] {
	return s._Scheme.KnownTypes() // Goland
}

func (t *configTypeScheme) CreateConfig(obj runtime.TypedObject) (Config, error) {
	return t._Scheme.Convert(obj)
}

func (t *configTypeScheme) Usage() string {
	found := map[string]bool{}

	s := "\nThe following configuration types are supported:\n"
	for _, n := range t.KnownTypeNames() {
		ct := t.GetType(n)
		u := ct.Usage()
		if strings.TrimSpace(u) == "" || found[u] {
			continue
		}
		found[u] = true
		for strings.HasSuffix(u, "\n") {
			u = u[:len(u)-1]
		}
		s = fmt.Sprintf("%s\n- <code>%s</code>\n%s", s, ct.GetKind(), stringutils.IndentLines(u, "  "))
	}
	return s + "\n"
}

// DefaultConfigTypeScheme contains all globally known access serializer.
var DefaultConfigTypeScheme = NewConfigTypeScheme(nil)

////////////////////////////////////////////////////////////////////////////////

type Evaluator interface {
	Evaluate(ctx Context) (Config, error)
}

type GenericConfig struct {
	runtime.UnstructuredVersionedTypedObject `json:",inline"`
	unknown                                  bool
}

func IsGeneric(cfg Config) bool {
	_, ok := cfg.(*GenericConfig)
	return ok
}

func NewGenericConfig(data []byte, unmarshaler runtime.Unmarshaler) (Config, error) {
	unstr := &runtime.UnstructuredVersionedTypedObject{}
	if unmarshaler == nil {
		unmarshaler = runtime.DefaultYAMLEncoding
	}
	err := unmarshaler.Unmarshal(data, unstr)
	if err != nil {
		return nil, err
	}
	return &GenericConfig{*unstr, false}, nil
}

func ToGenericConfig(c Config) (*GenericConfig, error) {
	if reflect2.IsNil(c) {
		return nil, nil
	}
	if g, ok := c.(*GenericConfig); ok {
		return g, nil
	}
	u, err := runtime.ToUnstructuredVersionedTypedObject(c)
	if err != nil {
		return nil, err
	}
	return &GenericConfig{*u, false}, nil
}

func (s *GenericConfig) IsUnknown() bool {
	return s.unknown
}

func (s *GenericConfig) Evaluate(ctx Context) (Config, error) {
	raw, err := s.GetRaw()
	if err != nil {
		return nil, err
	}
	cfg, err := ctx.ConfigTypes().Decode(raw, runtime.DefaultJSONEncoding)
	if err != nil {
		return nil, err
	}
	if IsGeneric(cfg) {
		s.unknown = true
		return nil, errors.ErrUnknown(KIND_CONFIGTYPE, s.GetType())
	} else {
		s.unknown = false
	}
	return cfg, nil
}

func (s *GenericConfig) ApplyTo(ctx Context, target interface{}) error {
	spec, err := s.Evaluate(ctx)
	if err != nil {
		return err
	}
	return spec.ApplyTo(ctx, target)
}

var _ Config = &GenericConfig{}

////////////////////////////////////////////////////////////////////////////////
