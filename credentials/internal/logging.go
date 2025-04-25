package internal

import (
	ctxlog "github.com/mandelsoft/ctxmgmt/logging"
)

var (
	REALM = ctxlog.DefineSubRealm("Credentials", "credentials")
	log   = ctxlog.DynamicLogger(REALM)
)
