package directcreds

import (
	"github.com/mandelsoft/datacontext/credentials/cpi"
	"github.com/mandelsoft/datacontext/utils"
)

func NewCredentials(props utils.Properties) cpi.CredentialsSpec {
	return cpi.NewCredentialsSpec(Type, NewRepositorySpec(props))
}
