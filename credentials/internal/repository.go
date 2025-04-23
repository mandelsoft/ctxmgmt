package internal

import (
	"github.com/mandelsoft/goutils/set"

	"github.com/mandelsoft/datacontext/utils"
)

type Repository interface {
	ExistsCredentials(name string) (bool, error)
	LookupCredentials(name string) (Credentials, error)
	WriteCredentials(name string, creds Credentials) (Credentials, error)
}

type Credentials interface {
	CredentialsSource
	ExistsProperty(name string) bool
	GetProperty(name string) string
	PropertyNames() set.Set[string]
	Properties() utils.Properties
}

type DirectCredentials utils.Properties

var _ Credentials = (*DirectCredentials)(nil)

func NewCredentials(props utils.Properties) DirectCredentials {
	if props == nil {
		props = utils.Properties{}
	} else {
		props = props.Copy()
	}
	return DirectCredentials(props)
}

func (c DirectCredentials) ExistsProperty(name string) bool {
	_, ok := c[name]
	return ok
}

func (c DirectCredentials) GetProperty(name string) string {
	return c[name]
}

func (c DirectCredentials) PropertyNames() set.Set[string] {
	return utils.Properties(c).Names()
}

func (c DirectCredentials) Properties() utils.Properties {
	return utils.Properties(c).Copy()
}

func (c DirectCredentials) Credentials(Context, ...CredentialsSource) (Credentials, error) {
	return c, nil
}

func (c DirectCredentials) Copy() DirectCredentials {
	return DirectCredentials(utils.Properties(c).Copy())
}

func (c DirectCredentials) String() string {
	return utils.Properties(c).String()
}
