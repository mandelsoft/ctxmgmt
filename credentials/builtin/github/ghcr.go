package github

import (
	"os"

	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	identity "github.com/mandelsoft/ctxmgmt/credentials/identity/oci"
	"github.com/mandelsoft/ctxmgmt/utils"
)

const HOST = "ghcr.io"

func init() {
	t := os.Getenv("GITHUB_TOKEN")
	if t != "" {
		host := os.Getenv("GITHUB_HOST")
		if host == "" {
			host = HOST
		}
		id := cpi.NewConsumerIdentity(identity.CONSUMER_TYPE, identity.ID_HOSTNAME, host)
		user := os.Getenv("GITHUB_REPOSITORY_OWNER")
		if user == "" {
			user = "any"
		}
		if src, err := cpi.DefaultContext.GetCredentialsForConsumer(id); err != nil || src == nil {
			creds := cpi.NewCredentials(utils.Properties{cpi.ATTR_IDENTITY_TOKEN: t, cpi.ATTR_USERNAME: user})
			cpi.DefaultContext.SetCredentialsForConsumer(id, creds)
		}
	}
}
