package directcreds

import (
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	"github.com/mandelsoft/ctxmgmt/utils"
	"github.com/mandelsoft/ctxmgmt/utils/runtime"
)

const (
	Type   = "Credentials"
	TypeV1 = Type + runtime.VersionSeparator + "v1"
)

func init() {
	cpi.RegisterRepositoryType(cpi.NewRepositoryType[*RepositorySpec](Type))
	cpi.RegisterRepositoryType(cpi.NewRepositoryType[*RepositorySpec](TypeV1, cpi.WithDescription(usage), cpi.WithFormatSpec(format)))
}

// RepositorySpec describes a repository interface for single direct credentials.
type RepositorySpec struct {
	runtime.ObjectVersionedType `json:",inline"`
	Properties                  utils.Properties `json:"properties"`
}

var (
	_ cpi.RepositorySpec  = &RepositorySpec{}
	_ cpi.CredentialsSpec = &RepositorySpec{}
)

// NewRepositorySpec creates a new RepositorySpec.
func NewRepositorySpec(credentials utils.Properties) *RepositorySpec {
	return &RepositorySpec{
		ObjectVersionedType: runtime.NewVersionedTypedObject(Type),
		Properties:          credentials,
	}
}

func (a *RepositorySpec) GetType() string {
	return Type
}

func (a *RepositorySpec) Repository(ctx cpi.Context, creds cpi.Credentials) (cpi.Repository, error) {
	return NewRepository(cpi.NewCredentials(a.Properties)), nil
}

func (a *RepositorySpec) Credentials(context cpi.Context, source ...cpi.CredentialsSource) (cpi.Credentials, error) {
	return cpi.NewCredentials(a.Properties), nil
}

func (a *RepositorySpec) GetCredentialsName() string {
	return ""
}

func (a *RepositorySpec) GetRepositorySpec(context cpi.Context) cpi.RepositorySpec {
	return a
}
