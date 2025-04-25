package runtime

import (
	"strings"
	"sync"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/maputils"
)

func init() {
	DefaultTypeProviderRegistry.Register("type", &DefaultProvider{})
	DefaultTypeProviderRegistry.Register("kubernetes", &KubernetesProvider{})
}

////////////////////////////////////////////////////////////////////////////////

type TypeProvider interface {
	GetTypeFor([]byte, Unmarshaler) (string, bool)
}

type TypeProviderRegistry interface {
	TypeProvider
	Register(name string, p TypeProvider)
	Get(name string) TypeProvider
	Names() []string

	Reset()
	AddAll(TypeProviderRegistry)
}

var DefaultTypeProviderRegistry = NewTypeProviderRegistry()

type _TypeProviderRegistry struct {
	lock sync.Mutex

	base      TypeProviderRegistry
	providers map[string]TypeProvider
}

func NewTypeProviderRegistry(base ...TypeProviderRegistry) TypeProviderRegistry {
	return &_TypeProviderRegistry{
		base:      general.Optional(base...),
		providers: map[string]TypeProvider{},
	}
}

func (r *_TypeProviderRegistry) Names() []string {
	r.lock.Lock()
	defer r.lock.Unlock()

	return maputils.OrderedKeys(r.providers)
}

func (r *_TypeProviderRegistry) AddAll(o TypeProviderRegistry) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, n := range o.Names() {
		r.providers[n] = o.Get(n)
	}
}

func (r *_TypeProviderRegistry) Reset() {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.providers = map[string]TypeProvider{}
}

func (r *_TypeProviderRegistry) Register(name string, provider TypeProvider) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.providers[name] = provider
}

func (r *_TypeProviderRegistry) Get(name string) TypeProvider {
	r.lock.Lock()
	defer r.lock.Unlock()

	p := r.providers[name]
	if p == nil && r.base != nil {
		p = r.base.Get(name)
	}
	return p
}

func (r *_TypeProviderRegistry) GetTypeFor(data []byte, m Unmarshaler) (string, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, p := range r.providers {
		if t, ok := p.GetTypeFor(data, m); ok {
			return t, ok
		}
	}
	if r.base != nil {
		return r.base.GetTypeFor(data, m)
	}
	return "", false
}

////////////////////////////////////////////////////////////////////////////////

// DefaultProvider uses the field type to represent the document type.
type DefaultProvider struct {
}

var _ TypeProvider = (*DefaultProvider)(nil)

func (d DefaultProvider) GetTypeFor(data []byte, m Unmarshaler) (string, bool) {
	un := &UnstructuredTypedObject{}

	err := m.Unmarshal(data, un)
	if err != nil {
		return "", false
	}
	return un.GetType(), un.GetType() != ""
}

////////////////////////////////////////////////////////////////////////////////

// KubernetesProvider uses the Kubernetes manifest fields apiVersion and kind
// to represent the type. The document type is kind[.group]/version
type KubernetesProvider struct {
}

var _ TypeProvider = (*KubernetesProvider)(nil)

type manifest struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

func (d *KubernetesProvider) GetTypeFor(data []byte, m Unmarshaler) (string, bool) {
	un := &manifest{}

	err := m.Unmarshal(data, un)
	if err != nil {
		return "", false
	}

	return MapK8SManifestInfoToType(un.ApiVersion, un.Kind)
}

func MapK8SManifestInfoToType(apiVersion, kind string) (string, bool) {
	g, v := KindVersion(apiVersion)
	if v == "" {
		v = g
		g = ""
	}

	if kind == "" || v == "" {
		return "", false
	}
	if g != "" {
		return kind + "." + g + VersionSeparator + v, true
	}
	return kind + VersionSeparator + v, true
}

func MapTypeToK8SManifestInfo(typ string) (apiVersion, kind string) {
	g, v := KindVersion(typ)

	if v == "" {
		v = "v1"
	}
	i := strings.Index(g, ".")
	if i == -1 {
		return v, g
	}
	return g[i+1:] + VersionSeparator + v, g[:i]
}
