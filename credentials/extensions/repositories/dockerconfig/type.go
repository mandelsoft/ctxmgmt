package dockerconfig

import (
	"encoding/json"
	"fmt"

	"github.com/mandelsoft/goutils/general"
	"github.com/mandelsoft/goutils/generics"
	"github.com/mandelsoft/goutils/optionutils"

	"github.com/mandelsoft/datacontext/credentials/cpi"
	"github.com/mandelsoft/datacontext/utils/runtime"
)

const (
	Type   = "DockerConfig"
	TypeV1 = Type + runtime.VersionSeparator + "v1"
)

func init() {
	cpi.RegisterRepositoryType(cpi.NewRepositoryType[*RepositorySpec](Type))
	cpi.RegisterRepositoryType(cpi.NewRepositoryType[*RepositorySpec](TypeV1, cpi.WithDescription(usage), cpi.WithFormatSpec(format)))
}

// RepositorySpec describes a docker config based credential repository interface.
type RepositorySpec struct {
	runtime.ObjectVersionedType `json:",inline"`
	DockerConfigFile            string          `json:"dockerConfigFile,omitempty"`
	DockerConfig                json.RawMessage `json:"dockerConfig,omitempty"`
	PropagateConsumerIdentity   *bool           `json:"propagateConsumerIdentity,omitempty"`
}

func (s RepositorySpec) WithConsumerPropagation(propagate bool) *RepositorySpec {
	s.PropagateConsumerIdentity = &propagate
	return &s
}

// NewRepositorySpec creates a new memory RepositorySpec.
func NewRepositorySpec(path string, prop ...bool) *RepositorySpec {
	var p *bool
	if len(prop) > 0 {
		p = generics.PointerTo(general.Optional(prop...))
	}
	if path == "" {
		path = "~/.docker/config.json"
	}
	return &RepositorySpec{
		ObjectVersionedType:       runtime.NewVersionedTypedObject(Type),
		DockerConfigFile:          path,
		PropagateConsumerIdentity: p,
	}
}

func NewRepositorySpecForConfig(data []byte, prop ...bool) *RepositorySpec {
	var p *bool
	if len(prop) > 0 {
		p = generics.PointerTo(general.Optional(prop...))
	}
	return &RepositorySpec{
		ObjectVersionedType:       runtime.NewVersionedTypedObject(Type),
		DockerConfig:              data,
		PropagateConsumerIdentity: p,
	}
}

func (a *RepositorySpec) GetType() string {
	return Type
}

func (a *RepositorySpec) Repository(ctx cpi.Context, creds cpi.Credentials) (cpi.Repository, error) {
	r := ctx.GetAttributes().GetOrCreateAttribute(ATTR_REPOS, newRepositories)
	repos, ok := r.(*Repositories)
	if !ok {
		return nil, fmt.Errorf("failed to assert type %T to Repositories", r)
	}
	return repos.GetRepository(ctx, a.DockerConfigFile, a.DockerConfig, optionutils.AsBool(a.PropagateConsumerIdentity, true))
}
