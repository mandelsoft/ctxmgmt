package oci

import (
	"github.com/mandelsoft/ctxmgmt/credentials/cpi"
	"github.com/mandelsoft/ctxmgmt/credentials/identity/hostpath"
	"github.com/mandelsoft/ctxmgmt/utils/listformat"
)

// CONSUMER_TYPE is the OCI registry type.
const CONSUMER_TYPE = "OCIRegistry"

// used identity properties.
const (
	ID_TYPE       = hostpath.ID_TYPE
	ID_HOSTNAME   = hostpath.ID_HOSTNAME
	ID_PORT       = hostpath.ID_PORT
	ID_PATHPREFIX = hostpath.ID_PATHPREFIX
	ID_SCHEME     = hostpath.ID_SCHEME
)

// used credential properties.
const (
	ATTR_USERNAME              = cpi.ATTR_USERNAME
	ATTR_PASSWORD              = cpi.ATTR_PASSWORD
	ATTR_IDENTITY_TOKEN        = cpi.ATTR_IDENTITY_TOKEN
	ATTR_CERTIFICATE_AUTHORITY = cpi.ATTR_CERTIFICATE_AUTHORITY
)

func init() {
	attrs := listformat.FormatListElements("", listformat.StringElementDescriptionList{
		ATTR_USERNAME, "the basic auth username",
		ATTR_PASSWORD, "the basic auth password",
		ATTR_IDENTITY_TOKEN, "the bearer token used for non-basic auth authorization",
		ATTR_CERTIFICATE_AUTHORITY, "the certificate authority certificate used to verify certificates",
	})

	cpi.RegisterStandardIdentity(CONSUMER_TYPE, IdentityMatcher, `OCI registry credential matcher

It matches the <code>`+CONSUMER_TYPE+`</code> consumer type and additionally acts like 
the <code>`+hostpath.IDENTITY_TYPE+`</code> type.`,
		attrs)
}

var identityMatcher = hostpath.IdentityMatcher(CONSUMER_TYPE)

func IdentityMatcher(pattern, cur, id cpi.ConsumerIdentity) bool {
	return identityMatcher(pattern, cur, id)
}
