package oci

import (
	"path"

	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
)

func SimpleCredentials(user, passwd string) cpi.Credentials {
	return cpi.SimpleCredentials(user, passwd)
}

func GetCredentials(ctx cpi.ContextProvider, locator, repo string) (cpi.Credentials, error) {
	return cpi.CredentialsForConsumer(ctx.CredentialsContext(), GetConsumerId(locator, repo), identityMatcher)
}

func GetConsumerId(locator, repo string) cpi.ConsumerIdentity {
	host, port, base := SplitLocator(locator)
	id := cpi.NewConsumerIdentity(CONSUMER_TYPE, ID_HOSTNAME, host)
	if port != "" {
		id[ID_PORT] = port
	}
	if repo == "" {
		id[ID_PATHPREFIX] = base
	} else {
		id[ID_PATHPREFIX] = path.Join(base, repo)
	}
	return id
}
