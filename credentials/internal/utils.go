package internal

import (
	"github.com/mandelsoft/goutils/errors"
)

func CredentialsForConsumer(ctx ContextProvider, id ConsumerIdentity, unknownAsError bool, matchers ...IdentityMatcher) (Credentials, error) {
	cctx := ctx.CredentialsContext()

	src, err := cctx.GetCredentialsForConsumer(id, matchers...)
	if err != nil {
		if !errors.IsErrUnknown(err) {
			return nil, errors.Wrapf(err, "lookup credentials failed for %s", id)
		}
		if unknownAsError {
			return nil, err
		}
		return nil, nil
	}
	creds, err := src.Credentials(cctx)
	if err != nil {
		return nil, errors.Wrapf(err, "lookup credentials failed for %s", id)
	}
	return creds, nil
}

func SimpleCredentials(user, passwd string) Credentials {
	return DirectCredentials{
		ATTR_USERNAME: user,
		ATTR_PASSWORD: passwd,
	}
}
