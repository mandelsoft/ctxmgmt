package directcreds

import (
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	"github.com/mandelsoft/ctxmgmt/utils"
)

func NewCredentials(props utils.Properties) cpi.CredentialsSpec {
	return cpi.NewCredentialsSpec(Type, NewRepositorySpec(props))
}
