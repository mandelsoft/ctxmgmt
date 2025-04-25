package identity

import (
	"github.com/mandelsoft/ctxmgmt/credentials"
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	"github.com/mandelsoft/ctxmgmt/credentials/identity/hostpath"
	"github.com/mandelsoft/ctxmgmt/credentials/identity/oci"
	"github.com/mandelsoft/ctxmgmt/utils/listformat"
)

const CONSUMER_TYPE = "service.acme.corp"

// identity properties.
const (
	ID_HOSTNAME   = hostpath.ID_HOSTNAME
	ID_PORT       = hostpath.ID_PORT
	ID_PATHPREFIX = hostpath.ID_PATHPREFIX
)

// credential properties.
const (
	ATTR_USERNAME = credentials.ATTR_USERNAME
	ATTR_PASSWORD = credentials.ATTR_PASSWORD
)

var IdentityMatcher = hostpath.IdentityMatcher(CONSUMER_TYPE)

func init() {
	attrs := listformat.FormatListElements("", listformat.StringElementDescriptionList{
		ATTR_USERNAME, "user name",
		ATTR_PASSWORD, "password",
	})
	ids := listformat.FormatListElements("", listformat.StringElementDescriptionList{
		ID_HOSTNAME, "vault server host",
		ID_PORT, "(optional) server port",
		ID_PATHPREFIX, "path prefix for secret",
	})
	cpi.RegisterStandardIdentity(CONSUMER_TYPE, IdentityMatcher,
		`Service from acme.corp

This matcher matches credentials for a acme.corp service.
It uses the following identity attributes:
`+ids, attrs)
}

func GetConsumerId(addr string) credentials.ConsumerIdentity {
	h, p, s := oci.SplitLocator(addr)
	return credentials.NewConsumerIdentity(CONSUMER_TYPE,
		hostpath.ID_HOSTNAME, h,
		hostpath.ID_PORT, p,
		hostpath.ID_PATHPREFIX, s,
	)
}
