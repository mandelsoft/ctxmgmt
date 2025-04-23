package github

import (
	"os"

	"github.com/mandelsoft/datacontext/credentials/cpi"
	identity "github.com/mandelsoft/datacontext/credentials/identity/github"
	"github.com/mandelsoft/datacontext/utils"
)

func init() {
	t := os.Getenv("GITHUB_TOKEN")
	if t != "" {
		us := os.Getenv("GITHUB_SERVER_URL")
		id := identity.GetConsumerId(us)

		if src, err := cpi.DefaultContext.GetCredentialsForConsumer(id); err != nil || src == nil {
			creds := cpi.NewCredentials(utils.Properties{cpi.ATTR_TOKEN: t})
			cpi.DefaultContext.SetCredentialsForConsumer(id, creds)
		}
	}
}
